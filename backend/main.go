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

	// CRITICAL FIX: Setup a robust, browser-compliant CORS configuration.
	// Since we are using an explicit array for AllowCredentials, we whitelist the domains.
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"*", // Allow all origins
			"http://localhost:3000",
			"https://task-manager-app-beta-steel.vercel.app",
		},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Requested-With"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// ROUTER PREFLIGHT FIX: Explicitly short-circuit any manual/implicit OPTIONS requests at the root level
	r.Options("/*", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

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
	// FIX: Pass the chi router instance `r` here instead of `nil` to ensure Go listens to your routes!
	http.ListenAndServe(":"+port, r)
}
