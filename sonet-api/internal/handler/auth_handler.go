package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/NikitaMurugov/sonote-api/internal/service"
	"github.com/NikitaMurugov/sonote-api/pkg/response"
	"github.com/NikitaMurugov/sonote-api/pkg/validator"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type registerRequest struct {
	Email       string `json:"email"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type logoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if !validator.ValidateEmail(req.Email) {
		response.Error(w, http.StatusBadRequest, "invalid email format")
		return
	}
	if !validator.ValidateUsername(req.Username) {
		response.Error(w, http.StatusBadRequest, "username must be 3-64 characters, alphanumeric, hyphens or underscores")
		return
	}
	if !validator.ValidatePassword(req.Password) {
		response.Error(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}
	if !validator.ValidateRequired(req.DisplayName) {
		req.DisplayName = req.Username
	}

	result, err := h.authService.Register(r.Context(), req.Email, req.Username, req.Password, req.DisplayName)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmailTaken):
			response.Error(w, http.StatusConflict, err.Error())
		case errors.Is(err, service.ErrUsernameTaken):
			response.Error(w, http.StatusConflict, err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, "registration failed")
		}
		return
	}

	response.JSON(w, http.StatusCreated, result)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if !validator.ValidateEmail(req.Email) || !validator.ValidateRequired(req.Password) {
		response.Error(w, http.StatusBadRequest, "email and password are required")
		return
	}

	result, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			response.Error(w, http.StatusUnauthorized, err.Error())
		} else {
			response.Error(w, http.StatusInternalServerError, "login failed")
		}
		return
	}

	response.JSON(w, http.StatusOK, result)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.RefreshToken == "" {
		response.Error(w, http.StatusBadRequest, "refresh_token is required")
		return
	}

	tokens, err := h.authService.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, service.ErrInvalidToken) {
			response.Error(w, http.StatusUnauthorized, err.Error())
		} else {
			response.Error(w, http.StatusInternalServerError, "token refresh failed")
		}
		return
	}

	response.JSON(w, http.StatusOK, tokens)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req logoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.authService.Logout(r.Context(), req.RefreshToken); err != nil {
		response.Error(w, http.StatusInternalServerError, "logout failed")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "logged out"})
}
