package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	appMiddleware "github.com/kapilsinghsisodiya/task-manager/internal/middleware"
	"github.com/kapilsinghsisodiya/task-manager/internal/models"
)

// Test 1: Meaningful Unit Test verifying Password Hashing Logic
func TestPasswordHashingAndVerification(t *testing.T) {
	password := "superSecure123"

	hash, err := models.HashPassword(password)
	if err != nil {
		t.Fatalf("Expected password hash generation to pass, got error: %v", err)
	}

	if hash == password {
		t.Errorf("Security flaw: Hash should not match plaintext password string")
	}

	if !models.CheckPasswordHash(password, hash) {
		t.Errorf("Password verification failed for valid credentials matching")
	}

	if models.CheckPasswordHash("wrongPassword", hash) {
		t.Errorf("Security flaw: Verification passed for incorrect password text matching")
	}
}

// Test 2: Meaningful API Guard Test checking Task validation constraints
func TestCreateTaskHandlerMissingTitle(t *testing.T) {
	// Prepare sample empty JSON payload (missing the mandatory Title field)
	jsonPayload := []byte(`{"description": "Testing missing title context", "priority": "high"}`)

	req, err := http.NewRequest("POST", "/api/tasks", bytes.NewBuffer(jsonPayload))
	if err != nil {
		t.Fatalf("Could not construct test request instance: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Inject a mock authorized User Context directly to bypass middleware auth check
	ctx := context.WithValue(req.Context(), appMiddleware.UserIDKey, 1)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateTaskHandler)

	handler.ServeHTTP(rr, req)

	// We expect a 422 Unprocessable Entity code because the task title is missing
	if rr.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status code %d, but received %d instead", http.StatusUnprocessableEntity, rr.Code)
	}
}

// Test 3: Meaningful Context Security Test verifying identity extraction stability
func TestContextRetrievalHelpers(t *testing.T) {
	ctx := context.Background()

	// Verify it returns false when no user context parameters exist
	_, ok := appMiddleware.GetUserIDFromContext(ctx)
	if ok {
		t.Errorf("Expected retrieval helper to return false on empty context structures")
	}

	// Verify it successfully extracts data when correctly populated
	expectedUserID := 42
	ctxWithUser := context.WithValue(ctx, appMiddleware.UserIDKey, expectedUserID)

	val, ok := appMiddleware.GetUserIDFromContext(ctxWithUser)
	if !ok || val != expectedUserID {
		t.Errorf("Expected user ID %d to be retrieved cleanly, got %d (ok: %t)", expectedUserID, val, ok)
	}
}
