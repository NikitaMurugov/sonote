package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/NikitaMurugov/sonote-api/internal/model"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	query := `INSERT INTO users (email, username, display_name, password_hash) VALUES (?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query, user.Email, user.Username, user.DisplayName, user.PasswordHash)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = uint64(id)
	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uint64) (*model.User, error) {
	var user model.User
	err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = ?", id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &user, err
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE email = ?", email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &user, err
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE username = ?", username)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &user, err
}

func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	query := `UPDATE users SET display_name = ?, avatar_url = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, user.DisplayName, user.AvatarURL, user.ID)
	return err
}

func (r *UserRepository) UpdatePassword(ctx context.Context, id uint64, hash string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE users SET password_hash = ? WHERE id = ?", hash, id)
	return err
}

func (r *UserRepository) UpdateEncryption(ctx context.Context, id uint64, userSalt, encryptedPrivateKey, publicKey, recoveryDEK string) error {
	query := `UPDATE users SET user_salt = ?, encrypted_private_key = ?, public_key = ?, recovery_dek = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, userSalt, encryptedPrivateKey, publicKey, recoveryDEK, id)
	return err
}

func (r *UserRepository) GetPublicKey(ctx context.Context, id uint64) (string, error) {
	var key string
	err := r.db.GetContext(ctx, &key, "SELECT COALESCE(public_key, '') FROM users WHERE id = ?", id)
	return key, err
}
