// health.go
package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok"}`))
}

func RegisterHealthCheck(router *mux.Router) {
	router.HandleFunc("/api/v1/health", HealthHandler).Methods("GET")
}
