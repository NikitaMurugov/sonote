package model

import "time"

type NoteLink struct {
	ID             uint64    `json:"id" db:"id"`
	SourceNoteID   uint64    `json:"source_note_id" db:"source_note_id"`
	TargetNoteID   uint64    `json:"target_note_id" db:"target_note_id"`
	ContextSnippet *string   `json:"context_snippet" db:"context_snippet"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}
