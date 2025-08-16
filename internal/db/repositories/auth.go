package repositories

import (
	"context"
	"errors"
	"go-boilerplate/internal/db"
	"go-boilerplate/internal/lib"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type AuthRepository struct {
	dbService *db.DBService
	logger    *slog.Logger
}

func NewAuthRepository(dbService *db.DBService, parentLogger *slog.Logger) *AuthRepository {
	return &AuthRepository{
		dbService: dbService,
		logger:    parentLogger.With("repository", "auth"),
	}
}

func (r *AuthRepository) SigninWithPassword(ctx context.Context, loginIdentifier string, password string) (*db.LookupAccountForAuthRow, error) {
	accountInfo, err := r.dbService.LookupAccountForAuth(ctx, pgtype.Text{String: loginIdentifier, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAccountDoesNotExist
		}

		r.logger.Error("error looking up account for password auth", "error", err)
		return nil, err
	}

	currentTimestamp := time.Now()
	if !accountInfo.IsActive {
		return nil, ErrAccountIsNotActive
	}

	if !accountInfo.IsVerified {
		return nil, ErrAccountIsNotVerified
	}

	if accountInfo.BannedAt.Valid && accountInfo.BannedAt.Time.Before(currentTimestamp) {
		return nil, ErrAccountIsBanned
	}

	if accountInfo.SuspendedUntil.Valid && accountInfo.SuspendedUntil.Time.After(currentTimestamp) {
		return nil, ErrAccountIsSuspended
	}

	if !accountInfo.PasswordHash.Valid {
		return nil, ErrInvalidCredentials
	}

	passwordMatch, err := lib.VerifyPassword(password, accountInfo.PasswordHash.String)
	if err != nil {
		return nil, err
	}

	if !passwordMatch {
		return nil, ErrInvalidCredentials
	}

	return &accountInfo, nil
}
