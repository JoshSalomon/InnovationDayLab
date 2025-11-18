package handlers

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "modernc.org/sqlite"
)

func TestHealthHandler_GetHealth(t *testing.T) {
	// Create in-memory database
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	defer db.Close()

	handler := NewHealthHandler(db)

	t.Run("health check returns 200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
		rr := httptest.NewRecorder()

		handler.GetHealth(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("GetHealth() status = %v, want %v", status, http.StatusOK)
		}

		// Check response body contains expected content
		body := rr.Body.String()
		if body == "" {
			t.Error("GetHealth() returned empty body")
		}
	})

	t.Run("health check with database connection", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
		rr := httptest.NewRecorder()

		handler.GetHealth(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("GetHealth() status = %v, want %v", status, http.StatusOK)
		}
	})
}
