package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/NikitaMurugov/sonote-api/internal/model"
	"github.com/NikitaMurugov/sonote-api/internal/repository"
	"github.com/NikitaMurugov/sonote-api/pkg/hash"
	"github.com/NikitaMurugov/sonote-api/pkg/jwt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailTaken         = errors.New("email already registered")
	ErrUsernameTaken      = errors.New("username already taken")
	ErrInvalidToken       = errors.New("invalid or expired token")
)

type AuthService struct {
	userRepo  *repository.UserRepository
	tokenRepo *repository.TokenRepository
	wsRepo    *repository.WorkspaceRepository
	jwtSecret string
	accessExp time.Duration
	refreshExp time.Duration
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthResponse struct {
	User   *model.User `json:"user"`
	Tokens TokenPair   `json:"tokens"`
}

func NewAuthService(
	userRepo *repository.UserRepository,
	tokenRepo *repository.TokenRepository,
	wsRepo *repository.WorkspaceRepository,
	jwtSecret string,
	accessExp, refreshExp time.Duration,
) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		wsRepo:     wsRepo,
		jwtSecret:  jwtSecret,
		accessExp:  accessExp,
		refreshExp: refreshExp,
	}
}

func (s *AuthService) Register(ctx context.Context, email, username, password, displayName string) (*AuthResponse, error) {
	existing, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrEmailTaken
	}

	existing, err = s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrUsernameTaken
	}

	passwordHash, err := hash.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &model.User{
		Email:        email,
		Username:     username,
		DisplayName:  displayName,
		PasswordHash: passwordHash,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Create personal workspace
	slug := strings.ToLower(username) + "-personal"
	ws := &model.Workspace{
		Name:       "Personal",
		Slug:       slug,
		OwnerID:    user.ID,
		IsPersonal: true,
	}
	if err := s.wsRepo.Create(ctx, ws); err != nil {
		return nil, fmt.Errorf("failed to create personal workspace: %w", err)
	}

	// Add owner as admin member (needed for encrypted_dek storage)
	member := &model.WorkspaceMember{
		WorkspaceID: ws.ID,
		UserID:      user.ID,
		Role:        "admin",
	}
	if err := s.wsRepo.AddMember(ctx, member); err != nil {
		return nil, fmt.Errorf("failed to add owner as member: %w", err)
	}

	tokens, err := s.generateTokens(ctx, user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{User: user, Tokens: *tokens}, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if !hash.CheckPassword(password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	tokens, err := s.generateTokens(ctx, user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{User: user, Tokens: *tokens}, nil
}

func (s *AuthService) Refresh(ctx context.Context, rawToken string) (*TokenPair, error) {
	tokenHash := hashToken(rawToken)

	stored, err := s.tokenRepo.GetByHash(ctx, tokenHash)
	if err != nil {
		return nil, err
	}
	if stored == nil {
		return nil, ErrInvalidToken
	}
	if stored.ExpiresAt.Before(time.Now()) {
		return nil, ErrInvalidToken
	}

	// Revoke old token (rotation)
	if err := s.tokenRepo.Revoke(ctx, stored.ID); err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(ctx, stored.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidToken
	}

	return s.generateTokens(ctx, user)
}

func (s *AuthService) Logout(ctx context.Context, rawToken string) error {
	tokenHash := hashToken(rawToken)
	stored, err := s.tokenRepo.GetByHash(ctx, tokenHash)
	if err != nil {
		return err
	}
	if stored == nil {
		return nil
	}
	return s.tokenRepo.Revoke(ctx, stored.ID)
}

func (s *AuthService) generateTokens(ctx context.Context, user *model.User) (*TokenPair, error) {
	accessToken, err := jwt.GenerateAccessToken(user.ID, user.Email, s.jwtSecret, s.accessExp)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	rawRefresh, err := generateRandomToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	refreshHash := hashToken(rawRefresh)
	rt := &model.RefreshToken{
		UserID:    user.ID,
		TokenHash: refreshHash,
		ExpiresAt: time.Now().Add(s.refreshExp),
	}

	if err := s.tokenRepo.Create(ctx, rt); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: rawRefresh,
	}, nil
}

func generateRandomToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
