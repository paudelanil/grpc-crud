package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/paudelanil/grpc-crud/models"
	pb "github.com/paudelanil/grpc-crud/pb"
	"gorm.io/gorm"
	"time"
)

// AccountService implements the account management service
type AccountService struct {
	pb.UnimplementedAccountServiceServer
	DB *gorm.DB
}

// CreateUser creates a new user account
func (s *AccountService) CreateUser(ctx context.Context, req *pb.CreateCustomerRequest) (*pb.CreateCustomerResponse, error) {
	customer := models.Customer{
		ID:        uuid.New().String(),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Address:   req.Address,
		Phone:     req.PhoneNumber,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result := s.DB.Create(&customer)
	if result.Error != nil {
		return nil, result.Error
	}

	return &pb.CreateCustomerResponse{
		CustomerId: customer.ID,
		Message:    "Successfully Craeted",
	}, nil
}
