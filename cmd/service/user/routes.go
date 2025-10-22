package user

import (
	"fmt"
	"net/http"

	"github.com/Mazin-emad/todo-backend/cmd/service/auth"
	"github.com/Mazin-emad/todo-backend/config"
	"github.com/Mazin-emad/todo-backend/types"
	"github.com/Mazin-emad/todo-backend/utils"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/login", h.handleLogin)
	router.POST("/register", h.handleRegister)
	router.GET("/current-user", h.handleCurrentUser)
	router.POST("/logout", h.handleLogout)
	router.GET("/users", h.handleGetUsers)
}

func (h *Handler) handleLogin(c *gin.Context) {
	var payload types.LoginUserPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.WriteErrorGin(c, http.StatusBadRequest, err)
		return
	}

	// validate the payload
	if err := utils.Validate.Struct(payload); err != nil {
	utils.WriteErrorGin(c, http.StatusBadRequest, fmt.Errorf("invalid payload: %s", err.Error()))
		return
	}
  
	user, err := h.store.GetUserByUsername(payload.UserName)
	if err != nil {
		// Check if it's a "user not found" error or a database error
		if err.Error() == fmt.Sprintf("no user found with username: %s", payload.UserName) {
			utils.WriteErrorGin(c, http.StatusUnauthorized, fmt.Errorf("invalid username or password"))
		} else {
			utils.WriteErrorGin(c, http.StatusInternalServerError, err)
		}
		return
	}

	if user.Password != payload.Password {
		utils.WriteErrorGin(c, http.StatusUnauthorized, fmt.Errorf("invalid username or password"))
		return
	}

	secret := []byte(config.ConfigAmigoo.JWTKey)
	token, err := auth.GenerateToken(secret, user.ID)
	if err != nil {
		utils.WriteErrorGin(c, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJsonGin(c, http.StatusOK, map[string]any{"token": token, "user": user})
}

func (h *Handler) handleRegister(c *gin.Context) {
	// parse the request body
	var payload types.RegisterUserPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.WriteErrorGin(c, http.StatusBadRequest, err)
		return
	}

	// validate the payload
	if err := utils.Validate.Struct(payload); err != nil {
	utils.WriteErrorGin(c, http.StatusBadRequest, fmt.Errorf("invalid payload: %s", err.Error()))
		return
	}

	_, err := h.store.GetUserByUsername(payload.UserName)
	if err == nil {
		utils.WriteErrorGin(c, http.StatusBadRequest, fmt.Errorf("%s already exists", payload.UserName))
		return
	}

	err = h.store.CreateUser(&types.User{
		UserName: payload.UserName,
		Password: payload.Password,
	})

	if err != nil {
		utils.WriteErrorGin(c, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJsonGin(c, http.StatusCreated, map[string]string{"message": "User created successfully"})
}


func (h *Handler) handleCurrentUser(c *gin.Context) {
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
	user, err := h.store.GetUserByID(userID)
	if err != nil {
		utils.WriteErrorGin(c, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJsonGin(c, http.StatusOK, map[string]any{"user": user})
}

func (h *Handler) handleLogout(c *gin.Context) {
	// For JWT tokens, logout is typically handled client-side by removing the token
	// This endpoint just confirms the logout was successful
	utils.WriteJsonGin(c, http.StatusOK, map[string]string{"message": "Logged out successfully"})
}

func (h *Handler) handleGetUsers(c *gin.Context) {
	users, err := h.store.GetAllUsers()
	if err != nil {
		utils.WriteErrorGin(c, http.StatusInternalServerError, err)
		return
	}
	
	utils.WriteJsonGin(c, http.StatusOK, map[string]any{"users": users})
}


