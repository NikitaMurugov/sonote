package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/NikitaMurugov/sonote-api/internal/middleware"
	"github.com/NikitaMurugov/sonote-api/internal/repository"
	"github.com/NikitaMurugov/sonote-api/pkg/response"
)

type SessionHandler struct {
	sessionRepo *repository.SessionRepository
}

func NewSessionHandler(sessionRepo *repository.SessionRepository) *SessionHandler {
	return &SessionHandler{sessionRepo: sessionRepo}
}

func (h *SessionHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	sessions, err := h.sessionRepo.ListByUser(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to list sessions")
		return
	}
	response.JSON(w, http.StatusOK, sessions)
}

func (h *SessionHandler) Revoke(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	sessionID, err := strconv.ParseUint(chi.URLParam(r, "sessionId"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid session id")
		return
	}

	if err := h.sessionRepo.Delete(r.Context(), sessionID, userID); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to revoke session")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "session revoked"})
}

type revokeAllRequest struct {
	CurrentSessionID uint64 `json:"current_session_id"`
}

func (h *SessionHandler) RevokeAll(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req revokeAllRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.sessionRepo.DeleteAllExcept(r.Context(), userID, req.CurrentSessionID); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to revoke sessions")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "all other sessions revoked"})
}
