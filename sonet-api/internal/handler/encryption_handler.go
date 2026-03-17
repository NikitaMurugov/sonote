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

type EncryptionHandler struct {
	userRepo *repository.UserRepository
}

func NewEncryptionHandler(userRepo *repository.UserRepository) *EncryptionHandler {
	return &EncryptionHandler{userRepo: userRepo}
}

type setupKeysRequest struct {
	UserSalt            string `json:"user_salt"`
	EncryptedPrivateKey string `json:"encrypted_private_key"`
	PublicKey           string `json:"public_key"`
	RecoveryDEK         string `json:"recovery_dek"`
}

func (h *EncryptionHandler) SetupKeys(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req setupKeysRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.UserSalt == "" || req.EncryptedPrivateKey == "" || req.PublicKey == "" || req.RecoveryDEK == "" {
		response.Error(w, http.StatusBadRequest, "all fields are required")
		return
	}

	if err := h.userRepo.UpdateEncryption(r.Context(), userID, req.UserSalt, req.EncryptedPrivateKey, req.PublicKey, req.RecoveryDEK); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to setup encryption keys")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "encryption keys saved"})
}

func (h *EncryptionHandler) GetKeys(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	user, err := h.userRepo.GetByID(r.Context(), userID)
	if err != nil || user == nil {
		response.Error(w, http.StatusNotFound, "user not found")
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"user_salt":             user.UserSalt,
		"encrypted_private_key": user.EncryptedPrivateKey,
		"public_key":            user.PublicKey,
		"has_encryption":        user.PublicKey != nil && *user.PublicKey != "",
	})
}

func (h *EncryptionHandler) GetPublicKey(w http.ResponseWriter, r *http.Request) {
	targetUserID, err := strconv.ParseUint(chi.URLParam(r, "userId"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid user id")
		return
	}

	key, err := h.userRepo.GetPublicKey(r.Context(), targetUserID)
	if err != nil || key == "" {
		response.Error(w, http.StatusNotFound, "public key not found")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"public_key": key})
}

func (h *EncryptionHandler) GetPublicKeyByEmail(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		response.Error(w, http.StatusBadRequest, "email is required")
		return
	}

	user, err := h.userRepo.GetByEmail(r.Context(), email)
	if err != nil || user == nil || user.PublicKey == nil || *user.PublicKey == "" {
		response.Error(w, http.StatusNotFound, "public key not found")
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"user_id":    user.ID,
		"public_key": *user.PublicKey,
	})
}

type recoverRequest struct {
	RecoveryKey         string `json:"recovery_key"`
	NewEncryptedPrivateKey string `json:"new_encrypted_private_key"`
	NewUserSalt         string `json:"new_user_salt"`
	NewRecoveryDEK      string `json:"new_recovery_dek"`
}

func (h *EncryptionHandler) RecoverKeys(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req recoverRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.userRepo.GetByID(r.Context(), userID)
	if err != nil || user == nil {
		response.Error(w, http.StatusNotFound, "user not found")
		return
	}

	if user.RecoveryDEK == nil || *user.RecoveryDEK == "" {
		response.Error(w, http.StatusBadRequest, "no recovery key configured")
		return
	}

	// Server stores the recovery-wrapped private key; the actual recovery verification
	// happens on the client side. Here we just update the keys with the new wrapping.
	publicKey := ""
	if user.PublicKey != nil {
		publicKey = *user.PublicKey
	}

	if err := h.userRepo.UpdateEncryption(r.Context(), userID,
		req.NewUserSalt, req.NewEncryptedPrivateKey, publicKey, req.NewRecoveryDEK); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to recover keys")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "keys recovered"})
}
