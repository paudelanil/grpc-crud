package middleware

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserContext holds user information from the JWT token
type UserContext struct {
	UserID   string
	Username string
	Email    string
}

// GetUserFromContext extracts user information from the context
func GetUserFromContext(ctx context.Context) (*UserContext, error) {
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	username, _ := ctx.Value("username").(string)
	email, _ := ctx.Value("email").(string)

	return &UserContext{
		UserID:   userID,
		Username: username,
		Email:    email,
	}, nil
}
