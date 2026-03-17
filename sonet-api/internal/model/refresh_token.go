package model

import "time"

type RefreshToken struct {
	ID         uint64    `json:"id" db:"id"`
	UserID     uint64    `json:"user_id" db:"user_id"`
	TokenHash  string    `json:"-" db:"token_hash"`
	DeviceInfo *string   `json:"device_info" db:"device_info"`
	IPAddress  *string   `json:"ip_address" db:"ip_address"`
	ExpiresAt  time.Time `json:"expires_at" db:"expires_at"`
	Revoked    bool      `json:"revoked" db:"revoked"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}
