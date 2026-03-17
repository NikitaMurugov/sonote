package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/NikitaMurugov/sonote-api/internal/model"
)

type FolderRepository struct {
	db *sqlx.DB
}

func NewFolderRepository(db *sqlx.DB) *FolderRepository {
	return &FolderRepository{db: db}
}

func (r *FolderRepository) Create(ctx context.Context, f *model.Folder) error {
	query := `INSERT INTO folders (workspace_id, parent_id, name, sort_order) VALUES (?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query, f.WorkspaceID, f.ParentID, f.Name, f.SortOrder)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	f.ID = uint64(id)
	return nil
}

func (r *FolderRepository) GetByID(ctx context.Context, id uint64) (*model.Folder, error) {
	var f model.Folder
	err := r.db.GetContext(ctx, &f, "SELECT * FROM folders WHERE id = ?", id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &f, err
}

func (r *FolderRepository) ListByWorkspace(ctx context.Context, workspaceID uint64) ([]model.Folder, error) {
	var folders []model.Folder
	err := r.db.SelectContext(ctx, &folders,
		"SELECT * FROM folders WHERE workspace_id = ? ORDER BY sort_order, name", workspaceID)
	return folders, err
}

func (r *FolderRepository) Update(ctx context.Context, f *model.Folder) error {
	query := `UPDATE folders SET name = ?, parent_id = ?, sort_order = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, f.Name, f.ParentID, f.SortOrder, f.ID)
	return err
}

func (r *FolderRepository) Delete(ctx context.Context, id uint64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM folders WHERE id = ?", id)
	return err
}
