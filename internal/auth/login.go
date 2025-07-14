package auth

import (
	"errors"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"

	"github.com/standard-group/mesa/internal/db"
	"github.com/standard-group/mesa/internal/jwt"
)

func LoginUser(username, serverDomain, password string) (string, error) {
	if username == "" || serverDomain == "" || password == "" {
		log.Warn().Msg("Login failed: missing required fields (username, server domain, or password)")
		return "", errors.New("username, server domain, and password are required")
	}

	user, err := db.GetUserByUsername(username, serverDomain)
	if err != nil {
		log.Warn().Str("username", username).Str("server_domain", serverDomain).Err(err).Msg("User not found or database error during login")
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		log.Warn().Str("username", username).Str("server_domain", serverDomain).Msg("Invalid credentials during password comparison")
		return "", errors.New("invalid credentials")
	}

	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		log.Error().Err(err).Str("username", username).Str("server_domain", serverDomain).Msg("Failed to generate token")
		return "", errors.New("failed to generate authentication token")
	}

	log.Info().Str("username", username).Str("server_domain", serverDomain).Msg("User logged in successfully")
	return token, nil
}
