package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/NikitaMurugov/sonote-api/internal/model"
)

type NoteRepository struct {
	db *sqlx.DB
}

func NewNoteRepository(db *sqlx.DB) *NoteRepository {
	return &NoteRepository{db: db}
}

func (r *NoteRepository) Create(ctx context.Context, n *model.Note) error {
	query := `INSERT INTO notes (workspace_id, folder_id, title, slug, content_md, content_html, content_json, author_id, word_count)
	           VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query,
		n.WorkspaceID, n.FolderID, n.Title, n.Slug,
		n.ContentMD, n.ContentHTML, n.ContentJSON, n.AuthorID, n.WordCount)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	n.ID = uint64(id)
	return nil
}

func (r *NoteRepository) GetByID(ctx context.Context, id uint64) (*model.Note, error) {
	var n model.Note
	err := r.db.GetContext(ctx, &n, "SELECT * FROM notes WHERE id = ?", id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &n, err
}

func (r *NoteRepository) GetBySlug(ctx context.Context, workspaceID uint64, slug string) (*model.Note, error) {
	var n model.Note
	err := r.db.GetContext(ctx, &n,
		"SELECT * FROM notes WHERE workspace_id = ? AND slug = ?", workspaceID, slug)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &n, err
}

func (r *NoteRepository) ListByWorkspace(ctx context.Context, workspaceID uint64, folderID *uint64, archived bool) ([]model.Note, error) {
	var notes []model.Note
	if folderID != nil {
		err := r.db.SelectContext(ctx, &notes,
			`SELECT * FROM notes WHERE workspace_id = ? AND folder_id = ? AND is_archived = ?
			 ORDER BY is_pinned DESC, updated_at DESC`, workspaceID, *folderID, archived)
		return notes, err
	}
	err := r.db.SelectContext(ctx, &notes,
		`SELECT * FROM notes WHERE workspace_id = ? AND is_archived = ?
		 ORDER BY is_pinned DESC, updated_at DESC`, workspaceID, archived)
	return notes, err
}

func (r *NoteRepository) Update(ctx context.Context, n *model.Note) error {
	query := `UPDATE notes SET title = ?, slug = ?, content_md = ?, content_html = ?,
	          content_json = ?, folder_id = ?, is_pinned = ?, is_archived = ?, word_count = ?
	          WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query,
		n.Title, n.Slug, n.ContentMD, n.ContentHTML, n.ContentJSON,
		n.FolderID, n.IsPinned, n.IsArchived, n.WordCount, n.ID)
	return err
}

func (r *NoteRepository) Delete(ctx context.Context, id uint64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM notes WHERE id = ?", id)
	return err
}

func (r *NoteRepository) Search(ctx context.Context, workspaceID uint64, query string) ([]model.Note, error) {
	var notes []model.Note
	err := r.db.SelectContext(ctx, &notes,
		`SELECT * FROM notes WHERE workspace_id = ? AND MATCH(title, content_md) AGAINST(? IN BOOLEAN MODE)
		 ORDER BY updated_at DESC LIMIT 50`, workspaceID, query)
	return notes, err
}

func (r *NoteRepository) ListAllSlugs(ctx context.Context, workspaceID uint64) (map[string]uint64, error) {
	type slugRow struct {
		ID   uint64 `db:"id"`
		Slug string `db:"slug"`
	}
	var rows []slugRow
	err := r.db.SelectContext(ctx, &rows,
		"SELECT id, slug FROM notes WHERE workspace_id = ?", workspaceID)
	if err != nil {
		return nil, err
	}
	result := make(map[string]uint64, len(rows))
	for _, row := range rows {
		result[row.Slug] = row.ID
	}
	return result, nil
}
