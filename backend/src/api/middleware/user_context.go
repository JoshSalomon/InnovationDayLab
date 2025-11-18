package middleware

import (
	"context"
)

// UserContext provides helper functions to extract user information from context
// This file complements auth.go by providing additional context utilities

// GetUserIDFromContext extracts user ID from context (alias for GetUserID)
func GetUserIDFromContext(ctx context.Context) (int, bool) {
	return GetUserID(ctx)
}

// GetUserTypeFromContext extracts user type from context (alias for GetUserType)
func GetUserTypeFromContext(ctx context.Context) (string, bool) {
	return GetUserType(ctx)
}
