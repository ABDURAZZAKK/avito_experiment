package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ABDURAZZAKK/avito_experiment/config"
	"github.com/ABDURAZZAKK/avito_experiment/internal/repo"
	"github.com/ABDURAZZAKK/avito_experiment/internal/service"

	v1 "github.com/ABDURAZZAKK/avito_experiment/internal/controller/http/v1"
	"github.com/ABDURAZZAKK/avito_experiment/pkg/broker"
	"github.com/ABDURAZZAKK/avito_experiment/pkg/httpserver"
	"github.com/ABDURAZZAKK/avito_experiment/pkg/postgres"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func Run(configPath string) {
	// Configuration
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Logger
	SetLogrus(cfg.Log.Level)

	// Repositories
	log.Info("Initializing postgres...")
	// pg, err := postgres.New("postgres://postgres:postgres@localhost:5433/postgres", postgres.MaxPoolSize(cfg.PG.MaxPoolSize))
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.MaxPoolSize))
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - pgdb.NewServices: %w", err))
	}
	defer pg.Close()

	// Repositories
	log.Info("Initializing repositories...")
	repositories := repo.NewRepositories(pg)

	// Services dependencies
	log.Info("Initializing services...")
	deps := service.ServicesDependencies{
		Repos: repositories,
	}
	services := service.NewServices(deps)

	// Echo handler
	log.Info("Initializing handlers routes and RabbitMQ...")
	handler := echo.New()
	// rabbit, err := broker.NewRabbitMQ("amqp://guest:guest@localhost:5672")
	rabbit, err := broker.NewRabbitMQ(cfg.BROKER.URL)
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - NewRabbitMQ: %w", err))
	}
	defer rabbit.Close()

	// setup handler validator as lib validator
	v1.NewRouter(handler, services, rabbit)

	// HTTP server
	log.Info("Starting http server...")
	log.Debugf("Server port: %s", cfg.HTTP.Port)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	log.Info("Configuring graceful shutdown...")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		log.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Graceful shutdown
	log.Info("Shutting down...")
	err = httpServer.Shutdown()
	if err != nil {
		log.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
