package server

import (
	"context"
	"go-boilerplate/internal/cache"
	"go-boilerplate/internal/config"
	"go-boilerplate/internal/contracts"
	"go-boilerplate/internal/db"
	"go-boilerplate/internal/job"
	jobhandlers "go-boilerplate/internal/job/handlers"
	"go-boilerplate/internal/mailer"
	"go-boilerplate/internal/redisutil"
	"go-boilerplate/internal/server/api"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
)

type Server struct {
	Cache       *cache.CacheService
	Job         *job.JobService
	Redis       *redisutil.RedisService
	logger      *slog.Logger
	mux         *http.ServeMux
	settings    config.Settings
	cacheClient *redis.Client
	jobClient   *redis.Client
	redisClient *redis.Client
	DB          *db.DBService
}

func NewServer(logger *slog.Logger, cfg config.Settings) (*Server, error) {
	cacheConfig := cfg.Cache
	jobConfig := cfg.Job
	redisConfig := cfg.Redis

	cacheClient, err := redisutil.NewRedisClient(cacheConfig.ConnectionConfig)
	if err != nil {
		return nil, err
	}

	jobClient, err := redisutil.NewRedisClient(jobConfig.ConnectionConfig)
	if err != nil {
		return nil, err
	}

	redisClient, err := redisutil.NewRedisClient(redisConfig)
	if err != nil {
		return nil, err
	}

	cacheService, err := cache.NewCacheService(logger, cacheClient)
	if err != nil {
		return nil, err
	}

	emailService, err := mailer.NewEmailService(cfg, logger)
	if err != nil {
		return nil, err
	}

	handlerCtx := &job.JobHandlerContext{
		EmailService: emailService,
	}

	jobService, err := job.NewJobService(logger, &jobConfig, jobClient, handlerCtx)
	if err != nil {
		return nil, err
	}

	jobService.RegisterHandler(contracts.TypeVerificationEmail, jobhandlers.HandleVerificationEmail)

	redisService, err := redisutil.NewRedisService(logger, redisClient)
	if err != nil {
		return nil, err
	}

	dbService, err := db.NewDBService(context.Background(), &cfg.Database, logger)
	if err != nil {
		return nil, err
	}

	// handlers.SetupGoth(&cfg)
	mux := http.NewServeMux()

	return &Server{
		Cache:       cacheService,
		Job:         jobService,
		Redis:       redisService,
		cacheClient: cacheClient,
		jobClient:   jobClient,
		redisClient: redisClient,
		logger:      logger,
		mux:         mux,
		settings:    cfg,
		DB:          dbService,
	}, nil
}

func (s *Server) setupRoutes() http.Handler {
	// s.mux.HandleFunc("/auth/std/login/google", httphandlers.GoogleLoginHandler(&s.settings))
	// s.mux.HandleFunc("/auth/std/callback/google", httphandlers.GoogleCallbackHandler(&s.settings))
	// s.mux.HandleFunc("/auth/std/login/facebook", httphandlers.FacebookLoginHandler(&s.settings))
	// s.mux.HandleFunc("/auth/std/callback/facebook", httphandlers.FacebookCallbackHandler(&s.settings))

	// s.mux.HandleFunc("/auth/goth/login/google", handlers.GothGoogleLoginHandler)
	// s.mux.HandleFunc("/auth/goth/callback/google", handlers.GothGoogleCallbackHandler)
	// s.mux.HandleFunc("/auth/goth/login/facebook", handlers.GothFacebookLoginHandler)
	// s.mux.HandleFunc("/auth/goth/callback/facebook", handlers.GothFacebookCallbackHandler)
	wrapper := api.SetupAPI(s.mux, s.settings, s.DB, s.Cache, s.Job, s.Redis, s.logger)

	return wrapper.SetupMiddleware(s.mux)
}

func (s *Server) cleanup() {
	s.logger.Info("Cleaning up resources...")
	s.Job.Stop()
	s.cacheClient.Close()
	s.jobClient.Close()
	s.redisClient.Close()
	s.DB.Close()
}

func (s *Server) Start() {
	address := net.JoinHostPort(s.settings.Server.Host, strconv.Itoa(s.settings.Server.Port))
	handler := s.setupRoutes()
	server := &http.Server{
		Addr:    address,
		Handler: handler,
	}

	s.Job.Start()
	serverErrors := make(chan error, 1)

	go func() {
		s.logger.Info("Starting HTTP server", "address", address)
		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		s.logger.Error("Server error", "error", err)
		os.Exit(1)
	case sig := <-shutdown:
		s.logger.Info("Received signal", "signal", sig, "message", "Starting graceful shutdown...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			s.logger.Error("Server shutdown failed", "error", err)
		} else {
			s.logger.Info("Server gracefully stopped")
		}

		s.cleanup()
	}
}
