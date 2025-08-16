package config

import (
	"fmt"
	"go-boilerplate/internal/lib"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type LogStyle string

const (
	ColouredLogStyle LogStyle = "colour"
	JSONLogStyle     LogStyle = "json"
	PlainLogStyle    LogStyle = "plain"
)

type AppConfig struct {
	AppName string `env:"APP_NAME" envDefault:"go-boilerplate"`
}

type AuthConfig struct {
	SecretKey                 string        `env:"SECRET_KEY,required"`
	JWTAudience               string        `env:"JWT_AUDIENCE" envDefault:"ethnocopia-app"`
	JWTIssuer                 string        `env:"JWT_ISSUER" envDefault:"ethnocopia-auth"`
	TokenLifetime             time.Duration `env:"TOKEN_LIFETIME" envDefault:"24h"`
	VerificationTokenLifetime time.Duration `env:"VERIFICATION_TOKEN_TTL" envDefault:"30m"`
}

type CacheConfig struct {
	ConnectionConfig RedisConfig `envPrefix:"CACHE_REDIS_"`
}

type DatabaseConfig struct {
	DSN             string        `env:"DATABASE_URL"`
	Driver          string        `env:"DB_DRIVER" envDefault:"postgres"`
	Host            string        `env:"DB_HOST" envDefault:"localhost"`
	Port            int           `env:"DB_PORT" envDefault:"5432"`
	User            string        `env:"DB_USER" envDefault:"user"`
	Password        string        `env:"DB_PASSWORD" envDefault:"password"`
	DBName          string        `env:"DB_NAME" envDefault:"ethnocopia"`
	SSLMode         string        `env:"DB_SSLMODE" envDefault:"disable"`
	TimeZone        string        `env:"DB_TIMEZONE" envDefault:"Asia/Shanghai"`
	MaxOpenConns    int32         `env:"DB_MAX_OPEN_CONNS" envDefault:"25"`
	MaxIdleConns    int32         `env:"DB_MAX_IDLE_CONNS" envDefault:"5"`
	ConnMaxIdleTime time.Duration `env:"DB_CONN_MAX_IDLE_TIME" envDefault:"5m"`
	ConnMaxLifetime time.Duration `env:"DB_CONN_MAX_LIFETIME" envDefault:"1h"`
	EchoSQL         bool          `env:"ECHO_SQL" envDefault:"false"`
}

type LogConfig struct {
	LogLevel string   `env:"LOG_LEVEL" envDefault:"info"`
	Style    LogStyle `env:"LOG_STYLE" envDefault:"json"`
}

type RedisConfig struct {
	Host     string `env:"REDIS_HOST" envDefault:"localhost"`
	Port     int    `env:"REDIS_PORT" envDefault:"6379"`
	DB       int    `env:"REDIS_DB" envDefault:"0"`
	User     string `env:"REDIS_USER" envDefault:""`
	Password string `env:"REDIS_PASSWORD" envDefault:""`
	UseSSL   *bool  `env:"REDIS_USE_SSL"`
}

type RoutingConfig struct {
	BaseURL         string `env:"BASE_URL"`
	FrontendBaseURL string `env:"FRONTEND_BASE_URL"`
}

type ServerConfig struct {
	Host           string   `env:"HOST"`
	Port           int      `env:"PORT"`
	AllowedOrigins []string `env:"ALLOWED_ORIGINS" envSeparator:"," envDefault:""`
	AllowedHeaders []string `env:"ALLOWED_HEADERS" envSeparator:"," envDefault:""`
}

type SMTPConfig struct {
	Host          string        `env:"MAIL_SERVER" envDefault:"localhost"`
	Port          int           `env:"MAIL_PORT" envDefault:"25"`
	Username      string        `env:"MAIL_USERNAME" envDefault:""`
	Password      string        `env:"MAIL_PASSWORD" envDefault:""`
	UseSSL        bool          `env:"MAIL_USE_SSL" envDefault:"false"`
	UseTLS        bool          `env:"MAIL_USE_TLS" envDefault:"false"`
	Timeout       time.Duration `env:"MAIL_TIMEOUT" envDefault:"5s"`
	DefaultSender string        `env:"MAIL_SENDER" envDefault:"noreply@localhost"`
}

type OAuthConfig struct {
	ClientID     string   `env:"CLIENT_ID"`
	ClientSecret string   `env:"CLIENT_SECRET"`
	RedirectURL  string   `env:"REDIRECT_URL"`
	Scopes       []string `env:"SCOPES" envSeparator:","`
}

type JobConfig struct {
	Concurrency      int         `env:"JOB_CONCURRENCY" envDefault:"10"`
	Queues           []string    `env:"JOB_QUEUES" envSeparator:"," envDefault:"critical:6,default:3,low:1"`
	ConnectionConfig RedisConfig `envPrefix:"JOB_REDIS_"`
}

type Settings struct {
	App           AppConfig
	Auth          AuthConfig
	Cache         CacheConfig
	Database      DatabaseConfig
	Email         SMTPConfig
	Job           JobConfig
	Logging       LogConfig
	Redis         RedisConfig
	Routing       RoutingConfig
	Server        ServerConfig `envPrefix:"SERVER_"`
	GoogleOAuth   OAuthConfig  `envPrefix:"GOOGLE_OAUTH_"`
	FacebookOAuth OAuthConfig  `envPrefix:"FACEBOOK_OAUTH_"`
}

var Configuration Settings
var onceLoad sync.Once

func LoadConfig(paths ...string) (err error) {
	onceLoad.Do(func() {
		for _, path := range paths {
			if path == "" {
				continue
			}

			if loadErr := godotenv.Load(path); loadErr != nil && !os.IsNotExist(loadErr) {
				err = loadErr
				return
			}
		}
		err = env.Parse(&Configuration)
	})
	return err
}

func GetSettings(paths ...string) Settings {
	err := LoadConfig(paths...)
	if err != nil {
		panic(fmt.Sprintf("Configuration load error: %s", err))
	}

	return Configuration
}

func (c *DatabaseConfig) BuildDSN() (string, error) {
	if c.DSN != "" {
		return c.DSN, nil
	}

	urlBuilder := lib.NewURLBuilder()
	urlBuilder.WithScheme("postgresql").WithHost(c.Host)
	if c.Port != 0 {
		urlBuilder.WithPort(c.Port)
	}
	if c.User != "" {
		urlBuilder.WithCredentials(c.User, c.Password)
	}
	urlBuilder.WithPath(c.DBName)
	if c.SSLMode != "" {
		urlBuilder.WithQuery("sslmode", c.SSLMode)
	}

	return urlBuilder.Build()
}

func (c *RedisConfig) BuildConnectionString() (string, error) {
	// Validate required fields
	if c.Host == "" {
		return "", fmt.Errorf("REDIS_HOST is required")
	}

	if c.Port <= 0 || c.Port > 65535 {
		return "", fmt.Errorf("REDIS_PORT must be between 1 and 65535")
	}

	if c.DB < 0 {
		return "", fmt.Errorf("REDIS_DB cannot be negative")
	}

	// Determine scheme based on SSL setting
	scheme := "redis"
	if c.UseSSL != nil && *c.UseSSL {
		scheme = "rediss"
	}

	urlBuilder := lib.NewURLBuilder()
	urlBuilder.WithScheme(scheme).WithHost(c.Host).WithPort(c.Port)

	if c.User != "" || c.Password != "" {
		urlBuilder.WithCredentials(c.User, c.Password)
	}

	urlBuilder.WithPath(strconv.Itoa(c.DB))

	return urlBuilder.Build()
}

type RoutingConfigProvider interface {
	GetBaseURL() string
	GetFrontendBaseURL() string
}

func (rc RoutingConfig) GetBaseURL() string {
	return rc.BaseURL
}

func (rc RoutingConfig) GetFrontendBaseURL() string {
	return rc.FrontendBaseURL
}
