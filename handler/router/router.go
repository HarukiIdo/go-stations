package router

import (
	"database/sql"
	"net/http"

	"github.com/TechBowl-japan/go-stations/handler"
	"github.com/TechBowl-japan/go-stations/service"
)

func NewRouter(todoDB *sql.DB) *http.ServeMux {
	// register routes
	mux := http.NewServeMux()

	// router or http://localhost:8080/healthz
	mux.Handle("/healthz", handler.NewHealthzHandler())

	// router or http://localhost:8080/todos
	todo := service.NewTODOService(todoDB)

	mux.Handle("/todos", handler.NewTODOHandler(todo))

	return mux
}
