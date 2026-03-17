package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/NikitaMurugov/sonote-api/internal/model"
)

type WorkspaceRepository struct {
	db *sqlx.DB
}

func NewWorkspaceRepository(db *sqlx.DB) *WorkspaceRepository {
	return &WorkspaceRepository{db: db}
}

func (r *WorkspaceRepository) Create(ctx context.Context, ws *model.Workspace) error {
	query := `INSERT INTO workspaces (name, slug, description, owner_id, is_personal, icon) VALUES (?, ?, ?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query, ws.Name, ws.Slug, ws.Description, ws.OwnerID, ws.IsPersonal, ws.Icon)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	ws.ID = uint64(id)
	return nil
}

func (r *WorkspaceRepository) GetByID(ctx context.Context, id uint64) (*model.Workspace, error) {
	var ws model.Workspace
	err := r.db.GetContext(ctx, &ws, "SELECT * FROM workspaces WHERE id = ?", id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &ws, err
}

func (r *WorkspaceRepository) GetBySlug(ctx context.Context, slug string) (*model.Workspace, error) {
	var ws model.Workspace
	err := r.db.GetContext(ctx, &ws, "SELECT * FROM workspaces WHERE slug = ?", slug)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &ws, err
}

func (r *WorkspaceRepository) ListByUserID(ctx context.Context, userID uint64) ([]model.Workspace, error) {
	var workspaces []model.Workspace
	query := `SELECT w.* FROM workspaces w
	          LEFT JOIN workspace_members wm ON w.id = wm.workspace_id
	          WHERE w.owner_id = ? OR wm.user_id = ?
	          GROUP BY w.id ORDER BY w.name`
	err := r.db.SelectContext(ctx, &workspaces, query, userID, userID)
	return workspaces, err
}

func (r *WorkspaceRepository) Update(ctx context.Context, ws *model.Workspace) error {
	query := `UPDATE workspaces SET name = ?, description = ?, icon = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, ws.Name, ws.Description, ws.Icon, ws.ID)
	return err
}

func (r *WorkspaceRepository) Delete(ctx context.Context, id uint64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM workspaces WHERE id = ?", id)
	return err
}

func (r *WorkspaceRepository) AddMember(ctx context.Context, member *model.WorkspaceMember) error {
	query := `INSERT INTO workspace_members (workspace_id, user_id, role, invited_by) VALUES (?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query, member.WorkspaceID, member.UserID, member.Role, member.InvitedBy)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	member.ID = uint64(id)
	return nil
}

func (r *WorkspaceRepository) GetMember(ctx context.Context, workspaceID, userID uint64) (*model.WorkspaceMember, error) {
	var member model.WorkspaceMember
	err := r.db.GetContext(ctx, &member,
		"SELECT * FROM workspace_members WHERE workspace_id = ? AND user_id = ?",
		workspaceID, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &member, err
}

func (r *WorkspaceRepository) ListMembers(ctx context.Context, workspaceID uint64) ([]model.WorkspaceMember, error) {
	var members []model.WorkspaceMember
	err := r.db.SelectContext(ctx, &members,
		"SELECT * FROM workspace_members WHERE workspace_id = ?", workspaceID)
	return members, err
}

func (r *WorkspaceRepository) UpdateMemberRole(ctx context.Context, workspaceID, userID uint64, role string) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE workspace_members SET role = ? WHERE workspace_id = ? AND user_id = ?",
		role, workspaceID, userID)
	return err
}

func (r *WorkspaceRepository) RemoveMember(ctx context.Context, workspaceID, userID uint64) error {
	_, err := r.db.ExecContext(ctx,
		"DELETE FROM workspace_members WHERE workspace_id = ? AND user_id = ?",
		workspaceID, userID)
	return err
}
