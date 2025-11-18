package middleware

import (
	"context"
	"net/http"

	"github.com/gorilla/sessions"
)

type contextKey string

const userIDKey contextKey = "userID"
const userTypeKey contextKey = "userType"

// AuthMiddleware checks if user is authenticated
func AuthMiddleware(store *sessions.CookieStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := store.Get(r, "session")
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			userID, ok := session.Values["userID"].(int)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			userType, _ := session.Values["userType"].(string)

			// Add user info to context
			ctx := context.WithValue(r.Context(), userIDKey, userID)
			ctx = context.WithValue(ctx, userTypeKey, userType)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID extracts user ID from context
func GetUserID(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(userIDKey).(int)
	return userID, ok
}

// GetUserType extracts user type from context
func GetUserType(ctx context.Context) (string, bool) {
	userType, ok := ctx.Value(userTypeKey).(string)
	return userType, ok
}
