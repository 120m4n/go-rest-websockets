package handlers

import (
	"encoding/json"
	"net/http"
	"rest-ws/models"
	"rest-ws/server"
	"rest-ws/repository"

	"github.com/segmentio/ksuid"
)

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupResponse struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
	Id      int64  `json:"id"`
	Email   string `json:"email"`
}

func SignupHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SignupRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := ksuid.NewRandom()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user := &models.User{
			Email:    req.Email,
			Password: req.Password,
			ID:       int64(id.Time().Unix()),
		}

		err = repository.InsertUser(r.Context(), user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(SignupResponse{
			Message: "User created successfully",
			Status:  true,
			Id:      user.ID,
			Email:   user.Email,
		})
	}
}