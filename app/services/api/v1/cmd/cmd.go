package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	db "github.com/hpetrov29/resttemplate/business/data/dbsql/mysql"
	v1 "github.com/hpetrov29/resttemplate/business/web/v1"
	"github.com/hpetrov29/resttemplate/business/web/v1/auth"
	"github.com/hpetrov29/resttemplate/internal/keystore"
	"github.com/hpetrov29/resttemplate/internal/logger"
	"github.com/hpetrov29/resttemplate/internal/web"
	"github.com/rs/cors"
	"github.com/sethvargo/go-envconfig"
)

func Main(routeAdder v1.RouteAdder) error {
	var log *logger.Logger

	events := logger.Events{
		Error: func(ctx context.Context, r logger.Record) {
			log.Info(ctx, "******* SEND ALERT ******")
		},
	}

	traceIDFunc := func(ctx context.Context) string {
		return web.GetTraceID(ctx)
	}

	// Logger will disregard logs of category lower than the one specified here
	log = logger.NewWithEvents(os.Stdout, logger.LevelInfo, "API", traceIDFunc, events)

	ctx := context.Background()

	if err := run(ctx, log, "v1", routeAdder); err != nil {
		log.Error(ctx, "startup", "msg", err)
		return err
	}
	return nil
}

func run(ctx context.Context, log *logger.Logger, build string, routeAdder v1.RouteAdder) error {

	// -------------------------------------------------------------------------
	// GOMAXPROCS
	
	log.Info(ctx, "service startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// -------------------------------------------------------------------------
	// Configuration

	var config GlobalConfig
	config.Version.Build = build
	if err := envconfig.Process(ctx, &config); err != nil {
		return fmt.Errorf("error while parsing env variables/config: %w", err)
	}
	
	// -------------------------------------------------------------------------
	// Set up database client conneciton

	log.Info(ctx, "DB startup", "status", "initializing database support", "host", config.DB.Host)

	mysqlClient, err := db.Open(db.Config{
		User:         config.DB.User,
		Password:     config.DB.Password,
		Host:         config.DB.Host,
		Name:         config.DB.Name,
		MaxIdleConns: config.DB.MaxIdleConns,
		MaxOpenConns: config.DB.MaxOpenConns,
		DisableTLS:   config.DB.DisableTLS,
	})
	if err != nil {
		return fmt.Errorf( "error connecting to db: %w", err)
	}
	defer func() {
		log.Info(ctx, "DB shutdown", "status", "stopping database support", "host", config.DB.Host)
		mysqlClient.Close()
	}()
	err = db.StatusCheck(ctx, mysqlClient); if err != nil {
		return fmt.Errorf("error database status check: %w", err)
	}

	// -------------------------------------------------------------------------
	// Initialize authentication support

	log.Info(ctx, "Auth startup", "status", "initializing authentication support")

	keystore, err := keystore.NewFS(os.DirFS(config.Auth.KeysFolder))
	if err != nil {
		return fmt.Errorf("error while retrieving keys: %w", err)
	}

	auth, err := auth.New(auth.Config{
		Log:       log,
		DB:        mysqlClient,
		Issuer:    config.Auth.Issuer,
		Vault: 	   keystore,
	})
	if err != nil {
		return fmt.Errorf("error constructing Auth service: %w", err)
	}

	// -------------------------------------------------------------------------
	// Start API

	log.Info(ctx, "API startup", "version", build)

	// Only the signals explicitly provided (SIGINT and SIGTERM) will be captured.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	muxConfig := v1.APIMuxConfig{
		Build: build,
		Shutdown: shutdown,
		Log: log,
		Auth: auth,
		DB: mysqlClient,
	}

	apiMux := v1.NewAPIMux(muxConfig, routeAdder)

	c := cors.New(cors.Options{
        AllowedOrigins:   config.CORS.AllowedOrigins,
        AllowCredentials: true,
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Authorization", "Content-Type", "Set-Cookie"},
    })

	api := &http.Server{
        Addr:    config.Web.APIHost,
        Handler: c.Handler(apiMux),
		ReadTimeout:  config.Web.ReadTimeout,
		WriteTimeout: config.Web.WriteTimeout,
		IdleTimeout:  config.Web.IdleTimeout,
    }

	serverErrors := make(chan error, 1)

	go func() {
		serverErrors <- api.ListenAndServe()
	}()

	select{
		case err := <-serverErrors:
			return fmt.Errorf("server error: %w", err)
		case sig := <-shutdown:
			log.Info(ctx, "API shutdown", "status", "shutdown started", "signal", sig)
			defer log.Info(ctx, "API shutdown", "status", "shutdown complete", "signal", sig)
	
			ctx, cancel := context.WithTimeout(ctx, config.Web.ShutdownTimeout)
			defer cancel()
	
			if err := api.Shutdown(ctx); err != nil {
				api.Close()
				return fmt.Errorf("could not stop server gracefully: %w", err)
			}
	}

	return nil
}