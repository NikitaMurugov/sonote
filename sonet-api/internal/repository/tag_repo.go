package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/NikitaMurugov/sonote-api/internal/model"
)

type TagRepository struct {
	db *sqlx.DB
}

func NewTagRepository(db *sqlx.DB) *TagRepository {
	return &TagRepository{db: db}
}

func (r *TagRepository) Create(ctx context.Context, tag *model.Tag) error {
	query := `INSERT INTO tags (workspace_id, name, color) VALUES (?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query, tag.WorkspaceID, tag.Name, tag.Color)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	tag.ID = uint64(id)
	return nil
}

func (r *TagRepository) GetByID(ctx context.Context, id uint64) (*model.Tag, error) {
	var tag model.Tag
	err := r.db.GetContext(ctx, &tag, "SELECT * FROM tags WHERE id = ?", id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &tag, err
}

func (r *TagRepository) ListByWorkspace(ctx context.Context, workspaceID uint64) ([]model.Tag, error) {
	var tags []model.Tag
	err := r.db.SelectContext(ctx, &tags,
		"SELECT * FROM tags WHERE workspace_id = ? ORDER BY name", workspaceID)
	return tags, err
}

func (r *TagRepository) Update(ctx context.Context, tag *model.Tag) error {
	_, err := r.db.ExecContext(ctx, "UPDATE tags SET name = ?, color = ? WHERE id = ?",
		tag.Name, tag.Color, tag.ID)
	return err
}

func (r *TagRepository) Delete(ctx context.Context, id uint64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM tags WHERE id = ?", id)
	return err
}

func (r *TagRepository) AttachToNote(ctx context.Context, noteID, tagID uint64) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT IGNORE INTO note_tags (note_id, tag_id) VALUES (?, ?)", noteID, tagID)
	return err
}

func (r *TagRepository) DetachFromNote(ctx context.Context, noteID, tagID uint64) error {
	_, err := r.db.ExecContext(ctx,
		"DELETE FROM note_tags WHERE note_id = ? AND tag_id = ?", noteID, tagID)
	return err
}

func (r *TagRepository) GetTagsForNote(ctx context.Context, noteID uint64) ([]model.Tag, error) {
	var tags []model.Tag
	err := r.db.SelectContext(ctx, &tags,
		`SELECT t.* FROM tags t
		 JOIN note_tags nt ON t.id = nt.tag_id
		 WHERE nt.note_id = ? ORDER BY t.name`, noteID)
	return tags, err
}
