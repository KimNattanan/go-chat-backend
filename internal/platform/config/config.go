package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type (
	// Config -.
	Config struct {
		App     App
		HTTP    HTTP
		Log     Log
		DB      DB
		Redis   Redis
		GRPC    GRPC
		RMQ     RMQ
		Swagger Swagger
		JWT     JWT
	}

	// App -.
	App struct {
		ENV     string `env:"APP_ENV,required"`
		Name    string `env:"APP_NAME,required"`
		Version string `env:"APP_VERSION,required"`
	}

	// HTTP -.
	HTTP struct {
		Port           string `env:"HTTP_PORT,required"`
		AllowedOrigins string `env:"HTTP_ALLOWED_ORIGINS" envDefault:"http://localhost:3000"`
	}

	// Log -.
	Log struct {
		Level string `env:"LOG_LEVEL,required"`
	}

	// DB -.
	DB struct {
		Host     string `env:"DB_HOST,required"`
		Port     string `env:"DB_PORT,required"`
		User     string `env:"DB_USER,required"`
		Password string `env:"DB_PASSWORD,required"`
		Name     string `env:"DB_NAME,required"`
		DSN      string
	}

	// Redis -.
	Redis struct {
		Address  string `env:"REDIS_ADDRESS,required"`
		Password string `env:"REDIS_PASSWORD,required"`
		DB       int    `env:"REDIS_DB,required"`
	}

	// GRPC -.
	GRPC struct {
		Host string `env:"GRPC_HOST,required"`
		Port string `env:"GRPC_PORT,required"`
	}

	// RMQ -.
	RMQ struct {
		ServerExchange string `env:"RMQ_RPC_SERVER,required"`
		ClientExchange string `env:"RMQ_RPC_CLIENT,required"`
		URL            string `env:"RMQ_URL,required"`
	}

	// Swagger -.
	Swagger struct {
		Enabled bool `env:"SWAGGER_ENABLED" envDefault:"false"`
	}

	// JWT -.
	JWT struct {
		Secret     string `env:"JWT_SECRET,required"`
		AccessTTL  int    `env:"ACCESS_TTL" envDefault:"3600"`    // 1 hour, in seconds
		RefreshTTL int    `env:"REFRESH_TTL" envDefault:"604800"` // 7 days
	}
)

// NewConfig returns app config.
func NewConfig(appEnv string) (*Config, error) {
	envFile := ".env"
	if appEnv != "" {
		envFile = ".env." + appEnv
	}
	if err := godotenv.Load(envFile); err != nil {
		log.Printf("Warning: could not load .env file: %v", err)
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}
	cfg.DB.DSN = fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name,
	)

	return cfg, nil
}
