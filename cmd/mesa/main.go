package main

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/standard-group/mesa/internal/config"
	"github.com/standard-group/mesa/internal/server"
)

const defaultConfigPath = "config/data/config.toml"

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// If user did not set MESA_CONFIG_PATH, fall back to the generated file.
	if os.Getenv("MESA_CONFIG_PATH") == "" {
		if err := config.Ensure(defaultConfigPath); err != nil {
			log.Fatal().Err(err).Msg("cannot create default config")
		}
		_ = os.Setenv("MESA_CONFIG_PATH", defaultConfigPath)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("invalid configuration")
	}

	zerolog.SetGlobalLevel(cfg.LogLevel())

	// if user did not set PORT, fall back to the config value. else if user did set environment, then use it.
	portStr := os.Getenv("PORT")
	var port int
	if portStr != "" {
		port, err = strconv.Atoi(portStr)
		if err != nil {
			log.Fatal().Err(err).Msg("Invalid PORT environment variable")
		}
	} else {
		port = cfg.HTTP.Port
	}
	srv := server.New(port)

	// trap SIGINT/SIGTERM for graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal().Err(err).Msg("Server failed to start")
		}
	}()

	<-stop
	srv.Shutdown(context.Background())
}
