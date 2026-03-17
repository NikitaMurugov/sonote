package model

import "time"

type User struct {
	ID                  uint64    `json:"id" db:"id"`
	Email               string    `json:"email" db:"email"`
	Username            string    `json:"username" db:"username"`
	DisplayName         string    `json:"display_name" db:"display_name"`
	PasswordHash        string    `json:"-" db:"password_hash"`
	AvatarURL           *string   `json:"avatar_url" db:"avatar_url"`
	EmailVerified       bool      `json:"email_verified" db:"email_verified"`
	UserSalt            *string   `json:"user_salt,omitempty" db:"user_salt"`
	EncryptedPrivateKey *string   `json:"encrypted_private_key,omitempty" db:"encrypted_private_key"`
	PublicKey           *string   `json:"public_key,omitempty" db:"public_key"`
	RecoveryDEK         *string   `json:"-" db:"recovery_dek"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}
