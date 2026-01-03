package middleware

import (
	"context"
	"strings"

	"github.com/paudelanil/grpc-crud/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// AuthInterceptor validates JWT tokens for protected endpoints
func AuthInterceptor(authService service.IAuthService) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		// Skip authentication for public methods
		if isPublicMethod(info.FullMethod) {
			return handler(ctx, req)
		}

		// Extract metadata from incoming context
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		// Get Authorization header
		authHeaders := md.Get("authorization")
		if len(authHeaders) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization token")
		}

		// Expect "Bearer <token>"
		token := authHeaders[0]
		if !strings.HasPrefix(token, "Bearer ") {
			return nil, status.Error(
				codes.Unauthenticated,
				"invalid authorization format, expected 'Bearer <token>'",
			)
		}

		token = strings.TrimPrefix(token, "Bearer ")

		// Validate token 
		claims, err := authService.ValidateToken(token)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid or expired token")
		}

		// Add user info to context
		ctx = context.WithValue(ctx, "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "email", claims.Email)
		ctx = context.WithValue(ctx, "username", claims.Username)

		// Continue request
		return handler(ctx, req)
	}
}

// isPublicMethod checks if the gRPC method does not require authentication
func isPublicMethod(method string) bool {
	publicMethods := map[string]bool{
		"/grpc_crud.LoginService/Register": true,
		"/grpc_crud.LoginService/Login":    true,
	}

	return publicMethods[method]
}
