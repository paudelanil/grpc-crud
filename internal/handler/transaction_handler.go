package handler

import (
	"context"

	"github.com/paudelanil/grpc-crud/internal/service"
	"github.com/paudelanil/grpc-crud/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


type TransactionHandler struct {
	pb.UnimplementedTransactionServiceServer
	transactionService service.ITransactionService
}

// NewTransactionHandler creates a new instance of TransactionHandler
func NewTransactionHandler(transactionService service.ITransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

// Transaction operations


func (h* TransactionHandler) Deposit(ctx context.Context, req *pb.TransactionRequest) (*pb.TransactionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	response, err := h.transactionService.Deposit(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return response, nil
}

// WithDraw handles withdraw requests
func (h *TransactionHandler) WithDraw(ctx context.Context, req *pb.TransactionRequest) (*pb.TransactionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	response, err := h.transactionService.WithDraw(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return response, nil
}	

// TransferMoney handles transfer money requests
func (h *TransactionHandler) TransferMoney(ctx context.Context, req *pb.TransferMoneyRequest) (*pb.TransferMoneyResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	response, err := h.transactionService.TransferMoney(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return response, nil
}

// GetStatement handles get statement requests
func (h *TransactionHandler) GetStatement(ctx context.Context, req *pb.GetStatementRequest) (*pb.GetStatementResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	response, err := h.transactionService.GetStatement(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return response, nil
}

// GenerateStatement handles generate statement requests
func (h *TransactionHandler) GenerateStatement(ctx context.Context, req *pb.GenerateStatementRequest) (*pb.GenerateStatementResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	response, err := h.transactionService.GenerateStatement(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return response, nil
}

// ListStatements handles list statements requests
func (h *TransactionHandler) ListStatements(ctx context.Context, req *pb.ListStatementsRequest) (*pb.ListStatementsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	response, err := h.transactionService.ListStatements(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return response, nil
}	