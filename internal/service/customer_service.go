package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/paudelanil/grpc-crud/internal/repository"
	"github.com/paudelanil/grpc-crud/models"
	"github.com/paudelanil/grpc-crud/pb"
)

// ICustomerService defines the interface for customer operations
type ICustomerService interface {
	CreateCustomer(ctx context.Context, req *pb.CreateCustomerRequest) (*pb.CreateCustomerResponse, error)
	GetCustomer(ctx context.Context, req *pb.GetCustomerRequest) (*pb.GetCustomerResponse, error)
	UpdateCustomer(ctx context.Context, req *pb.UpdateCustomerRequest) (*pb.UpdateCustomerResponse, error)
	DeleteCustomer(ctx context.Context, req *pb.DeleteCustomerRequest) (*pb.DeleteCustomerResponse, error)
	ListCustomers(ctx context.Context, req *pb.ListCustomerRequest) (*pb.ListCustomerResponse, error)
}

// CustomerService implements ICustomerService interface
type CustomerService struct {
	customerRepo repository.ICustomerRepository
}

// NewCustomerService creates a new instance of CustomerService
func NewCustomerService(customerRepo repository.ICustomerRepository) ICustomerService {
	return &CustomerService{
		customerRepo: customerRepo,
	}
}

// CreateCustomer creates a new customer
func (s *CustomerService) CreateCustomer(ctx context.Context, req *pb.CreateCustomerRequest) (*pb.CreateCustomerResponse, error) {
	// Validate input
	if req.FirstName == "" || req.LastName == "" {
		return nil, errors.New("first name and last name are required")
	}

	if req.Email == "" {
		return nil, errors.New("email is required")
	}

	if req.PhoneNumber == "" {
		return nil, errors.New("phone number is required")
	}

	// Check if email is already taken
	emailTaken, err := s.customerRepo.IsEmailTaken(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if emailTaken {
		return nil, errors.New("email is already in use")
	}

	// Check if phone is already taken
	phoneTaken, err := s.customerRepo.IsPhoneTaken(ctx, req.PhoneNumber)
	if err != nil {
		return nil, err
	}
	if phoneTaken {
		return nil, errors.New("phone number is already in use")
	}

	// Create customer model
	customer := &models.Customer{
		ID:        uuid.New().String(),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.PhoneNumber,
		Address:   req.Address,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save to database
	if err := s.customerRepo.Create(ctx, customer); err != nil {
		return nil, errors.New("failed to create customer")
	}

	return &pb.CreateCustomerResponse{
		CustomerId: customer.ID,
		Message:    "Customer created successfully",
	}, nil
}

// GetCustomer retrieves a customer by ID
func (s *CustomerService) GetCustomer(ctx context.Context, req *pb.GetCustomerRequest) (*pb.GetCustomerResponse, error) {
	if req.CustomerId == "" {
		return nil, errors.New("customer ID is required")
	}

	customer, err := s.customerRepo.FindByID(ctx, req.CustomerId)
	if err != nil {
		return nil, err
	}

	return &pb.GetCustomerResponse{
		CustomerId:  customer.ID,
		FirstName:   customer.FirstName,
		LastName:    customer.LastName,
		Email:       customer.Email,
		PhoneNumber: customer.Phone,
		Address:     customer.Address,
		CreatedAt:   customer.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   customer.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// UpdateCustomer updates an existing customer
func (s *CustomerService) UpdateCustomer(ctx context.Context, req *pb.UpdateCustomerRequest) (*pb.UpdateCustomerResponse, error) {
	if req.CustomerId == "" {
		return nil, errors.New("customer ID is required")
	}

	// Find existing customer
	customer, err := s.customerRepo.FindByID(ctx, req.CustomerId)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.FirstName != "" {
		customer.FirstName = req.FirstName
	}
	if req.LastName != "" {
		customer.LastName = req.LastName
	}
	if req.Email != "" {
		customer.Email = req.Email
	}
	if req.PhoneNumber != "" {
		customer.Phone = req.PhoneNumber
	}
	if req.Address != "" {
		customer.Address = req.Address
	}
	customer.UpdatedAt = time.Now()

	// Save changes
	if err := s.customerRepo.Update(ctx, customer); err != nil {
		return nil, errors.New("failed to update customer")
	}

	return &pb.UpdateCustomerResponse{
		Message: "Customer updated successfully",
		Customer: &pb.GetCustomerResponse{
			CustomerId:  customer.ID,
			FirstName:   customer.FirstName,
			LastName:    customer.LastName,
			Email:       customer.Email,
			PhoneNumber: customer.Phone,
			Address:     customer.Address,
			CreatedAt:   customer.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   customer.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

// DeleteCustomer deletes a customer by ID
func (s *CustomerService) DeleteCustomer(ctx context.Context, req *pb.DeleteCustomerRequest) (*pb.DeleteCustomerResponse, error) {
	if req.CustomerId == "" {
		return nil, errors.New("customer ID is required")
	}

	if err := s.customerRepo.Delete(ctx, req.CustomerId); err != nil {
		return nil, err
	}

	return &pb.DeleteCustomerResponse{
		Message: "Customer deleted successfully",
	}, nil
}

// ListCustomers lists all customers with pagination
func (s *CustomerService) ListCustomers(ctx context.Context, req *pb.ListCustomerRequest) (*pb.ListCustomerResponse, error) {
	// Set default pagination values
	pageSize := int(req.PageSize)
	if pageSize <= 0 {
		pageSize = 10
	}

	pageNumber := int(req.PageNumber)
	if pageNumber <= 0 {
		pageNumber = 1
	}

	offset := (pageNumber - 1) * pageSize

	// Fetch customers
	customers, err := s.customerRepo.FindAll(ctx, pageSize, offset)
	if err != nil {
		return nil, errors.New("failed to retrieve customers")
	}

	// Convert to response format
	var customerResponses []*pb.GetCustomerResponse
	for _, customer := range customers {
		customerResponses = append(customerResponses, &pb.GetCustomerResponse{
			CustomerId:  customer.ID,
			FirstName:   customer.FirstName,
			LastName:    customer.LastName,
			Email:       customer.Email,
			PhoneNumber: customer.Phone,
			Address:     customer.Address,
			CreatedAt:   customer.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   customer.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &pb.ListCustomerResponse{
		Customers: customerResponses,
	}, nil
}
