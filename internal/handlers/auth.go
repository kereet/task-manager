package handlers

import (
	"encoding/json"
	"net/http"
	"task-manager/internal/auth"
	"task-manager/internal/database"
	"task-manager/internal/models"
)

type AuthHandlers struct {
	userStore *database.UserStore
}

func NewAuthHandlers(userStore *database.UserStore) *AuthHandlers {
	return &AuthHandlers{userStore: userStore}
}

func (h *AuthHandlers) Register(w http.ResponseWriter, r *http.Request) {
	var input models.RegisterInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if input.Name == "" {
		respondWithError(w, http.StatusBadRequest, "Name is required")
		return
	}

	if input.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Email is required")
		return
	}

	if len(input.Password) < 6 {
		respondWithError(w, http.StatusBadRequest, "Password must be at least 6 characters")
		return
	}

	user, err := h.userStore.Create(input)
	if err != nil {
		respondWithError(w, http.StatusConflict, "User with this email already exists")
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	var response models.AuthResponse
	response.Token = token
	response.User.ID = user.ID
	response.User.Name = user.Name
	response.User.Email = user.Email

	respondWithJSON(w, http.StatusCreated, response)
}

func (h *AuthHandlers) Login(w http.ResponseWriter, r *http.Request) {
	var input models.LoginInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.userStore.GetByEmail(input.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	if !auth.CheckPassword(input.Password, user.PasswordHash) {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	var response models.AuthResponse
	response.Token = token
	response.User.ID = user.ID
	response.User.Name = user.Name
	response.User.Email = user.Email

	respondWithJSON(w, http.StatusOK, response)
}
