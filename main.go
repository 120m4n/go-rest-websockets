package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"rest-ws/handlers"
	"rest-ws/middleware"
	"rest-ws/server"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("error loading .env file")
	}

	PORT :=  os.Getenv("PORT")
	if PORT == "" {
		log.Fatalf("PORT is required")
	}
	JWT_SECRET := os.Getenv("JWT_SECRET")
	if JWT_SECRET == "" {
		log.Fatalf("JWT_SECRET is required")
	}
	DATABASE_URL := os.Getenv("DATABASE_URL")
	if DATABASE_URL == "" {
		log.Fatalf("DATABASE_URL is required")
	}

	s, err := server.NewServer(context.Background(), &server.Config{
		Port:        PORT,
		JWTSecret:   JWT_SECRET,
		DatabaseUrl: DATABASE_URL,
	})
	if err != nil {
		log.Fatalf("error creating server: %v", err)
	}

	s.Start(BindRoutes)
	
}

func BindRoutes(s server.Server, r *mux.Router){

	api := r.PathPrefix("/api/v1").Subrouter()
	api.Use(middleware.CheckAuthMiddleware(s))
	// User Handlers
	r.HandleFunc("/", handlers.HomeHandler(s)).Methods(http.MethodGet)
	r.HandleFunc("/health", handlers.HealthHandler(s)).Methods(http.MethodGet)
	r.HandleFunc("/signup", handlers.SignupHandler(s)).Methods(http.MethodPost)
	r.HandleFunc("/login", handlers.LoginHandler(s)).Methods(http.MethodPost)
	r.HandleFunc("/me", handlers.MeHandlers(s)).Methods(http.MethodGet)
	// Post Handlers
	api.HandleFunc("/posts", handlers.InsertPostHandler(s)).Methods(http.MethodPost)
	api.HandleFunc("/posts/{id}", handlers.GetPostByIdHandler(s)).Methods(http.MethodGet)
	api.HandleFunc("/posts/{id}", handlers.UpdatePostHandler(s)).Methods(http.MethodPut)
	api.HandleFunc("/posts/{id}", handlers.DeletePostHandler(s)).Methods(http.MethodDelete)
	api.HandleFunc("/posts", handlers.ListPostsHandler(s)).Methods(http.MethodGet)
	//ws
	r.HandleFunc("/ws", s.Hub().HandleWebSocket)

}