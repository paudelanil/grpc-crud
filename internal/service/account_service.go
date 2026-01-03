package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/paudelanil/grpc-crud/internal/repository"
	"github.com/paudelanil/grpc-crud/models"
	"github.com/paudelanil/grpc-crud/pb"
)

// IAccountService defines the interface for account operations
type IAccountService interface {
	CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error)
	GetAccount(ctx context.Context, req *pb.GetAccountRequest) (*pb.GetAccountResponse, error)
	UpdateAccount(ctx context.Context, req *pb.UpdateAccountRequest) (*pb.UpdateAccountResponse, error)
	DeleteAccount(ctx context.Context, req *pb.DeleteAccountRequest) (*pb.DeleteAccountResponse, error)
	ListAccounts(ctx context.Context, req *pb.ListAccountRequest) (*pb.ListAccountResponse, error)
}

// AccountServiceImpl implements IAccountService interface
type AccountServiceImpl struct {
	accountRepo  repository.IAccountRepository
	customerRepo repository.ICustomerRepository
}

// NewAccountService creates a new instance of AccountService
func NewAccountService(accountRepo repository.IAccountRepository, customerRepo repository.ICustomerRepository) IAccountService {
	return &AccountServiceImpl{
		accountRepo:  accountRepo,
		customerRepo: customerRepo,
	}
}

// CreateAccount creates a new account for a customer
func (s *AccountServiceImpl) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	if req.CustomerId == "" {
		return nil, errors.New("customer ID is required")
	}

	// Verify customer exists
	customer, err := s.customerRepo.FindByID(ctx, req.CustomerId)
	if err != nil {
		return nil, errors.New("customer not found")
	}

	// Generate unique account number
	accountNumber := fmt.Sprintf("%s-%d", customer.Phone, time.Now().Unix())

	// Create account model
	account := &models.Account{
		ID:            uuid.New().String(),
		AccountNumber: accountNumber,
		CustomerID:    req.CustomerId,
		Status:        "active",
		Balance:       0.0,
		OpenedAt:      time.Now(),
		Currency:      "NPR",
		AccountType:   "savings",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Save to database
	if err := s.accountRepo.Create(ctx, account); err != nil {
		return nil, errors.New("failed to create account")
	}

	return &pb.CreateAccountResponse{
		AccountId:     account.ID,
		AccountNumber: account.AccountNumber,
		Message:       "Account created successfully",
	}, nil
}

// GetAccount retrieves an account by ID
func (s *AccountServiceImpl) GetAccount(ctx context.Context, req *pb.GetAccountRequest) (*pb.GetAccountResponse, error) {
	if req.AccountId == "" {
		return nil, errors.New("account ID is required")
	}

	account, err := s.accountRepo.FindByID(ctx, req.AccountId)
	if err != nil {
		return nil, err
	}

	return &pb.GetAccountResponse{
		AccountId:     account.ID,
		CustomerId:    account.CustomerID,
		AccountNumber: account.AccountNumber,
		AccountType:   account.AccountType,
		Balance:       account.Balance,
		Currency:      account.Currency,
		Status:        account.Status,
		CreatedAt:     account.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     account.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// UpdateAccount updates an existing account
func (s *AccountServiceImpl) UpdateAccount(ctx context.Context, req *pb.UpdateAccountRequest) (*pb.UpdateAccountResponse, error) {
	if req.AccountId == "" {
		return nil, errors.New("account ID is required")
	}

	// Find existing account
	account, err := s.accountRepo.FindByID(ctx, req.AccountId)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.AccountType != "" {
		account.AccountType = req.AccountType
	}
	if req.Status != "" {
		account.Status = req.Status
	}
	account.UpdatedAt = time.Now()

	// Save changes
	if err := s.accountRepo.Update(ctx, account); err != nil {
		return nil, errors.New("failed to update account")
	}

	return &pb.UpdateAccountResponse{
		Message: "Account updated successfully",
		Account: &pb.GetAccountResponse{
			AccountId:     account.ID,
			AccountNumber: account.AccountNumber,
			CustomerId:    account.CustomerID,
			AccountType:   account.AccountType,
			Balance:       account.Balance,
			Currency:      account.Currency,
			Status:        account.Status,
			CreatedAt:     account.CreatedAt.Format(time.RFC3339),
			UpdatedAt:     account.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

// DeleteAccount deletes an account by ID
func (s *AccountServiceImpl) DeleteAccount(ctx context.Context, req *pb.DeleteAccountRequest) (*pb.DeleteAccountResponse, error) {
	if req.AccountId == "" {
		return nil, errors.New("account ID is required")
	}

	if err := s.accountRepo.Delete(ctx, req.AccountId); err != nil {
		return nil, err
	}

	return &pb.DeleteAccountResponse{
		Message: "Account deleted successfully",
	}, nil
}

// ListAccounts lists all accounts with pagination
func (s *AccountServiceImpl) ListAccounts(ctx context.Context, req *pb.ListAccountRequest) (*pb.ListAccountResponse, error) {
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

	// Fetch accounts
	accounts, err := s.accountRepo.FindAll(ctx, pageSize, offset)
	if err != nil {
		return nil, errors.New("failed to retrieve accounts")
	}

	// Convert to response format
	var accountResponses []*pb.GetAccountResponse
	for _, account := range accounts {
		accountResponses = append(accountResponses, &pb.GetAccountResponse{
			AccountId:     account.ID,
			CustomerId:    account.CustomerID,
			AccountNumber: account.AccountNumber,
			AccountType:   account.AccountType,
			Balance:       account.Balance,
			Currency:      account.Currency,
			Status:        account.Status,
			CreatedAt:     account.CreatedAt.Format(time.RFC3339),
			UpdatedAt:     account.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &pb.ListAccountResponse{
		Accounts: accountResponses,
	}, nil
}
