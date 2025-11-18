package middleware

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/sessions"
	_ "modernc.org/sqlite"
)

func TestAuthMiddleware(t *testing.T) {
	// Create a test session store
	store := sessions.NewCookieStore([]byte("test-secret-key-min-32-characters-long"))

	// Create a test database (not used in auth middleware, but required for handler)
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	defer db.Close()

	middleware := AuthMiddleware(store)

	t.Run("unauthenticated request returns 401", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/tasks", nil)
		rr := httptest.NewRecorder()

		handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusUnauthorized {
			t.Errorf("AuthMiddleware() status = %v, want %v", status, http.StatusUnauthorized)
		}
	})

	t.Run("authenticated request passes through", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/tasks", nil)
		rr := httptest.NewRecorder()

		// Create a session and set userID (note: key must match what auth.go expects)
		session, _ := store.Get(req, "session")
		session.Values["userID"] = 1
		session.Values["userType"] = "regular"
		session.Save(req, rr)

		// Create new request with session cookie
		req = httptest.NewRequest(http.MethodGet, "/api/tasks", nil)
		for _, cookie := range rr.Result().Cookies() {
			req.AddCookie(cookie)
		}
		rr = httptest.NewRecorder()

		handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify user_id is in context
			userID, ok := GetUserID(r.Context())
			if !ok || userID == 0 {
				t.Error("UserID not found in request context")
			}
			w.WriteHeader(http.StatusOK)
		}))

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("AuthMiddleware() status = %v, want %v", status, http.StatusOK)
		}
	})
}
