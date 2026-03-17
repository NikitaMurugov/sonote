package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/NikitaMurugov/sonote-api/internal/model"
)

type SessionRepository struct {
	db *sqlx.DB
}

func NewSessionRepository(db *sqlx.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(ctx context.Context, s *model.Session) error {
	query := `INSERT INTO user_sessions (user_id, refresh_token_id, device_name, device_type, os, browser, ip_address, location, is_current)
	           VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query,
		s.UserID, s.RefreshTokenID, s.DeviceName, s.DeviceType,
		s.OS, s.Browser, s.IPAddress, s.Location, s.IsCurrent)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	s.ID = uint64(id)
	return nil
}

func (r *SessionRepository) ListByUser(ctx context.Context, userID uint64) ([]model.Session, error) {
	var sessions []model.Session
	err := r.db.SelectContext(ctx, &sessions,
		"SELECT * FROM user_sessions WHERE user_id = ? ORDER BY last_active_at DESC", userID)
	return sessions, err
}

func (r *SessionRepository) Delete(ctx context.Context, id, userID uint64) error {
	_, err := r.db.ExecContext(ctx,
		"DELETE FROM user_sessions WHERE id = ? AND user_id = ?", id, userID)
	return err
}

func (r *SessionRepository) DeleteAllExcept(ctx context.Context, userID, exceptID uint64) error {
	_, err := r.db.ExecContext(ctx,
		"DELETE FROM user_sessions WHERE user_id = ? AND id != ?", userID, exceptID)
	return err
}

func (r *SessionRepository) UpdateLastActive(ctx context.Context, id uint64) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE user_sessions SET last_active_at = NOW() WHERE id = ?", id)
	return err
}
