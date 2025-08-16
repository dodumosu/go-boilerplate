package api

import (
	"go-boilerplate/internal/cache"
	"go-boilerplate/internal/config"
	"go-boilerplate/internal/db"
	"go-boilerplate/internal/db/repositories"
	"go-boilerplate/internal/job"
	"go-boilerplate/internal/lib"
	"go-boilerplate/internal/redisutil"
	"log/slog"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
)

type APIWrapper struct {
	config        config.Settings
	logger        *slog.Logger
	cacheService  *cache.CacheService
	dbService     *db.DBService
	redisService  *redisutil.RedisService
	accountRepo   *repositories.AccountRepository
	authRepo      *repositories.AuthRepository
	jobService    *job.JobService
	authenticator *lib.JWTAuthenticator
}

func NewAPIWrapper(cfg config.Settings, dbService *db.DBService, cacheService *cache.CacheService, jobService *job.JobService, redisService *redisutil.RedisService, rootLogger *slog.Logger) *APIWrapper {
	accountRepository := repositories.NewAccountRepository(dbService, rootLogger)
	authRepository := repositories.NewAuthRepository(dbService, rootLogger)

	authenticator := lib.NewJWTAuthenticator(cfg.Auth.SecretKey, cfg.Auth.JWTAudience, cfg.Auth.JWTIssuer)

	return &APIWrapper{
		config:        cfg,
		cacheService:  cacheService,
		dbService:     dbService,
		jobService:    jobService,
		logger:        rootLogger,
		redisService:  redisService,
		accountRepo:   accountRepository,
		authRepo:      authRepository,
		authenticator: authenticator,
	}
}

func (w *APIWrapper) registerAuthRoutes(api huma.API) {
	authGroup := huma.NewGroup(api, AccountBase)
	huma.Register(authGroup, huma.Operation{
		Description: "Sign up for an account",
		Method:      http.MethodPost,
		OperationID: "signup",
		Path:        strings.TrimPrefix(AccountSignUp, AccountBase),
		Summary:     "Sign up",
		Tags:        []string{"Accounts"},
	}, w.SignUp)
}

func (w *APIWrapper) registerAPIRoutes(api huma.API) {
	w.registerAuthRoutes(api)
}

func SetupAPI(mux *http.ServeMux, cfg config.Settings, dbService *db.DBService, cacheService *cache.CacheService, jobService *job.JobService, redisService *redisutil.RedisService, rootLogger *slog.Logger) *APIWrapper {
	humaConfig := huma.DefaultConfig("Boilerplate App", "1.0.0")
	humaConfig.DocsPath = APIDocs
	humaConfig.OpenAPI.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"BearerAuth": {
			Type:         "http",
			Scheme:       "bearer",
			BearerFormat: "JWT",
		},
	}

	api := humago.New(mux, humaConfig)
	wrapper := NewAPIWrapper(cfg, dbService, cacheService, jobService, redisService, rootLogger)
	wrapper.registerAPIRoutes(api)

	return wrapper
}
