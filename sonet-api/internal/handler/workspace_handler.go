package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/NikitaMurugov/sonote-api/internal/middleware"
	"github.com/NikitaMurugov/sonote-api/internal/model"
	"github.com/NikitaMurugov/sonote-api/internal/repository"
	"github.com/NikitaMurugov/sonote-api/pkg/response"
)

type WorkspaceHandler struct {
	wsRepo   *repository.WorkspaceRepository
	userRepo *repository.UserRepository
}

func NewWorkspaceHandler(wsRepo *repository.WorkspaceRepository, userRepo *repository.UserRepository) *WorkspaceHandler {
	return &WorkspaceHandler{wsRepo: wsRepo, userRepo: userRepo}
}

func (h *WorkspaceHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	workspaces, err := h.wsRepo.ListByUserID(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to list workspaces")
		return
	}
	response.JSON(w, http.StatusOK, workspaces)
}

type createWorkspaceRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Icon        *string `json:"icon"`
}

func (h *WorkspaceHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	var req createWorkspaceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		response.Error(w, http.StatusBadRequest, "name is required")
		return
	}

	slug := generateSlug(req.Name)
	ws := &model.Workspace{
		Name:        req.Name,
		Slug:        slug,
		Description: req.Description,
		OwnerID:     userID,
		Icon:        req.Icon,
	}

	if err := h.wsRepo.Create(r.Context(), ws); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to create workspace")
		return
	}

	// Add owner as admin member
	member := &model.WorkspaceMember{
		WorkspaceID: ws.ID,
		UserID:      userID,
		Role:        "admin",
	}
	_ = h.wsRepo.AddMember(r.Context(), member)

	response.JSON(w, http.StatusCreated, ws)
}

func (h *WorkspaceHandler) Get(w http.ResponseWriter, r *http.Request) {
	wsID, _ := strconv.ParseUint(chi.URLParam(r, "wsId"), 10, 64)
	ws, err := h.wsRepo.GetByID(r.Context(), wsID)
	if err != nil || ws == nil {
		response.Error(w, http.StatusNotFound, "workspace not found")
		return
	}
	response.JSON(w, http.StatusOK, ws)
}

type updateWorkspaceRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Icon        *string `json:"icon"`
}

func (h *WorkspaceHandler) Update(w http.ResponseWriter, r *http.Request) {
	wsID, _ := strconv.ParseUint(chi.URLParam(r, "wsId"), 10, 64)
	var req updateWorkspaceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	ws, err := h.wsRepo.GetByID(r.Context(), wsID)
	if err != nil || ws == nil {
		response.Error(w, http.StatusNotFound, "workspace not found")
		return
	}

	if req.Name != nil {
		ws.Name = *req.Name
	}
	if req.Description != nil {
		ws.Description = req.Description
	}
	if req.Icon != nil {
		ws.Icon = req.Icon
	}

	if err := h.wsRepo.Update(r.Context(), ws); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to update workspace")
		return
	}

	response.JSON(w, http.StatusOK, ws)
}

func (h *WorkspaceHandler) Delete(w http.ResponseWriter, r *http.Request) {
	wsID, _ := strconv.ParseUint(chi.URLParam(r, "wsId"), 10, 64)
	userID := middleware.GetUserID(r.Context())

	ws, err := h.wsRepo.GetByID(r.Context(), wsID)
	if err != nil || ws == nil {
		response.Error(w, http.StatusNotFound, "workspace not found")
		return
	}

	if ws.OwnerID != userID {
		response.Error(w, http.StatusForbidden, "only the owner can delete a workspace")
		return
	}

	if ws.IsPersonal {
		response.Error(w, http.StatusBadRequest, "cannot delete personal workspace")
		return
	}

	if err := h.wsRepo.Delete(r.Context(), wsID); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to delete workspace")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "workspace deleted"})
}

func (h *WorkspaceHandler) ListMembers(w http.ResponseWriter, r *http.Request) {
	wsID, _ := strconv.ParseUint(chi.URLParam(r, "wsId"), 10, 64)
	members, err := h.wsRepo.ListMembers(r.Context(), wsID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to list members")
		return
	}
	response.JSON(w, http.StatusOK, members)
}

type inviteMemberRequest struct {
	Email        string  `json:"email"`
	Role         string  `json:"role"`
	EncryptedDEK *string `json:"encrypted_dek"`
}

func (h *WorkspaceHandler) InviteMember(w http.ResponseWriter, r *http.Request) {
	wsID, _ := strconv.ParseUint(chi.URLParam(r, "wsId"), 10, 64)
	userID := middleware.GetUserID(r.Context())

	var req inviteMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Role != "viewer" && req.Role != "editor" && req.Role != "admin" {
		response.Error(w, http.StatusBadRequest, "role must be viewer, editor, or admin")
		return
	}

	invitee, err := h.userRepo.GetByEmail(r.Context(), req.Email)
	if err != nil || invitee == nil {
		response.Error(w, http.StatusNotFound, "user not found")
		return
	}

	member := &model.WorkspaceMember{
		WorkspaceID:  wsID,
		UserID:       invitee.ID,
		Role:         req.Role,
		InvitedBy:    &userID,
		EncryptedDEK: req.EncryptedDEK,
	}

	if err := h.wsRepo.AddMember(r.Context(), member); err != nil {
		response.Error(w, http.StatusConflict, "user is already a member")
		return
	}

	response.JSON(w, http.StatusCreated, member)
}

func (h *WorkspaceHandler) UpdateMemberRole(w http.ResponseWriter, r *http.Request) {
	wsID, _ := strconv.ParseUint(chi.URLParam(r, "wsId"), 10, 64)
	memberUserID, _ := strconv.ParseUint(chi.URLParam(r, "userId"), 10, 64)

	var req struct {
		Role string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Role != "viewer" && req.Role != "editor" && req.Role != "admin" {
		response.Error(w, http.StatusBadRequest, "role must be viewer, editor, or admin")
		return
	}

	if err := h.wsRepo.UpdateMemberRole(r.Context(), wsID, memberUserID, req.Role); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to update role")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "role updated"})
}

func (h *WorkspaceHandler) UpdateMemberDEK(w http.ResponseWriter, r *http.Request) {
	wsID, _ := strconv.ParseUint(chi.URLParam(r, "wsId"), 10, 64)
	userID := middleware.GetUserID(r.Context())

	var req struct {
		EncryptedDEK string `json:"encrypted_dek"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.wsRepo.UpdateMemberDEK(r.Context(), wsID, userID, req.EncryptedDEK); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to update dek")
		return
	}

	// Also mark workspace as encrypted
	_ = h.wsRepo.SetEncrypted(r.Context(), wsID, true)

	response.JSON(w, http.StatusOK, map[string]string{"message": "dek updated"})
}

func (h *WorkspaceHandler) GetMyMember(w http.ResponseWriter, r *http.Request) {
	wsID, _ := strconv.ParseUint(chi.URLParam(r, "wsId"), 10, 64)
	userID := middleware.GetUserID(r.Context())

	member, err := h.wsRepo.GetMember(r.Context(), wsID, userID)
	if err != nil || member == nil {
		response.Error(w, http.StatusNotFound, "not a member")
		return
	}

	response.JSON(w, http.StatusOK, member)
}

func (h *WorkspaceHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	wsID, _ := strconv.ParseUint(chi.URLParam(r, "wsId"), 10, 64)
	memberUserID, _ := strconv.ParseUint(chi.URLParam(r, "userId"), 10, 64)

	if err := h.wsRepo.RemoveMember(r.Context(), wsID, memberUserID); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to remove member")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "member removed"})
}

var translitMap = map[rune]string{
	'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e", 'ё': "yo",
	'ж': "zh", 'з': "z", 'и': "i", 'й': "y", 'к': "k", 'л': "l", 'м': "m",
	'н': "n", 'о': "o", 'п': "p", 'р': "r", 'с': "s", 'т': "t", 'у': "u",
	'ф': "f", 'х': "kh", 'ц': "ts", 'ч': "ch", 'ш': "sh", 'щ': "sch",
	'ъ': "", 'ы': "y", 'ь': "", 'э': "e", 'ю': "yu", 'я': "ya",
}

func generateSlug(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	var b strings.Builder
	for _, r := range name {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			b.WriteRune(r)
		} else if r == ' ' || r == '-' || r == '_' {
			b.WriteRune('-')
		} else if t, ok := translitMap[r]; ok {
			b.WriteString(t)
		}
	}
	slug := b.String()
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}
	slug = strings.Trim(slug, "-")
	if slug == "" {
		slug = fmt.Sprintf("note-%d", time.Now().UnixMilli())
	}
	return slug
}
