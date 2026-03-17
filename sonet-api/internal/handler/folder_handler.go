package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/NikitaMurugov/sonote-api/internal/model"
	"github.com/NikitaMurugov/sonote-api/internal/repository"
	"github.com/NikitaMurugov/sonote-api/pkg/response"
)

type FolderHandler struct {
	folderRepo *repository.FolderRepository
}

func NewFolderHandler(folderRepo *repository.FolderRepository) *FolderHandler {
	return &FolderHandler{folderRepo: folderRepo}
}

func (h *FolderHandler) List(w http.ResponseWriter, r *http.Request) {
	wsID, _ := strconv.ParseUint(chi.URLParam(r, "wsId"), 10, 64)
	folders, err := h.folderRepo.ListByWorkspace(r.Context(), wsID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to list folders")
		return
	}
	response.JSON(w, http.StatusOK, folders)
}

type createFolderRequest struct {
	Name     string  `json:"name"`
	ParentID *uint64 `json:"parent_id"`
}

func (h *FolderHandler) Create(w http.ResponseWriter, r *http.Request) {
	wsID, _ := strconv.ParseUint(chi.URLParam(r, "wsId"), 10, 64)
	var req createFolderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		response.Error(w, http.StatusBadRequest, "name is required")
		return
	}

	folder := &model.Folder{
		WorkspaceID: wsID,
		ParentID:    req.ParentID,
		Name:        req.Name,
	}

	if err := h.folderRepo.Create(r.Context(), folder); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to create folder")
		return
	}

	response.JSON(w, http.StatusCreated, folder)
}

type updateFolderRequest struct {
	Name      *string `json:"name"`
	ParentID  *uint64 `json:"parent_id"`
	SortOrder *int    `json:"sort_order"`
}

func (h *FolderHandler) Update(w http.ResponseWriter, r *http.Request) {
	fID, _ := strconv.ParseUint(chi.URLParam(r, "fId"), 10, 64)
	var req updateFolderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	folder, err := h.folderRepo.GetByID(r.Context(), fID)
	if err != nil || folder == nil {
		response.Error(w, http.StatusNotFound, "folder not found")
		return
	}

	if req.Name != nil {
		folder.Name = *req.Name
	}
	if req.ParentID != nil {
		folder.ParentID = req.ParentID
	}
	if req.SortOrder != nil {
		folder.SortOrder = *req.SortOrder
	}

	if err := h.folderRepo.Update(r.Context(), folder); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to update folder")
		return
	}

	response.JSON(w, http.StatusOK, folder)
}

func (h *FolderHandler) Delete(w http.ResponseWriter, r *http.Request) {
	fID, _ := strconv.ParseUint(chi.URLParam(r, "fId"), 10, 64)
	if err := h.folderRepo.Delete(r.Context(), fID); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to delete folder")
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "folder deleted"})
}
