package auth

import (
	"errors"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"

	"github.com/standard-group/mesa/internal/db"
	"github.com/standard-group/mesa/internal/jwt"
)

func LoginUser(username, password string) (string, error) {
	user, err := db.GetUserByUsername(username)
	if err != nil {
		log.Warn().Str("username", username).Err(err).Msg("User not found")
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		log.Warn().Str("username", username).Msg("Invalid credentials")
		return "", errors.New("invalid credentials")
	}

	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		log.Error().Err(err).Str("username", username).Msg("Failed to generate token")
		return "", err
	}

	log.Info().Str("username", username).Msg("User logged in successfully")
	return token, nil
}
