package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/NikitaMurugov/sonote-api/internal/model"
)

type NoteLinkRepository struct {
	db *sqlx.DB
}

func NewNoteLinkRepository(db *sqlx.DB) *NoteLinkRepository {
	return &NoteLinkRepository{db: db}
}

func (r *NoteLinkRepository) Upsert(ctx context.Context, link *model.NoteLink) error {
	query := `INSERT INTO note_links (source_note_id, target_note_id, context_snippet)
	          VALUES (?, ?, ?)
	          ON DUPLICATE KEY UPDATE context_snippet = VALUES(context_snippet)`
	_, err := r.db.ExecContext(ctx, query, link.SourceNoteID, link.TargetNoteID, link.ContextSnippet)
	return err
}

func (r *NoteLinkRepository) DeleteBySource(ctx context.Context, sourceNoteID uint64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM note_links WHERE source_note_id = ?", sourceNoteID)
	return err
}

func (r *NoteLinkRepository) GetOutgoing(ctx context.Context, noteID uint64) ([]model.NoteLink, error) {
	var links []model.NoteLink
	err := r.db.SelectContext(ctx, &links,
		"SELECT * FROM note_links WHERE source_note_id = ?", noteID)
	return links, err
}

func (r *NoteLinkRepository) GetBacklinks(ctx context.Context, noteID uint64) ([]model.NoteLink, error) {
	var links []model.NoteLink
	err := r.db.SelectContext(ctx, &links,
		"SELECT * FROM note_links WHERE target_note_id = ?", noteID)
	return links, err
}

func (r *NoteLinkRepository) GetAllForWorkspace(ctx context.Context, workspaceID uint64) ([]model.NoteLink, error) {
	var links []model.NoteLink
	err := r.db.SelectContext(ctx, &links,
		`SELECT nl.* FROM note_links nl
		 JOIN notes n ON nl.source_note_id = n.id
		 WHERE n.workspace_id = ?`, workspaceID)
	return links, err
}
