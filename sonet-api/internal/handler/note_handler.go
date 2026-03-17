package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/NikitaMurugov/sonote-api/internal/middleware"
	"github.com/NikitaMurugov/sonote-api/internal/model"
	"github.com/NikitaMurugov/sonote-api/internal/repository"
	"github.com/NikitaMurugov/sonote-api/pkg/response"
)

type NoteHandler struct {
	noteRepo *repository.NoteRepository
	linkRepo *repository.NoteLinkRepository
}

func NewNoteHandler(noteRepo *repository.NoteRepository, linkRepo *repository.NoteLinkRepository) *NoteHandler {
	return &NoteHandler{noteRepo: noteRepo, linkRepo: linkRepo}
}

func (h *NoteHandler) List(w http.ResponseWriter, r *http.Request) {
	wsID, _ := strconv.ParseUint(chi.URLParam(r, "wsId"), 10, 64)

	var folderID *uint64
	if fStr := r.URL.Query().Get("folder_id"); fStr != "" {
		fid, err := strconv.ParseUint(fStr, 10, 64)
		if err == nil {
			folderID = &fid
		}
	}

	archived := r.URL.Query().Get("archived") == "true"

	notes, err := h.noteRepo.ListByWorkspace(r.Context(), wsID, folderID, archived)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to list notes")
		return
	}

	response.JSON(w, http.StatusOK, notes)
}

type createNoteRequest struct {
	Title       string              `json:"title"`
	FolderID    *uint64             `json:"folder_id"`
	ContentMD   string              `json:"content_md"`
	ContentHTML string              `json:"content_html"`
	ContentJSON model.NullRawMessage `json:"content_json"`
}

func (h *NoteHandler) Create(w http.ResponseWriter, r *http.Request) {
	wsID, _ := strconv.ParseUint(chi.URLParam(r, "wsId"), 10, 64)
	userID := middleware.GetUserID(r.Context())

	var req createNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if strings.TrimSpace(req.Title) == "" {
		response.Error(w, http.StatusBadRequest, "title is required")
		return
	}

	note := &model.Note{
		WorkspaceID: wsID,
		FolderID:    req.FolderID,
		Title:       req.Title,
		Slug:        generateSlug(req.Title),
		ContentMD:   req.ContentMD,
		ContentHTML: req.ContentHTML,
		ContentJSON: req.ContentJSON,
		AuthorID:    userID,
		WordCount:   uint32(countWords(req.ContentMD)),
	}

	if err := h.noteRepo.Create(r.Context(), note); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to create note")
		return
	}

	response.JSON(w, http.StatusCreated, note)
}

func (h *NoteHandler) Get(w http.ResponseWriter, r *http.Request) {
	nID, _ := strconv.ParseUint(chi.URLParam(r, "nId"), 10, 64)
	note, err := h.noteRepo.GetByID(r.Context(), nID)
	if err != nil || note == nil {
		response.Error(w, http.StatusNotFound, "note not found")
		return
	}
	response.JSON(w, http.StatusOK, note)
}

type updateNoteRequest struct {
	Title       *string               `json:"title"`
	FolderID    *uint64               `json:"folder_id"`
	ContentMD   *string               `json:"content_md"`
	ContentHTML *string               `json:"content_html"`
	ContentJSON *model.NullRawMessage `json:"content_json"`
	IsPinned    *bool                 `json:"is_pinned"`
	IsArchived  *bool                 `json:"is_archived"`
}

func (h *NoteHandler) Update(w http.ResponseWriter, r *http.Request) {
	nID, _ := strconv.ParseUint(chi.URLParam(r, "nId"), 10, 64)

	var req updateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	note, err := h.noteRepo.GetByID(r.Context(), nID)
	if err != nil || note == nil {
		response.Error(w, http.StatusNotFound, "note not found")
		return
	}

	if req.Title != nil {
		note.Title = *req.Title
		note.Slug = generateSlug(*req.Title)
	}
	if req.FolderID != nil {
		note.FolderID = req.FolderID
	}
	if req.ContentMD != nil {
		note.ContentMD = *req.ContentMD
		note.WordCount = uint32(countWords(*req.ContentMD))
	}
	if req.ContentHTML != nil {
		note.ContentHTML = *req.ContentHTML
	}
	if req.ContentJSON != nil {
		note.ContentJSON = *req.ContentJSON
	}
	if req.IsPinned != nil {
		note.IsPinned = *req.IsPinned
	}
	if req.IsArchived != nil {
		note.IsArchived = *req.IsArchived
	}

	if err := h.noteRepo.Update(r.Context(), note); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to update note")
		return
	}

	response.JSON(w, http.StatusOK, note)
}

func (h *NoteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	nID, _ := strconv.ParseUint(chi.URLParam(r, "nId"), 10, 64)
	if err := h.noteRepo.Delete(r.Context(), nID); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to delete note")
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "note deleted"})
}

func (h *NoteHandler) GetLinks(w http.ResponseWriter, r *http.Request) {
	nID, _ := strconv.ParseUint(chi.URLParam(r, "nId"), 10, 64)
	links, err := h.linkRepo.GetOutgoing(r.Context(), nID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to get links")
		return
	}
	response.JSON(w, http.StatusOK, links)
}

func (h *NoteHandler) GetBacklinks(w http.ResponseWriter, r *http.Request) {
	nID, _ := strconv.ParseUint(chi.URLParam(r, "nId"), 10, 64)
	links, err := h.linkRepo.GetBacklinks(r.Context(), nID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to get backlinks")
		return
	}
	response.JSON(w, http.StatusOK, links)
}

func countWords(s string) int {
	return len(strings.Fields(s))
}
