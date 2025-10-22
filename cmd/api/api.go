package api

import (
	"database/sql"
	"log"

	"github.com/Mazin-emad/todo-backend/cmd/service/tasks"
	"github.com/Mazin-emad/todo-backend/cmd/service/user"
	"github.com/gin-gonic/gin"
)

type ApiServer struct {
	addr string
	db   *sql.DB
}


func NewApiServer(addr string, db *sql.DB) *ApiServer {
	return &ApiServer{
		addr: addr,
		db: db,
	}
}

func (s *ApiServer) Run() error {
	// Set Gin to release mode for production
	gin.SetMode(gin.ReleaseMode)
	
	router := gin.Default()
	
	// Create API v1 group
	v1 := router.Group("/api/v1")
	
	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(v1)

	taskStore := tasks.NewStore(s.db)
	taskHandler := tasks.NewHandler(taskStore)
	taskHandler.RegisterRoutes(v1)

	log.Println("Server is running on port", s.addr)
	return router.Run(s.addr)
}
