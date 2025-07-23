// Example main.go (ensure it's flexible with PORT and MESA_CONFIG_PATH)
package main

import (
	// "fmt"
	// "net/http"
	"os"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "8080"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatal().Err(err).Msg("Invalid PORT environment variable")
	}

	// srv := server.NewServer()
	// srv.Addr = fmt.Sprintf(":%d", port)

	log.Info().Int("port", port).Msg("Starting server")
	// if err := srv.ListenAndServe(); err != http.ErrServerClosed {
	//	log.Fatal().Err(err).Msg("Server failed to start")
	// }
}
