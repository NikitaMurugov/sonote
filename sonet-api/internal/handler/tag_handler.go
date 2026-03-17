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

type TagHandler struct {
	tagRepo *repository.TagRepository
}

func NewTagHandler(tagRepo *repository.TagRepository) *TagHandler {
	return &TagHandler{tagRepo: tagRepo}
}

func (h *TagHandler) List(w http.ResponseWriter, r *http.Request) {
	wsID, _ := strconv.ParseUint(chi.URLParam(r, "wsId"), 10, 64)
	tags, err := h.tagRepo.ListByWorkspace(r.Context(), wsID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to list tags")
		return
	}
	response.JSON(w, http.StatusOK, tags)
}

type createTagRequest struct {
	Name  string  `json:"name"`
	Color *string `json:"color"`
}

func (h *TagHandler) Create(w http.ResponseWriter, r *http.Request) {
	wsID, _ := strconv.ParseUint(chi.URLParam(r, "wsId"), 10, 64)
	var req createTagRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		response.Error(w, http.StatusBadRequest, "name is required")
		return
	}

	tag := &model.Tag{
		WorkspaceID: wsID,
		Name:        req.Name,
		Color:       req.Color,
	}

	if err := h.tagRepo.Create(r.Context(), tag); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to create tag")
		return
	}

	response.JSON(w, http.StatusCreated, tag)
}

func (h *TagHandler) Update(w http.ResponseWriter, r *http.Request) {
	tID, _ := strconv.ParseUint(chi.URLParam(r, "tId"), 10, 64)
	var req createTagRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	tag, err := h.tagRepo.GetByID(r.Context(), tID)
	if err != nil || tag == nil {
		response.Error(w, http.StatusNotFound, "tag not found")
		return
	}

	if req.Name != "" {
		tag.Name = req.Name
	}
	if req.Color != nil {
		tag.Color = req.Color
	}

	if err := h.tagRepo.Update(r.Context(), tag); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to update tag")
		return
	}

	response.JSON(w, http.StatusOK, tag)
}

func (h *TagHandler) Delete(w http.ResponseWriter, r *http.Request) {
	tID, _ := strconv.ParseUint(chi.URLParam(r, "tId"), 10, 64)
	if err := h.tagRepo.Delete(r.Context(), tID); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to delete tag")
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "tag deleted"})
}

type attachTagsRequest struct {
	TagIDs []uint64 `json:"tag_ids"`
}

func (h *TagHandler) AttachToNote(w http.ResponseWriter, r *http.Request) {
	nID, _ := strconv.ParseUint(chi.URLParam(r, "nId"), 10, 64)
	var req attachTagsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	for _, tagID := range req.TagIDs {
		if err := h.tagRepo.AttachToNote(r.Context(), nID, tagID); err != nil {
			response.Error(w, http.StatusInternalServerError, "failed to attach tag")
			return
		}
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "tags attached"})
}

func (h *TagHandler) DetachFromNote(w http.ResponseWriter, r *http.Request) {
	nID, _ := strconv.ParseUint(chi.URLParam(r, "nId"), 10, 64)
	tID, _ := strconv.ParseUint(chi.URLParam(r, "tId"), 10, 64)

	if err := h.tagRepo.DetachFromNote(r.Context(), nID, tID); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to detach tag")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "tag detached"})
}

func (h *TagHandler) GetNoteTags(w http.ResponseWriter, r *http.Request) {
	nID, _ := strconv.ParseUint(chi.URLParam(r, "nId"), 10, 64)
	tags, err := h.tagRepo.GetTagsForNote(r.Context(), nID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to get tags")
		return
	}
	response.JSON(w, http.StatusOK, tags)
}
