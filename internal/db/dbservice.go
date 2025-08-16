package db

import (
	"context"
	"fmt"
	"log/slog"

	"go-boilerplate/internal/config"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DBService encapsulates the database connection pool and sqlc Queries.
type DBService struct {
	logger   *slog.Logger
	pool     *pgxpool.Pool
	*Queries // Embed sqlc Queries for direct access to query methods
}

type myQueryTracer struct {
	log *slog.Logger
}

func (tracer *myQueryTracer) TraceQueryStart(
	ctx context.Context,
	_ *pgx.Conn,
	data pgx.TraceQueryStartData) context.Context {
	tracer.log.Info("Executing command", "sql", data.SQL, "args", data.Args)

	return ctx
}

func (tracer *myQueryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
}

// NewDBService creates a new DBService instance.
// It establishes a database connection pool and initializes sqlc Queries.
func NewDBService(ctx context.Context, config *config.DatabaseConfig, logger *slog.Logger) (*DBService, error) {
	if logger == nil {
		logger = slog.Default() // Fallback to default logger
	}
	logger = logger.With(slog.String("service", "DBService"))

	if config.DSN == "" {
		logger.Error("Database DSN is empty in configuration")
		return nil, fmt.Errorf("database DSN is empty")
	}

	poolConfig, err := pgxpool.ParseConfig(config.DSN)
	if err != nil {
		logger.Error("Failed to parse database DSN", slog.Any("error", err), slog.String("dsn", config.DSN))
		return nil, fmt.Errorf("parsing DSN: %w", err)
	}

	// Apply pool settings
	poolConfig.MaxConns = config.MaxOpenConns
	// pgxpool uses MinConns to maintain a minimum number of idle connections.
	// It's not exactly MaxIdleConns, but setting MinConns can help.
	// If MaxIdleConns is truly desired as a cap on idle connections,
	// pgxpool manages this more dynamically with MaxConnIdleTime.
	poolConfig.MinConns = config.MaxIdleConns
	poolConfig.MaxConnIdleTime = config.ConnMaxIdleTime
	poolConfig.MaxConnLifetime = config.ConnMaxLifetime
	// Consider adding health check period
	// poolConfig.HealthCheckPeriod = 1 * time.Minute
	if config.EchoSQL {
		poolConfig.ConnConfig.Tracer = &myQueryTracer{
			log: logger,
		}
	}

	logger.Info("Attempting to connect to database...", slog.String("host", poolConfig.ConnConfig.Host), slog.String("database", poolConfig.ConnConfig.Database))

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		logger.Error("Failed to create database connection pool", slog.Any("error", err))
		return nil, fmt.Errorf("creating connection pool: %w", err)
	}

	// Ping the database to verify the connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close() // Close the pool if ping fails
		logger.Error("Failed to ping database", slog.Any("error", err))
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	logger.Info("Successfully connected to the database.")

	return &DBService{
		logger:  logger,
		pool:    pool,
		Queries: New(pool), // Initialize sqlc Queries with the pool
	}, nil
}

// Close gracefully closes the database connection pool.
func (s *DBService) Close() {
	if s.pool != nil {
		s.logger.Info("Closing database connection pool...")
		s.pool.Close()
		s.logger.Info("Database connection pool closed.")
	}
}

// GetPool returns the underlying pgxpool.Pool.
func (s *DBService) GetPool() *pgxpool.Pool {
	return s.pool
}

// BeginTx starts a new database transaction.
// The caller is responsible for committing or rolling back the transaction.
func (s *DBService) BeginTx(ctx context.Context) (pgx.Tx, *Queries, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("beginning transaction: %w", err)
	}
	// Create a new Querier that uses this transaction
	qtx := s.Queries.WithTx(tx) // Relies on sqlc's WithTx method being compatible with pgx.Tx
	return tx, qtx, nil
}

// WithinTransaction executes the given function within a database transaction.
// It automatically handles commit on success or rollback on error/panic.
func (s *DBService) WithinTransaction(ctx context.Context, fn func(qtx *Queries) error) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	// Defer a function to handle rollback in case of panic or error.
	// This defer needs to be carefully placed to ensure it runs correctly.
	successful := false
	defer func() {
		if r := recover(); r != nil { // A panic occurred
			s.logger.Error("Panic recovered during transaction, rolling back", slog.Any("panic_value", r))
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				s.logger.Error("Failed to rollback transaction after panic", slog.Any("rollback_error", rbErr))
			}
			panic(r) // Re-throw panic after attempting rollback
		} else if !successful { // An error occurred (but not a panic)
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				s.logger.Error("Failed to rollback transaction after error", slog.Any("rollback_error", rbErr))
				// Optionally, wrap the original error with the rollback error, or just log it.
			}
		}
	}()

	qtx := New(tx) // Create a new Querier for the transaction

	if err := fn(qtx); err != nil {
		s.logger.Debug("Error during transactional function, will rollback", slog.Any("error", err))
		// The defer will handle the rollback. Simply return the error.
		return err
	}

	// If we reach here, the function `fn` completed without error.
	successful = true // Mark as successful before attempting commit
	if err := tx.Commit(ctx); err != nil {
		s.logger.Error("Failed to commit transaction", slog.Any("error", err))
		// If commit fails, the transaction is implicitly rolled back by the database
		// (or should be treated as such).
		successful = false // Mark as unsuccessful if commit fails
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}
