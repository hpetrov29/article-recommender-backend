package cmd

import "time"

type GlobalConfig struct {
	Version struct {
		Build       string
		Description string
	}
	Web struct {
		APIHost         string        `env:"PORT, required"`
		ReadTimeout     time.Duration `env:"READ_TIMEOUT, default=5s"`
		WriteTimeout    time.Duration `env:"WRITE_TIMEOUT, default=10s"`
		IdleTimeout     time.Duration `env:"IDLE_TIMEOUT, default=120s"`
		ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT, default=20s"`
		//DebugHost       string          `conf:"default:0.0.0.0:4000"`
	}
	DB struct {
		User         string `env:"DB_USER, required"`
		Password     string `env:"DB_PASSWORD, required"`
		Host         string `env:"DB_HOST, required"`
		Name         string `env:"DB_NAME, required"`
		MaxIdleConns int    `env:"DB_MAX_IDLE_CONNECTIONS, default=2"`
		MaxOpenConns int    `env:"DB_MAX_OPEN_CONNECTIONS, default=0"`
		DisableTLS   bool   `env:"DB_DISABLE_TLS, default=true"`
	}
	Auth struct {
		KeysFolder string `env:"KEY_PATH, default=./zarf/keys/"`
		Issuer     string `env:"ISSUER_NAME, default=service"`
	}
	CORS struct {
		AllowedOrigins []string `env:"ALLOWED_ORIGINS, delimiter=;, required"`
	}
}