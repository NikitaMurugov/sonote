package model

import "time"

type Tag struct {
	ID          uint64    `json:"id" db:"id"`
	WorkspaceID uint64    `json:"workspace_id" db:"workspace_id"`
	Name        string    `json:"name" db:"name"`
	Color       *string   `json:"color" db:"color"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type NoteTag struct {
	NoteID    uint64    `json:"note_id" db:"note_id"`
	TagID     uint64    `json:"tag_id" db:"tag_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
