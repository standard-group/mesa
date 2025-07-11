package db

import (
	"database/sql"
	"errors"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	_ "modernc.org/sqlite"

	"github.com/standard-group/mesa/internal/models"
)

var DB *sql.DB

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

	stmt := `CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		username TEXT UNIQUE,
		password_hash TEXT,
		pubkey TEXT,
		created_at TEXT
	)`
	_, err = DB.Exec(stmt)
	if err != nil {
		return err
	}

	log.Info().Str("driver", driver).Msg("Database initialized")
	return nil
}

func SaveUser(u models.User) error {
	_, err := DB.Exec(`INSERT INTO users (id, username, password_hash, pubkey, created_at) VALUES (?, ?, ?, ?, ?)`,
		u.ID, u.Username, u.PasswordHash, u.PubKey, u.CreatedAt.Format("2006-01-02T15:04:05Z"))
	return err
}

func GetUserByUsername(username string) (models.User, error) {
	row := DB.QueryRow(`SELECT id, username, password_hash, pubkey, created_at FROM users WHERE username = ?`, username)
	var u models.User
	var created string
	if err := row.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.PubKey, &created); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return u, errors.New("user not found")
		}
		return u, err
	}
	t, err := time.Parse("2006-01-02T15:04:05Z", created)
	if err != nil {
		return u, err
	}
	u.CreatedAt = t
	return u, nil
}

type config struct {
	Driver      string `toml:"driver"`
	SQLitePath  string `toml:"sqlite_path"`
	PostgresDSN string `toml:"postgres_dsn"`
}

func loadConfig() (*config, error) {
	data, err := os.ReadFile("config/database.toml")
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
