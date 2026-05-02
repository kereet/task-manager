package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"task-manager/internal/auth"
	"task-manager/internal/database"
	"task-manager/internal/models"
)

type Handlers struct {
	store *database.TaskStore
}

func NewHandlers(store *database.TaskStore) *Handlers {
	return &Handlers{store}
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	respondWithJSON(w, statusCode, map[string]string{"error": message})
}

// GetAllTasks - проверяем токен внутри хендлера
func (h *Handlers) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	// Проверяем авторизацию
	claims, err := auth.GetUserFromRequest(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized: "+err.Error())
		return
	}

	// Используем userID из токена
	tasks, err := h.store.GetAll(claims.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get tasks")
		return
	}

	respondWithJSON(w, http.StatusOK, tasks)
}

// GetTask - проверяем токен внутри хендлера
func (h *Handlers) GetTask(w http.ResponseWriter, r *http.Request) {
	// Проверяем авторизацию
	claims, err := auth.GetUserFromRequest(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized: "+err.Error())
		return
	}

	// Извлекаем ID задачи из URL
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/tasks/"), "/")
	idStr := pathParts[0]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	task, err := h.store.GetByID(id, claims.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, task)
}

// CreateTask - проверяем токен внутри хендлера
func (h *Handlers) CreateTask(w http.ResponseWriter, r *http.Request) {
	// Проверяем авторизацию
	claims, err := auth.GetUserFromRequest(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized: "+err.Error())
		return
	}

	var input models.CreateTaskInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request data")
		return
	}

	if strings.TrimSpace(input.Title) == "" {
		respondWithError(w, http.StatusBadRequest, "Task title is required")
		return
	}

	task, err := h.store.Create(input, claims.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, task)
}

// UpdateTask - проверяем токен внутри хендлера
func (h *Handlers) UpdateTask(w http.ResponseWriter, r *http.Request) {
	// Проверяем авторизацию
	claims, err := auth.GetUserFromRequest(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized: "+err.Error())
		return
	}

	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/tasks/"), "/")
	idStr := pathParts[0]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	var input models.UpdateTaskInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if input.Title != nil && strings.TrimSpace(*input.Title) == "" {
		respondWithError(w, http.StatusBadRequest, "Task title is required")
		return
	}

	task, err := h.store.Update(id, claims.UserID, input)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondWithError(w, http.StatusNotFound, err.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, task)
}

// DeleteTask - проверяем токен внутри хендлера
func (h *Handlers) DeleteTask(w http.ResponseWriter, r *http.Request) {
	// Проверяем авторизацию
	claims, err := auth.GetUserFromRequest(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized: "+err.Error())
		return
	}

	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/tasks/"), "/")
	idStr := pathParts[0]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	err = h.store.Delete(id, claims.UserID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondWithError(w, http.StatusNotFound, err.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Task deleted successfully"})
}
