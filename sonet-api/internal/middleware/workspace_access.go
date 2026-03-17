package middleware

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/NikitaMurugov/sonote-api/internal/repository"
	"github.com/NikitaMurugov/sonote-api/pkg/response"
)

func WorkspaceAccess(wsRepo *repository.WorkspaceRepository, minRole string) func(http.Handler) http.Handler {
	roleLevel := map[string]int{
		"viewer": 1,
		"editor": 2,
		"admin":  3,
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := GetUserID(r.Context())
			if userID == 0 {
				response.Error(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			wsIDStr := chi.URLParam(r, "wsId")
			wsID, err := strconv.ParseUint(wsIDStr, 10, 64)
			if err != nil {
				response.Error(w, http.StatusBadRequest, "invalid workspace id")
				return
			}

			// Check if owner
			ws, err := wsRepo.GetByID(r.Context(), wsID)
			if err != nil {
				response.Error(w, http.StatusInternalServerError, "failed to check workspace")
				return
			}
			if ws == nil {
				response.Error(w, http.StatusNotFound, "workspace not found")
				return
			}

			if ws.OwnerID == userID {
				next.ServeHTTP(w, r)
				return
			}

			// Check membership
			member, err := wsRepo.GetMember(r.Context(), wsID, userID)
			if err != nil {
				response.Error(w, http.StatusInternalServerError, "failed to check access")
				return
			}
			if member == nil {
				response.Error(w, http.StatusForbidden, "no access to this workspace")
				return
			}

			if roleLevel[member.Role] < roleLevel[minRole] {
				response.Error(w, http.StatusForbidden, "insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
