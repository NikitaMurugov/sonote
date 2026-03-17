package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/NikitaMurugov/sonote-api/internal/repository"
	"github.com/NikitaMurugov/sonote-api/pkg/response"
)

type GraphHandler struct {
	noteRepo *repository.NoteRepository
	linkRepo *repository.NoteLinkRepository
}

func NewGraphHandler(noteRepo *repository.NoteRepository, linkRepo *repository.NoteLinkRepository) *GraphHandler {
	return &GraphHandler{noteRepo: noteRepo, linkRepo: linkRepo}
}

type graphNode struct {
	ID    uint64 `json:"id"`
	Title string `json:"title"`
	Slug  string `json:"slug"`
}

type graphEdge struct {
	Source uint64 `json:"source"`
	Target uint64 `json:"target"`
}

type graphResponse struct {
	Nodes []graphNode `json:"nodes"`
	Edges []graphEdge `json:"edges"`
}

func (h *GraphHandler) GetGraph(w http.ResponseWriter, r *http.Request) {
	wsID, _ := strconv.ParseUint(chi.URLParam(r, "wsId"), 10, 64)

	notes, err := h.noteRepo.ListByWorkspace(r.Context(), wsID, nil, false)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to get notes")
		return
	}

	links, err := h.linkRepo.GetAllForWorkspace(r.Context(), wsID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to get links")
		return
	}

	nodes := make([]graphNode, len(notes))
	for i, n := range notes {
		nodes[i] = graphNode{ID: n.ID, Title: n.Title, Slug: n.Slug}
	}

	edges := make([]graphEdge, len(links))
	for i, l := range links {
		edges[i] = graphEdge{Source: l.SourceNoteID, Target: l.TargetNoteID}
	}

	response.JSON(w, http.StatusOK, graphResponse{Nodes: nodes, Edges: edges})
}
