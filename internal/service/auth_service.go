package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/paudelanil/grpc-crud/internal/repository"
	"github.com/paudelanil/grpc-crud/models"
	"github.com/paudelanil/grpc-crud/pb"
	"golang.org/x/crypto/bcrypt"
)

// IAuthService defines the interface for authentication operations
type IAuthService interface {
	Login(ctx context.Context, req *pb.UserLoginRequest) (*pb.UserLoginResponse, error)
	Logout(ctx context.Context, req *pb.UserLogoutRequest) (*pb.UserLogoutResponse, error)
	RefreshToken(ctx context.Context, req *pb.TokenRequest) (*pb.TokenResponse, error)
	Register(ctx context.Context, username, email, password string) error
	ValidateToken(tokenString string) (*Claims, error)
}

// AuthService implements IAuthService interface
type AuthService struct {
	userRepo  repository.IUserRepository
	jwtSecret string
}

// Claims represents JWT claims
type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(userRepo repository.IUserRepository, jwtSecret string) IAuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

// Login authenticates a user and returns JWT tokens
func (s *AuthService) Login(
	ctx context.Context,
	req *pb.UserLoginRequest,
) (*pb.UserLoginResponse, error) {
	// Validate input
	if req.Username == "" || req.Password == "" {
		return nil, errors.New("username and password are required")
	}

	// Find user by username
	user, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("user account is inactive")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Generate access token (expires in 15 minutes)
	accessToken, err := s.generateToken(user, 15*time.Minute)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	// Generate refresh token (expires in 7 days)
	refreshToken, err := s.generateToken(user, 7*24*time.Hour)
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	return &pb.UserLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Message:      "Login successful",
	}, nil
}

// Logout handles user logout (in a real app, you'd blacklist the token)
func (s *AuthService) Logout(
	ctx context.Context,
	req *pb.UserLogoutRequest,
) (*pb.UserLogoutResponse, error) {
	// Validate the token
	if req.AccessToken == "" {
		return nil, errors.New("access token is required")
	}

	// In a production system, you would:
	// 1. Validate the token
	// 2. Add it to a blacklist/redis cache
	// 3. Set expiry on the blacklist entry

	return &pb.UserLogoutResponse{
		Message: "Logout successful",
	}, nil
}

// RefreshToken generates a new access token using a refresh token
func (s *AuthService) RefreshToken(
	ctx context.Context,
	req *pb.TokenRequest,
) (*pb.TokenResponse, error) {
	if req.RefreshToken == "" {
		return nil, errors.New("refresh token is required")
	}

	// Parse and validate the refresh token
	claims, err := s.ValidateToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	// Get user from database
	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if !user.IsActive {
		return nil, errors.New("user account is inactive")
	}

	// Generate new access token
	accessToken, err := s.generateToken(user, 15*time.Minute)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	// Generate new refresh token
	refreshToken, err := s.generateToken(user, 7*24*time.Hour)
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	return &pb.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Register creates a new user account
func (s *AuthService) Register(ctx context.Context, username, email, password string) error {
	// Validate input
	if username == "" || email == "" || password == "" {
		return errors.New("username, email, and password are required")
	}

	// Check if username is taken
	taken, err := s.userRepo.IsUsernameTaken(ctx, username)
	if err != nil {
		return err
	}
	if taken {
		return errors.New("username is already taken")
	}

	// Check if email is taken
	taken, err = s.userRepo.IsEmailTaken(ctx, email)
	if err != nil {
		return err
	}
	if taken {
		return errors.New("email is already taken")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	// Create user
	user := &models.User{
		ID:        uuid.New().String(),
		Username:  username,
		Email:     email,
		Password:  string(hashedPassword),
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s.userRepo.Create(ctx, user)
}

// ValidateToken validates a JWT token and returns the claims
func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

// generateToken generates a JWT token for a user
func (s *AuthService) generateToken(user *models.User, duration time.Duration) (string, error) {
	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
