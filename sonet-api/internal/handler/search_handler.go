package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/NikitaMurugov/sonote-api/internal/repository"
	"github.com/NikitaMurugov/sonote-api/pkg/response"
)

type SearchHandler struct {
	noteRepo *repository.NoteRepository
}

func NewSearchHandler(noteRepo *repository.NoteRepository) *SearchHandler {
	return &SearchHandler{noteRepo: noteRepo}
}

func (h *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	wsID, _ := strconv.ParseUint(chi.URLParam(r, "wsId"), 10, 64)
	query := r.URL.Query().Get("q")

	if query == "" {
		response.Error(w, http.StatusBadRequest, "query parameter 'q' is required")
		return
	}

	notes, err := h.noteRepo.Search(r.Context(), wsID, query)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "search failed")
		return
	}

	response.JSON(w, http.StatusOK, notes)
}
