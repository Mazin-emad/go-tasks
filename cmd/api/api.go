package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/Mazin-emad/todo-backend/cmd/service/tasks"
	"github.com/Mazin-emad/todo-backend/cmd/service/user"
	"github.com/gorilla/mux"
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
	router := mux.NewRouter()
	subroutes := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subroutes)

	taskStore := tasks.NewStore(s.db)
	taskHandler := tasks.NewHandler(taskStore)
	taskHandler.RegisterRoutes(subroutes)


	log.Println("Server is running on port", s.addr)
	return http.ListenAndServe(s.addr, router)
}
