package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/alexhool2/TimeCard/service/event"
	"github.com/alexhool2/TimeCard/service/users"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},

		AllowCredentials: true,
		Debug:            true,
	})

	router := mux.NewRouter()

	// Subrouter para "/api/v1"
	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	// Subrouter para usuários
	userSubRouter := apiRouter.PathPrefix("/users").Subrouter()
	userStore := users.NewStore(s.db)
	EventForUsers := event.NewStore(s.db)
	userHandler := users.NewHandler(userStore, EventForUsers)
	userHandler.RegisterRoutes(userSubRouter)

	// Subrouter para eventos
	eventSubRouter := apiRouter.PathPrefix("/event").Subrouter()
	eventStore := event.NewStore(s.db)
	eventHandler := event.NewHandler(eventStore, userStore)
	eventHandler.RegisterRoutes(eventSubRouter)

	handler := c.Handler(router)

	log.Println("listening on", s.addr)
	return http.ListenAndServe(s.addr, handler)

}
