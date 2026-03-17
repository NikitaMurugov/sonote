package model

import "time"

type Workspace struct {
	ID          uint64    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Slug        string    `json:"slug" db:"slug"`
	Description *string   `json:"description" db:"description"`
	OwnerID     uint64    `json:"owner_id" db:"owner_id"`
	IsPersonal  bool      `json:"is_personal" db:"is_personal"`
	Icon        *string   `json:"icon" db:"icon"`
	IsEncrypted bool      `json:"is_encrypted" db:"is_encrypted"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type WorkspaceMember struct {
	ID           uint64    `json:"id" db:"id"`
	WorkspaceID  uint64    `json:"workspace_id" db:"workspace_id"`
	UserID       uint64    `json:"user_id" db:"user_id"`
	Role         string    `json:"role" db:"role"`
	InvitedBy    *uint64   `json:"invited_by" db:"invited_by"`
	EncryptedDEK *string   `json:"encrypted_dek,omitempty" db:"encrypted_dek"`
	JoinedAt     time.Time `json:"joined_at" db:"joined_at"`
}
