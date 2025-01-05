package auth

import (
	"backend-bootcamp-assignment-2024/dto"
	"encoding/json"
	"net/http"
)

type AuthHandler struct {
	service AuthService
}

func NewAuthHandler(service AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) DummyLoginHandler(w http.ResponseWriter, r *http.Request) {

	userTypeStr := r.URL.Query().Get("user_type")

	userType := dto.UserType(userTypeStr)

	if userType != dto.Client && userType != dto.Moderator {
		http.Error(w, "Invalid user_type. Allowed values: client, moderator", http.StatusBadRequest)
		return
	}

	token, err := h.service.generateJWToken(userTypeStr)

	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.PostRegisterJSONRequestBody

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	ctx := r.Context()

	userUUID, err := h.service.registerUser(ctx, req)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"uuid": userUUID})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.PostLoginJSONRequestBody

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	token, err := h.service.loginUser(ctx, req)

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
	w.WriteHeader(http.StatusOK)
}
