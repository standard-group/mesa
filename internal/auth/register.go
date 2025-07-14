package auth

import (
	"errors"
	"time"

	"github.com/standard-group/mesa/internal/db"

	"github.com/standard-group/mesa/internal/models"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(username, serverdomain, password, pubkey string) error {
	if username == "" || serverdomain == "" || password == "" || pubkey == "" {
		log.Warn().Msg("Registration failed: missing required fields")
		return errors.New("username, serverdomain, password, and pubkey are required")
	}

	// TODO: add validation

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("Failed to hash password")
		return err
	}

	user := models.User{
		ID:           uuid.New().String(),
		Username:     username,
		ServerDomain: serverdomain,
		PasswordHash: string(hash),
		PubKey:       pubkey,
		CreatedAt:    time.Now(),
	}

	err = db.SaveUser(user)
	if err != nil {
		log.Error().Err(err).Str("username", username).Str("server_domain", serverdomain).Msg("Failed to save user")
		return err
	}

	log.Info().Str("username", username).Str("server_domain", serverdomain).Msg("User registered successfully")
	return nil
}
