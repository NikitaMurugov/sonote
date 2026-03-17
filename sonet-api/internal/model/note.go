package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// NullRawMessage handles NULL JSON columns from the database.
type NullRawMessage struct {
	json.RawMessage
	Valid bool
}

func (n *NullRawMessage) Scan(value interface{}) error {
	if value == nil {
		n.RawMessage = nil
		n.Valid = false
		return nil
	}
	n.Valid = true
	switch v := value.(type) {
	case []byte:
		n.RawMessage = make(json.RawMessage, len(v))
		copy(n.RawMessage, v)
	case string:
		n.RawMessage = json.RawMessage(v)
	}
	return nil
}

func (n NullRawMessage) Value() (driver.Value, error) {
	if !n.Valid || n.RawMessage == nil {
		return nil, nil
	}
	return []byte(n.RawMessage), nil
}

func (n NullRawMessage) MarshalJSON() ([]byte, error) {
	if !n.Valid || n.RawMessage == nil {
		return []byte("null"), nil
	}
	return n.RawMessage, nil
}

func (n *NullRawMessage) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Valid = false
		n.RawMessage = nil
		return nil
	}
	n.Valid = true
	n.RawMessage = make(json.RawMessage, len(data))
	copy(n.RawMessage, data)
	return nil
}

type Note struct {
	ID          uint64          `json:"id" db:"id"`
	WorkspaceID uint64          `json:"workspace_id" db:"workspace_id"`
	FolderID    *uint64         `json:"folder_id" db:"folder_id"`
	Title       string          `json:"title" db:"title"`
	Slug        string          `json:"slug" db:"slug"`
	ContentMD   string          `json:"content_md" db:"content_md"`
	ContentHTML string          `json:"content_html" db:"content_html"`
	ContentJSON NullRawMessage  `json:"content_json" db:"content_json"`
	AuthorID    uint64          `json:"author_id" db:"author_id"`
	IsPinned    bool            `json:"is_pinned" db:"is_pinned"`
	IsArchived  bool            `json:"is_archived" db:"is_archived"`
	WordCount   uint32          `json:"word_count" db:"word_count"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
}
