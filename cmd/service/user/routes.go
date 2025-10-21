package user

import (
	"fmt"
	"net/http"

	"github.com/Mazin-emad/todo-backend/cmd/service/auth"
	"github.com/Mazin-emad/todo-backend/config"
	"github.com/Mazin-emad/todo-backend/types"
	"github.com/Mazin-emad/todo-backend/utils"

	"github.com/gorilla/mux"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
	router.HandleFunc("/current-user", h.handleCurrentUser).Methods("GET")
	router.HandleFunc("/logout", h.handleLogout).Methods("POST")
	router.HandleFunc("/users", h.handleGetUsers).Methods("GET")
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var payload types.LoginUserPayload

	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate the payload
	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %s", err.Error()))
		return
	}
  
	user, err := h.store.GetUserByUsername(payload.UserName)
	if err != nil {
		// Check if it's a "user not found" error or a database error
		if err.Error() == fmt.Sprintf("no user found with username: %s", payload.UserName) {
			utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid username or password"))
		} else {
			utils.WriteError(w, http.StatusInternalServerError, err)
		}
		return
	}

	if user.Password != payload.Password {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid username or password"))
		return
	}

	secret := []byte(config.ConfigAmigoo.JWTKey)
	token, err := auth.GenerateToken(secret, user.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, map[string]any{"token": token, "user": user})

}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	// parse the request body
	var payload types.RegisterUserPayload

	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate the payload
	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %s", err.Error()))
		return
	}

	_, err := h.store.GetUserByUsername(payload.UserName)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("%s already exists", payload.UserName))
		return
	}

	err = h.store.CreateUser(&types.User{
		UserName: payload.UserName,
		Password: payload.Password,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusCreated, map[string]string{"message": "User created successfully"})

	
}


func (h *Handler) handleCurrentUser(w http.ResponseWriter, r *http.Request) {
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
	user, err := h.store.GetUserByID(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJson(w, http.StatusOK, map[string]any{"user": user})
}

func (h *Handler) handleLogout(w http.ResponseWriter, r *http.Request) {
	// For JWT tokens, logout is typically handled client-side by removing the token
	// This endpoint just confirms the logout was successful
	utils.WriteJson(w, http.StatusOK, map[string]string{"message": "Logged out successfully"})
}

func (h *Handler) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.store.GetAllUsers()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	
	utils.WriteJson(w, http.StatusOK, map[string]any{"users": users})
}


