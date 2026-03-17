package model

import "time"

type Folder struct {
	ID          uint64    `json:"id" db:"id"`
	WorkspaceID uint64    `json:"workspace_id" db:"workspace_id"`
	ParentID    *uint64   `json:"parent_id" db:"parent_id"`
	Name        string    `json:"name" db:"name"`
	SortOrder   int       `json:"sort_order" db:"sort_order"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
