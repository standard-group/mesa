package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	_ "modernc.org/sqlite"

	"github.com/standard-group/mesa/internal/models"
)

var DB *sql.DB
var localServerDomain string

// config struct for database and server configuration
type config struct {
	Driver       string `toml:"driver"`
	SQLitePath   string `toml:"sqlite_path"`
	PostgresDSN  string `toml:"postgres_dsn"`
	ServerDomain string `toml:"server_domain"`
}

func loadConfig() (*config, error) {
	data, err := os.ReadFile("config/main.toml")
	if err != nil {
		return nil, err
	}
	var cfg config
	err = toml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func InitDB() error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	var driver, dsn string
	switch cfg.Driver {
	case "postgres":
		driver = "postgres"
		dsn = cfg.PostgresDSN
	case "sqlite":
		fallthrough
	default:
		driver = "sqlite"
		dsn = cfg.SQLitePath
	}

	DB, err = sql.Open(driver, dsn)
	if err != nil {
		return err
	}

	localServerDomain = cfg.ServerDomain
	if localServerDomain == "" {
		log.Warn().Msg("ServerDomain not set in config/main.toml. Federation might not work correctly.")
	}

	// psql sjucks and hard i hate it
	stmt := `CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		username TEXT NOT NULL,
		server_domain TEXT NOT NULL,
		password_hash TEXT,
		pubkey TEXT,
		created_at TIMESTAMP WITH TIME ZONE, -- Changed to TIMESTAMP WITH TIME ZONE for PostgreSQL
		UNIQUE(username, server_domain)
	)`
	_, err = DB.Exec(stmt)
	if err != nil {
		return err
	}

	log.Info().Str("driver", driver).Msg("Database initialized")
	return nil
}

func SaveUser(u models.User) error {
	query := `INSERT INTO users (id, username, server_domain, password_hash, pubkey, created_at) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := DB.Exec(query,
		u.ID, u.Username, u.ServerDomain, u.PasswordHash, u.PubKey, u.CreatedAt.Format(time.RFC3339Nano)) // RFC3339Nano text editor whatt!!!
	if err != nil {
		log.Error().Err(err).Str("username", u.Username).Str("server_domain", u.ServerDomain).Msg("Failed to save user to local DB")
	}
	return err
}

func GetUserByUsername(username string, serverDomain string) (models.User, error) {
	var u models.User
	var created string

	// SELECT sigma
	row := DB.QueryRow(`SELECT id, username, server_domain, password_hash, pubkey, created_at FROM users WHERE username = $1 AND server_domain = $2`, username, serverDomain)
	err := row.Scan(&u.ID, &u.Username, &u.ServerDomain, &u.PasswordHash, &u.PubKey, &created)

	if err == nil {
		t, parseErr := time.Parse(time.RFC3339Nano, created)
		if parseErr != nil {
			log.Error().Err(parseErr).Str("created_at_string", created).Msg("Failed to parse created_at timestamp from DB")
			return u, parseErr
		}
		u.CreatedAt = t
		log.Debug().Str("username", username).Str("server_domain", serverDomain).Msg("User found in local DB")
		return u, nil
	}

	if !errors.Is(err, sql.ErrNoRows) {
		log.Error().Err(err).Str("username", username).Str("server_domain", serverDomain).Msg("Database query error for local user")
		return u, err
	}

	if serverDomain != localServerDomain && localServerDomain != "" {
		log.Info().Str("username", username).Str("server_domain", serverDomain).Msg("User not found locally, attempting federation")

		remoteURL := fmt.Sprintf("http://%s/api/v1/users/check?username=%s&server_domain=%s", serverDomain, username, serverDomain)

		resp, httpErr := http.Get(remoteURL)
		if httpErr != nil {
			log.Error().Err(httpErr).Str("remote_url", remoteURL).Msg("Failed to make HTTP request to remote server")
			return u, fmt.Errorf("failed to reach remote server: %w", httpErr)
		}
		defer resp.Body.Close()

		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			log.Error().Err(readErr).Msg("Failed to read response body from remote server")
			return u, fmt.Errorf("failed to read remote response: %w", readErr)
		}

		var remoteUserCheckResponse struct {
			Exists       bool   `json:"exists"`
			UserID       string `json:"user_id"`
			PubKey       string `json:"pub_key"`
			Username     string `json:"username"`
			ServerDomain string `json:"server_domain"`
			Message      string `json:"message"`
		}

		jsonErr := json.Unmarshal(body, &remoteUserCheckResponse)
		if jsonErr != nil {
			log.Error().Err(jsonErr).Str("response_body", string(body)).Msg("Failed to unmarshal JSON response from remote server")
			return u, fmt.Errorf("failed to parse remote response: %w", jsonErr)
		}

		if resp.StatusCode == http.StatusOK && remoteUserCheckResponse.Exists {
			log.Info().Str("username", username).Str("server_domain", serverDomain).Msg("User found on remote server via federation")
			u.ID = remoteUserCheckResponse.UserID
			u.Username = remoteUserCheckResponse.Username
			u.ServerDomain = remoteUserCheckResponse.ServerDomain
			u.PubKey = remoteUserCheckResponse.PubKey
			return u, nil
		} else if resp.StatusCode == http.StatusNotFound {
			log.Info().Str("username", username).Str("server_domain", serverDomain).Msg("User not found on remote server")
			return u, errors.New("user not found")
		} else {
			log.Warn().Str("username", username).Str("server_domain", serverDomain).
				Int("status_code", resp.StatusCode).Str("response_message", remoteUserCheckResponse.Message).
				Msg("Unexpected response from remote server during user check")
			return u, fmt.Errorf("remote server error: %s (status: %d)", remoteUserCheckResponse.Message, resp.StatusCode)
		}
	}

	log.Debug().Str("username", username).Str("server_domain", serverDomain).Msg("User not found locally and no federation attempted/needed")
	return u, errors.New("user not found")
}
