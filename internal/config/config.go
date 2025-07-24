package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/rs/zerolog"
)

// Config is the top-level configuration object.
type Config struct {
	HTTP  HTTPConfig  `toml:"http"`
	Log   LogConfig   `toml:"log"`
	Paths PathsConfig `toml:"paths"`
}

type HTTPConfig struct {
	Port int `toml:"port"`
}

type LogConfig struct {
	Level string `toml:"level"`
}

type PathsConfig struct {
	DataDir string `toml:"data_dir"`
}

// Load reads the configuration from:
//  1. TOML file pointed to by MESA_CONFIG_PATH (optional).
//  2. Environment variables (take precedence).
func Load() (*Config, error) {
	cfg := &Config{
		HTTP:  HTTPConfig{Port: 8080},
		Log:   LogConfig{Level: "info"},
		Paths: PathsConfig{DataDir: "./data"},
	}

	if path := os.Getenv("MESA_CONFIG_PATH"); path != "" {
		if _, err := toml.DecodeFile(path, cfg); err != nil {
			return nil, fmt.Errorf("decode %s: %w", path, err)
		}
	}

	overrideFromEnv(cfg)

	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

// overrideFromEnv applies env vars on top of any TOML values.
func overrideFromEnv(c *Config) {
	if p := os.Getenv("PORT"); p != "" {
		if v, err := strconv.Atoi(p); err == nil {
			c.HTTP.Port = v
		}
	}
	if l := os.Getenv("LOG_LEVEL"); l != "" {
		c.Log.Level = strings.ToLower(l)
	}
	if d := os.Getenv("MESA_DATA_DIR"); d != "" {
		c.Paths.DataDir = d
	}
}

// validate checks that all values are sane.
func (c *Config) validate() error {
	if c.HTTP.Port <= 0 || c.HTTP.Port > 65535 {
		return fmt.Errorf("invalid port: %d", c.HTTP.Port)
	}
	switch c.Log.Level {
	case "trace", "debug", "info", "warn", "error":
	default:
		return fmt.Errorf("invalid log level: %s", c.Log.Level)
	}
	if c.Paths.DataDir == "" {
		return fmt.Errorf("data_dir is required")
	}
	if abs, err := filepath.Abs(c.Paths.DataDir); err != nil {
		return fmt.Errorf("data_dir: %w", err)
	} else {
		c.Paths.DataDir = abs
	}
	return nil
}

// LogLevel returns the zerolog.Level corresponding to the config value.
func (c *Config) LogLevel() zerolog.Level {
	switch c.Log.Level {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	default:
		return zerolog.InfoLevel
	}
}

// Ensure creates a default config file at the given path if none exists.
// It also makes parent directories as needed.
func Ensure(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil // already exists
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(defaultConfig), 0o644)
}

// defaultConfig is embedded into the binary.
const defaultConfig = `
# Mesa server configuration.
# Edit it and save it as config.toml.

[http]
port = 8080

[log]
level = "info"

[paths]
data_dir = "./data"
`
