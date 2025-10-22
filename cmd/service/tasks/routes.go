package tasks

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Mazin-emad/todo-backend/cmd/service/auth"
	"github.com/Mazin-emad/todo-backend/config"
	"github.com/Mazin-emad/todo-backend/types"
	"github.com/Mazin-emad/todo-backend/utils"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	store types.TaskStore
}

func NewHandler(store types.TaskStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/tasks", h.handleCreateTask)
	router.GET("/tasks", h.handleGetTasks)
	router.DELETE("/tasks/:id", h.handleDeleteTask)
	router.PUT("/tasks/:id", h.handleUpdateTask)
}

func (h *Handler) handleCreateTask(c *gin.Context) {
	var payload types.CreateTaskPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.WriteErrorGin(c, http.StatusBadRequest, err)
		return
	}

	// validate the payload
	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteErrorGin(c, http.StatusBadRequest, fmt.Errorf("invalid payload: %s", err.Error()))
		return
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		utils.WriteErrorGin(c, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
		return
	}
	
	// Extract token from "Bearer <token>" format
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		utils.WriteErrorGin(c, http.StatusUnauthorized, fmt.Errorf("invalid authorization header format"))
		return
	}
	
	token := authHeader[7:] // Remove "Bearer " prefix
	userID, err := auth.GetUserIDFromToken([]byte(config.ConfigAmigoo.JWTKey), token)
	if err != nil {
		utils.WriteErrorGin(c, http.StatusUnauthorized, err)
		return
	}

	task := &types.Task{
		Title: payload.Title,
		Completed: false,
		UserID: userID,
	}

	err = h.store.CreateTask(task)
	if err != nil {
		utils.WriteErrorGin(c, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJsonGin(c, http.StatusCreated, map[string]string{"message": "Task created successfully"})
}

func (h *Handler) handleGetTasks(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		utils.WriteErrorGin(c, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
		return
	}
	
	// Extract token from "Bearer <token>" format
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		utils.WriteErrorGin(c, http.StatusUnauthorized, fmt.Errorf("invalid authorization header format"))
		return
	}
	
	token := authHeader[7:] // Remove "Bearer " prefix
	userID, err := auth.GetUserIDFromToken([]byte(config.ConfigAmigoo.JWTKey), token)
	if err != nil {
		utils.WriteErrorGin(c, http.StatusUnauthorized, err)
		return
	}

	tasks, err := h.store.GetTasks(userID)
	if err != nil {
		utils.WriteErrorGin(c, http.StatusInternalServerError, err)
		return
	}
	
	utils.WriteJsonGin(c, http.StatusOK, map[string]any{"tasks": tasks})
}

func (h *Handler) handleDeleteTask(c *gin.Context) {
	id := c.Param("id")
	
	// Convert string ID to integer
	taskID, err := strconv.Atoi(id)
	if err != nil {
		utils.WriteErrorGin(c, http.StatusBadRequest, fmt.Errorf("invalid task ID"))
		return
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		utils.WriteErrorGin(c, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
		return
	}
	
	// Extract token from "Bearer <token>" format
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		utils.WriteErrorGin(c, http.StatusUnauthorized, fmt.Errorf("invalid authorization header format"))
		return
	}
	
	token := authHeader[7:] // Remove "Bearer " prefix
	userID, err := auth.GetUserIDFromToken([]byte(config.ConfigAmigoo.JWTKey), token)
	if err != nil {
		utils.WriteErrorGin(c, http.StatusUnauthorized, err)
		return
	}

	// Check if the task belongs to the user
	task, err := h.store.GetTaskByID(taskID)
	if err != nil {
		utils.WriteErrorGin(c, http.StatusNotFound, fmt.Errorf("task not found"))
		return
	}

	if task.UserID != userID {
		utils.WriteErrorGin(c, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
		return
	}

	err = h.store.DeleteTask(taskID)
	if err != nil {
		utils.WriteErrorGin(c, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJsonGin(c, http.StatusOK, map[string]string{"message": "Task deleted successfully"})
}

func (h *Handler) handleUpdateTask(c *gin.Context) {
	id := c.Param("id")
	
	// Convert string ID to integer
	taskID, err := strconv.Atoi(id)
	if err != nil {
		utils.WriteErrorGin(c, http.StatusBadRequest, fmt.Errorf("invalid task ID"))
		return
	}
	
	var payload types.UpdateTaskPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.WriteErrorGin(c, http.StatusBadRequest, err)
		return
	}

	// validate the payload
	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteErrorGin(c, http.StatusBadRequest, fmt.Errorf("invalid payload: %s", err.Error()))
		return
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		utils.WriteErrorGin(c, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
		return
	}
	
	// Extract token from "Bearer <token>" format
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		utils.WriteErrorGin(c, http.StatusUnauthorized, fmt.Errorf("invalid authorization header format"))
		return
	}
	
	token := authHeader[7:] // Remove "Bearer " prefix
	userID, err := auth.GetUserIDFromToken([]byte(config.ConfigAmigoo.JWTKey), token)
	if err != nil {
		utils.WriteErrorGin(c, http.StatusUnauthorized, err)
		return
	}

	// Check if the task belongs to the user
	task, err := h.store.GetTaskByID(taskID)
	if err != nil {
		utils.WriteErrorGin(c, http.StatusNotFound, fmt.Errorf("task not found"))
		return
	}

	if task.UserID != userID {
		utils.WriteErrorGin(c, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
		return
	}

	task.Title = payload.Title
	task.Completed = payload.Completed

	err = h.store.UpdateTask(task)
	if err != nil {
		utils.WriteErrorGin(c, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJsonGin(c, http.StatusOK, map[string]string{"message": "Task updated successfully"})
}