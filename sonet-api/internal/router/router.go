package router

import (
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/NikitaMurugov/sonote-api/internal/handler"
	"github.com/NikitaMurugov/sonote-api/internal/middleware"
	"github.com/NikitaMurugov/sonote-api/internal/repository"
)

type Handlers struct {
	Auth       *handler.AuthHandler
	User       *handler.UserHandler
	Workspace  *handler.WorkspaceHandler
	Folder     *handler.FolderHandler
	Note       *handler.NoteHandler
	Tag        *handler.TagHandler
	Search     *handler.SearchHandler
	Graph      *handler.GraphHandler
	Encryption *handler.EncryptionHandler
	Session    *handler.SessionHandler
}

func New(h Handlers, wsRepo *repository.WorkspaceRepository, jwtSecret, corsOrigins string) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(middleware.CORS(corsOrigins))

	r.Route("/api/v1", func(r chi.Router) {
		// Public routes
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", h.Auth.Register)
			r.Post("/login", h.Auth.Login)
			r.Post("/refresh", h.Auth.Refresh)
			r.Post("/logout", h.Auth.Logout)
		})

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(jwtSecret))

			// Users
			r.Get("/users/me", h.User.GetMe)
			r.Patch("/users/me", h.User.UpdateMe)
			r.Put("/users/me/password", h.User.ChangePassword)

			// Encryption
			r.Post("/encryption/setup", h.Encryption.SetupKeys)
			r.Get("/encryption/keys", h.Encryption.GetKeys)
			r.Get("/encryption/public-key/{userId}", h.Encryption.GetPublicKey)
			r.Get("/encryption/public-key-by-email", h.Encryption.GetPublicKeyByEmail)
			r.Post("/encryption/recover", h.Encryption.RecoverKeys)

			// Sessions
			r.Get("/sessions", h.Session.List)
			r.Delete("/sessions/{sessionId}", h.Session.Revoke)
			r.Post("/sessions/revoke-all", h.Session.RevokeAll)

			// Workspaces
			r.Get("/workspaces", h.Workspace.List)
			r.Post("/workspaces", h.Workspace.Create)

			r.Route("/workspaces/{wsId}", func(r chi.Router) {
				r.Use(middleware.WorkspaceAccess(wsRepo, "viewer"))

				r.Get("/", h.Workspace.Get)
				r.Get("/members", h.Workspace.ListMembers)
				r.Get("/members/me", h.Workspace.GetMyMember)
				r.Put("/members/me/dek", h.Workspace.UpdateMemberDEK)
				r.Get("/graph", h.Graph.GetGraph)
				r.Get("/search", h.Search.Search)

				// Folders (viewer can read)
				r.Get("/folders", h.Folder.List)

				// Notes (viewer can read)
				r.Get("/notes", h.Note.List)
				r.Get("/notes/{nId}", h.Note.Get)
				r.Get("/notes/{nId}/links", h.Note.GetLinks)
				r.Get("/notes/{nId}/backlinks", h.Note.GetBacklinks)
				r.Get("/notes/{nId}/tags", h.Tag.GetNoteTags)

				// Tags (viewer can read)
				r.Get("/tags", h.Tag.List)

				// Editor+ routes
				r.Group(func(r chi.Router) {
					r.Use(middleware.WorkspaceAccess(wsRepo, "editor"))

					r.Post("/folders", h.Folder.Create)
					r.Patch("/folders/{fId}", h.Folder.Update)
					r.Delete("/folders/{fId}", h.Folder.Delete)

					r.Post("/notes", h.Note.Create)
					r.Patch("/notes/{nId}", h.Note.Update)
					r.Delete("/notes/{nId}", h.Note.Delete)

					r.Post("/tags", h.Tag.Create)
					r.Patch("/tags/{tId}", h.Tag.Update)
					r.Delete("/tags/{tId}", h.Tag.Delete)
					r.Post("/notes/{nId}/tags", h.Tag.AttachToNote)
					r.Delete("/notes/{nId}/tags/{tId}", h.Tag.DetachFromNote)
				})

				// Admin+ routes
				r.Group(func(r chi.Router) {
					r.Use(middleware.WorkspaceAccess(wsRepo, "admin"))

					r.Patch("/", h.Workspace.Update)
					r.Post("/members", h.Workspace.InviteMember)
					r.Patch("/members/{userId}", h.Workspace.UpdateMemberRole)
					r.Delete("/members/{userId}", h.Workspace.RemoveMember)
				})

				// Owner only
				r.Delete("/", h.Workspace.Delete)
			})
		})
	})

	return r
}
