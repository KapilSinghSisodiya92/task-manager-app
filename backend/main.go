package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/yourusername/task-manager/internal/database"
	"github.com/yourusername/task-manager/internal/handlers"
)

func main() {
	// Load environment variables from root or backend directory if present
	_ = godotenv.Load()

	// Initialize Database
	database.InitDB()
	defer database.DB.Close()

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Auth Public Routes
	r.Post("/api/auth/signup", handlers.SignupHandler)
	r.Post("/api/auth/login", handlers.LoginHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server smoothly running on port %s...\n", port)
	http.ListenAndServe(":"+port, r)
}
