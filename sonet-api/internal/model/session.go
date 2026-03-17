package model

import "time"

type Session struct {
	ID             uint64    `json:"id" db:"id"`
	UserID         uint64    `json:"user_id" db:"user_id"`
	RefreshTokenID *uint64   `json:"refresh_token_id" db:"refresh_token_id"`
	DeviceName     *string   `json:"device_name" db:"device_name"`
	DeviceType     string    `json:"device_type" db:"device_type"`
	OS             *string   `json:"os" db:"os"`
	Browser        *string   `json:"browser" db:"browser"`
	IPAddress      *string   `json:"ip_address" db:"ip_address"`
	Location       *string   `json:"location" db:"location"`
	LastActiveAt   time.Time `json:"last_active_at" db:"last_active_at"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	IsCurrent      bool      `json:"is_current" db:"is_current"`
}
