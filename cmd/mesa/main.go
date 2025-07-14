package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/standard-group/mesa/internal/db"
	"github.com/standard-group/mesa/internal/server"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if err := db.InitDB(); err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize database")
	}
	defer func() {
		if db.DB != nil {
			log.Info().Msg("Closing database connection")
			db.DB.Close()
		}
	}()

	srv := server.NewServer()
	log.Info().Msgf("Starting server on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal().Msg(err.Error())
	}
}
