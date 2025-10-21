package tasks

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Mazin-emad/todo-backend/cmd/service/auth"
	"github.com/Mazin-emad/todo-backend/config"
	"github.com/Mazin-emad/todo-backend/types"
	"github.com/Mazin-emad/todo-backend/utils"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.TaskStore
}

func NewHandler(store types.TaskStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/tasks", h.handleCreateTask).Methods("POST")
	router.HandleFunc("/tasks", h.handleGetTasks).Methods("GET")
	router.HandleFunc("/tasks/{id}", h.handleDeleteTask).Methods("DELETE")
	router.HandleFunc("/tasks/{id}", h.handleUpdateTask).Methods("PUT")
}

func (h *Handler) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	var payload types.CreateTaskPayload
	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate the payload
	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %s", err.Error()))
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
		return
	}

	// Extract token from "Bearer <token>" format
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid authorization header format"))
		return
	}
	
	token := authHeader[7:] // Remove "Bearer " prefix
	userID, err := auth.GetUserIDFromToken([]byte(config.ConfigAmigoo.JWTKey), token)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
		return
	}

	task := &types.Task{
		Title: payload.Title,
		Completed: false,
		UserID: userID,
	}

	err = h.store.CreateTask(task)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJson(w, http.StatusCreated, task)
}


func (h *Handler) handleGetTasks(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
		return
	}
	
	token := authHeader[7:] // Remove "Bearer " prefix
	userID, err := auth.GetUserIDFromToken([]byte(config.ConfigAmigoo.JWTKey), token)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	tasks, err := h.store.GetTasks(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJson(w, http.StatusOK, tasks)
}


func (h *Handler) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	
	taskID, err := strconv.Atoi(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid task ID"))
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
		return
	}
	
	// Extract token from "Bearer <token>" format
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid authorization header format"))
		return
	}
	
	token := authHeader[7:] // Remove "Bearer " prefix
	userID, err := auth.GetUserIDFromToken([]byte(config.ConfigAmigoo.JWTKey), token)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	// Check if the task belongs to the user
	task, err := h.store.GetTaskByID(taskID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("task not found"))
		return
	}

	if task.UserID != userID {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
		return
	}

	err = h.store.DeleteTask(taskID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJson(w, http.StatusOK, map[string]string{"message": "Task deleted successfully"})
}

func (h *Handler) handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	
	// Convert string ID to integer
	taskID, err := strconv.Atoi(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid task ID"))
		return
	}
	
	var payload types.UpdateTaskPayload
	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate the payload
	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %s", err.Error()))
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
		return
	}
	
	// Extract token from "Bearer <token>" format
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid authorization header format"))
		return
	}
	
	token := authHeader[7:] // Remove "Bearer " prefix
	userID, err := auth.GetUserIDFromToken([]byte(config.ConfigAmigoo.JWTKey), token)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	task := &types.Task{
		ID: taskID,
		Title: payload.Title,
		Completed: payload.Completed,
		UserID: userID,
	}

	err = h.store.UpdateTask(task)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJson(w, http.StatusOK, map[string]string{"message": "Task updated successfully"})
}











