package models

import "time"

type User struct {
	ID           string    `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	ServerDomain string    `json:"server_domain" db:"server_domain"`
	PasswordHash string    `json:"-" db:"password_hash"`
	PubKey       string    `json:"pub_key" db:"pub_key"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
