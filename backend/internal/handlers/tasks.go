package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/kapilsinghsisodiya/task-manager/internal/database"
	appMiddleware "github.com/kapilsinghsisodiya/task-manager/internal/middleware"
	"github.com/kapilsinghsisodiya/task-manager/internal/models"
)

type CreateTaskRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`   // todo, in_progress, completed
	Priority    string    `json:"priority"` // low, medium, high
	DueDate     time.Time `json:"due_date"`
}

// CreateTaskHandler handles POST /api/tasks
func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 1. Extract the authenticated user ID safely from the request context
	userID, ok := appMiddleware.GetUserIDFromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Unauthorized context access"})
		return
	}

	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid request payload"})
		return
	}

	// 2. Add strict validation rules
	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Task title is required"})
		return
	}

	// Set fallbacks for status and priority if they are empty strings
	if req.Status == "" {
		req.Status = "todo"
	}
	if req.Priority == "" {
		req.Priority = "medium"
	}

	// Validate allowed status configurations
	statusLower := strings.ToLower(req.Status)
	if statusLower != "todo" && statusLower != "in_progress" && statusLower != "completed" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Status must be either 'todo', 'in_progress', or 'completed'"})
		return
	}

	// Validate allowed priority levels
	priorityLower := strings.ToLower(req.Priority)
	if priorityLower != "low" && priorityLower != "medium" && priorityLower != "high" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Priority must be either 'low', 'medium', or 'high'"})
		return
	}

	var task models.Task
	query := `
		INSERT INTO tasks (user_id, title, description, status, priority, due_date)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, title, description, status, priority, due_date, created_at, updated_at
	`

	// 3. Persist securely into PostgreSQL database
	err := database.DB.QueryRow(
		query,
		userID,
		req.Title,
		req.Description,
		statusLower,
		priorityLower,
		req.DueDate,
	).Scan(
		&task.ID,
		&task.UserID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.DueDate,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to create task record: " + err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

// ListTasksHandler handles GET /api/tasks (With pagination and status filtering)
func ListTasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, ok := appMiddleware.GetUserIDFromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Unauthorized"})
		return
	}

	// 1. Read Query Parameters for Filtering & Pagination
	statusFilter := r.URL.Query().Get("status")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	// Set defaults for pagination
	limit := 10
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// 2. Build SQL Query Dynamically based on filters
	baseQuery := `SELECT id, user_id, title, description, status, priority, due_date, created_at, updated_at 
	              FROM tasks WHERE user_id = $1`

	var args []interface{}
	args = append(args, userID)
	argCount := 2

	if statusFilter != "" {
		baseQuery += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, strings.ToLower(statusFilter))
		argCount++
	}

	// Append sorting (default newest first) and pagination
	baseQuery += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, limit, offset)

	// 3. Execute Query
	rows, err := database.DB.Query(baseQuery, args...)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to fetch tasks: " + err.Error()})
		return
	}
	defer rows.Close()

	tasks := []models.Task{}
	for rows.Next() {
		var t models.Task
		err := rows.Scan(&t.ID, &t.UserID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.DueDate, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Error parsing task records"})
			return
		}
		tasks = append(tasks, t)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}
