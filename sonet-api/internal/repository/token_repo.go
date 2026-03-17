package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/NikitaMurugov/sonote-api/internal/model"
)

type TokenRepository struct {
	db *sqlx.DB
}

func NewTokenRepository(db *sqlx.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

func (r *TokenRepository) Create(ctx context.Context, token *model.RefreshToken) error {
	query := `INSERT INTO refresh_tokens (user_id, token_hash, device_info, ip_address, expires_at)
	           VALUES (?, ?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query,
		token.UserID, token.TokenHash, token.DeviceInfo, token.IPAddress, token.ExpiresAt)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	token.ID = uint64(id)
	return nil
}

func (r *TokenRepository) GetByHash(ctx context.Context, hash string) (*model.RefreshToken, error) {
	var token model.RefreshToken
	err := r.db.GetContext(ctx, &token,
		"SELECT * FROM refresh_tokens WHERE token_hash = ? AND revoked = FALSE", hash)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &token, err
}

func (r *TokenRepository) Revoke(ctx context.Context, id uint64) error {
	_, err := r.db.ExecContext(ctx, "UPDATE refresh_tokens SET revoked = TRUE WHERE id = ?", id)
	return err
}

func (r *TokenRepository) RevokeAllForUser(ctx context.Context, userID uint64) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE refresh_tokens SET revoked = TRUE WHERE user_id = ?", userID)
	return err
}
