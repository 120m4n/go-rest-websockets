package handlers

import (
	"encoding/json"
	"net/http"
	"rest-ws/server"
)

type HealthResponse struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
}

func HealthHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(HealthResponse{
			Message: "API is healthy",
			Status:  true,
		})
	}
}