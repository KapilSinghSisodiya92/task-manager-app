package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/kapilsinghsisodiya/task-manager/internal/database"
	"github.com/kapilsinghsisodiya/task-manager/internal/handlers"
	customMiddleware "github.com/kapilsinghsisodiya/task-manager/internal/middleware"
)

func main() {
	// Explicitly check if loading the .env file fails
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: No .env file found, relying on system environment variables")
	}

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

	// Protected Task Routes
	r.Group(func(protectedRouter chi.Router) {
		// Apply our secure authentication middleware gatekeeper
		protectedRouter.Use(customMiddleware.AuthMiddleware)

		// Map the explicit endpoints clearly without nested sub-route conflicts
		protectedRouter.Post("/api/tasks", handlers.CreateTaskHandler)
		protectedRouter.Get("/api/tasks", handlers.ListTasksHandler)
		protectedRouter.Get("/api/tasks/{id}", handlers.GetTaskHandler)
		protectedRouter.Patch("/api/tasks/{id}", handlers.UpdateTaskHandler)
		protectedRouter.Delete("/api/tasks/{id}", handlers.DeleteTaskHandler)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server smoothly running on port %s...\n", port)
	http.ListenAndServe(":"+port, r)
}
