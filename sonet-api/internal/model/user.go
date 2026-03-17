package model

import "time"

type User struct {
	ID            uint64    `json:"id" db:"id"`
	Email         string    `json:"email" db:"email"`
	Username      string    `json:"username" db:"username"`
	DisplayName   string    `json:"display_name" db:"display_name"`
	PasswordHash  string    `json:"-" db:"password_hash"`
	AvatarURL     *string   `json:"avatar_url" db:"avatar_url"`
	EmailVerified bool      `json:"email_verified" db:"email_verified"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}
