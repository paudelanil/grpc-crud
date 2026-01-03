package handler

import (
	"context"

	"github.com/paudelanil/grpc-crud/internal/service"
	"github.com/paudelanil/grpc-crud/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AuthHandler handles authentication gRPC requests
type AuthHandler struct {
	pb.UnimplementedLoginServiceServer
	authService service.IAuthService
}

// NewAuthHandler creates a new instance of AuthHandler
func NewAuthHandler(authService service.IAuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login handles user login requests
func (h *AuthHandler) Login(ctx context.Context, req *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {
	// Validate input
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	if req.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "username is required")
	}

	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	// Call service layer
	response, err := h.authService.Login(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return response, nil
}

// Register handles user registration requests
func (h *AuthHandler) Register(ctx context.Context, req *pb.UserRegisterRequest) (*pb.UserRegisterResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	if req.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "username is required")
	}

	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	// Call service layer
	err := h.authService.Register(ctx, req.Username, req.Email, req.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.UserRegisterResponse{Message: "User Registered Successfully"}, nil
}

// Logout handles user logout requests
func (h *AuthHandler) Logout(ctx context.Context, req *pb.UserLogoutRequest) (*pb.UserLogoutResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	if req.AccessToken == "" {
		return nil, status.Error(codes.InvalidArgument, "access token is required")
	}

	// Call service layer
	response, err := h.authService.Logout(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return response, nil
}

// RefreshToken handles token refresh requests
func (h *AuthHandler) RefreshToken(ctx context.Context, req *pb.TokenRequest) (*pb.TokenResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	if req.RefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh token is required")
	}

	// Call service layer
	response, err := h.authService.RefreshToken(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return response, nil
}
