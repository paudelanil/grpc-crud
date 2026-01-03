package handler

import (
	"context"

	"github.com/paudelanil/grpc-crud/internal/service"
	"github.com/paudelanil/grpc-crud/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AccountHandler handles account-related gRPC requests
type AccountHandler struct {
	pb.UnimplementedAccountServiceServer
	customerService service.ICustomerService
	accountService  service.IAccountService
}

// NewAccountHandler creates a new instance of AccountHandler
func NewAccountHandler(customerService service.ICustomerService, accountService service.IAccountService) *AccountHandler {
	return &AccountHandler{
		customerService: customerService,
		accountService:  accountService,
	}
}

// Customer operations

// CreateUser creates a new customer
func (h *AccountHandler) CreateUser(ctx context.Context, req *pb.CreateCustomerRequest) (*pb.CreateCustomerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	response, err := h.customerService.CreateCustomer(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return response, nil
}

// GetUser retrieves a customer by ID
func (h *AccountHandler) GetUser(ctx context.Context, req *pb.GetCustomerRequest) (*pb.GetCustomerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	response, err := h.customerService.GetCustomer(ctx, req)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return response, nil
}

// UpdateUser updates a customer
func (h *AccountHandler) UpdateUser(ctx context.Context, req *pb.UpdateCustomerRequest) (*pb.UpdateCustomerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	response, err := h.customerService.UpdateCustomer(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return response, nil
}

// DeleteUser deletes a customer
func (h *AccountHandler) DeleteUser(ctx context.Context, req *pb.DeleteCustomerRequest) (*pb.DeleteCustomerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	response, err := h.customerService.DeleteCustomer(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return response, nil
}

// ListUsers lists all customers with pagination
func (h *AccountHandler) ListUsers(ctx context.Context, req *pb.ListCustomerRequest) (*pb.ListCustomerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	response, err := h.customerService.ListCustomers(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return response, nil
}

// Account operations

// CreateAccount creates a new account for a customer
func (h *AccountHandler) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	response, err := h.accountService.CreateAccount(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return response, nil
}

// GetAccount retrieves an account by ID
func (h *AccountHandler) GetAccount(ctx context.Context, req *pb.GetAccountRequest) (*pb.GetAccountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	response, err := h.accountService.GetAccount(ctx, req)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return response, nil
}

// UpdateAccount updates an account
func (h *AccountHandler) UpdateAccount(ctx context.Context, req *pb.UpdateAccountRequest) (*pb.UpdateAccountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	response, err := h.accountService.UpdateAccount(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return response, nil
}

// DeleteAccount deletes an account
func (h *AccountHandler) DeleteAccount(ctx context.Context, req *pb.DeleteAccountRequest) (*pb.DeleteAccountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	response, err := h.accountService.DeleteAccount(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return response, nil
}

// ListAccounts lists all accounts with pagination
func (h *AccountHandler) ListAccounts(ctx context.Context, req *pb.ListAccountRequest) (*pb.ListAccountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	response, err := h.accountService.ListAccounts(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return response, nil
}
