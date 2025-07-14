package server

import (
	"encoding/json"
	"net/http"

	"errors"

	"github.com/rs/zerolog/log"
	"github.com/standard-group/mesa/internal/auth"
	"github.com/standard-group/mesa/internal/db"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username     string `json:"username"`
		ServerDomain string `json:"server_domain"`
		Password     string `json:"password"`
		PubKey       string `json:"pubkey"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Warn().Err(err).Msg("RegisterHandler: Invalid request body")
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := auth.RegisterUser(req.Username, req.ServerDomain, req.Password, req.PubKey); err != nil {
		log.Warn().Err(err).Str("username", req.Username).Str("server_domain", req.ServerDomain).Msg("Registration failed")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username     string `json:"username"`
		ServerDomain string `json:"server_domain"`
		Password     string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Warn().Err(err).Msg("LoginHandler: Invalid request body")
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Now passing username and server_domain separately to auth.LoginUser
	token, err := auth.LoginUser(req.Username, req.ServerDomain, req.Password)
	if err != nil {
		log.Warn().Err(err).Str("username", req.Username).Str("server_domain", req.ServerDomain).Msg("Login failed")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// usercheckhandler handles requests from other servers to check for a user's existence яlocally.
// expected request: GET /api/v1/users/check?username={username}&server_domain={server_domain}
func UserCheckHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	serverDomain := r.URL.Query().Get("server_domain")

	if username == "" || serverDomain == "" {
		log.Warn().Msg("UserCheckHandler: Missing username or server_domain query parameters")
		http.Error(w, "Missing username or server_domain", http.StatusBadRequest)
		return
	}

	// Check if the user exists in the local database
	user, err := db.GetUserByUsername(username, serverDomain)
	if err != nil {
		if errors.Is(err, errors.New("user not found")) { // Check for the specific "user not found" error
			log.Debug().Str("username", username).Str("server_domain", serverDomain).Msg("User not found locally")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]interface{}{"exists": false, "message": "User not found locally"})
			return
		}
		log.Error().Err(err).Str("username", username).Str("server_domain", serverDomain).Msg("Database error during user check")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// User found locally, return success
	log.Info().Str("username", username).Str("server_domain", serverDomain).Msg("User found locally via check handler")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"exists":        true,
		"user_id":       user.ID,
		"pub_key":       user.PubKey, // Include pubkey for federation if needed
		"username":      user.Username,
		"server_domain": user.ServerDomain,
	})
}
