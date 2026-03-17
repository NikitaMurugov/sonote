package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/NikitaMurugov/sonote-api/internal/config"
	"github.com/NikitaMurugov/sonote-api/internal/database"
	"github.com/NikitaMurugov/sonote-api/internal/handler"
	"github.com/NikitaMurugov/sonote-api/internal/repository"
	"github.com/NikitaMurugov/sonote-api/internal/router"
	"github.com/NikitaMurugov/sonote-api/internal/service"
)

func main() {
	// Logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	// Config
	cfg := config.Load()

	// Database
	db, err := database.Connect(cfg.Database.DSN())
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Repositories
	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	wsRepo := repository.NewWorkspaceRepository(db)
	folderRepo := repository.NewFolderRepository(db)
	noteRepo := repository.NewNoteRepository(db)
	linkRepo := repository.NewNoteLinkRepository(db)
	tagRepo := repository.NewTagRepository(db)
	sessionRepo := repository.NewSessionRepository(db)

	// Services
	authService := service.NewAuthService(
		userRepo, tokenRepo, wsRepo,
		cfg.JWT.Secret, cfg.JWT.AccessTokenExpiry, cfg.JWT.RefreshTokenExpiry,
	)

	// Handlers
	handlers := router.Handlers{
		Auth:       handler.NewAuthHandler(authService),
		User:       handler.NewUserHandler(userRepo),
		Workspace:  handler.NewWorkspaceHandler(wsRepo, userRepo),
		Folder:     handler.NewFolderHandler(folderRepo),
		Note:       handler.NewNoteHandler(noteRepo, linkRepo),
		Tag:        handler.NewTagHandler(tagRepo),
		Search:     handler.NewSearchHandler(noteRepo),
		Graph:      handler.NewGraphHandler(noteRepo, linkRepo),
		Encryption: handler.NewEncryptionHandler(userRepo),
		Session:    handler.NewSessionHandler(sessionRepo),
	}

	// Router
	r := router.New(handlers, wsRepo, cfg.JWT.Secret, cfg.Server.AllowOrigins)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	slog.Info("starting Sonet API", "addr", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
