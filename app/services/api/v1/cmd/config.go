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
	SQLDB struct {
		User         string `env:"SQLDB_USER, required"`
		Password     string `env:"SQLDB_PASSWORD, required"`
		Host         string `env:"SQLDB_HOST, required"`
		Name         string `env:"SQLDB_NAME, required"`
		MaxIdleConns int    `env:"SQLDB_MAX_IDLE_CONNECTIONS, default=2"`
		MaxOpenConns int    `env:"SQLDB_MAX_OPEN_CONNECTIONS, default=0"`
		DisableTLS   bool   `env:"SQLDB_DISABLE_TLS, default=true"`
	}
	NOSQLDB struct {
		User         string `env:"NOSQLDB_USER, required"`
		Password     string `env:"NOSQLDB_PASSWORD, required"`
		Host         string `env:"NOSQLDB_HOST, required"`
		Name         string `env:"NOSQLDB_NAME, required"`
		MaxOpenConns int    `env:"NOSQLDB_MAX_OPEN_CONNECTIONS, default=0"`
	}
	Messaging struct {
		User         string `env:"MESSAGING_USER, required"`
		Password     string `env:"MESSAGING_PASSWORD, required"`
		Host         string `env:"MESSAGING_HOST"`
	}
	Auth struct {
		KeysFolder string `env:"KEY_PATH, default=./zarf/keys/"`
		Issuer     string `env:"ISSUER_NAME, default=service"`
	}
	CORS struct {
		AllowedOrigins []string `env:"ALLOWED_ORIGINS, delimiter=;, required"`
	}
}