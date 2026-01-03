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


func (h* TransactionHandler) Deposit(ctx context.Context, req *pb.TransactionRequest) (*pb.TransactinoResponse, error) {
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
func (h *TransactionHandler) WithDraw(ctx context.Context, req *pb.TransactionRequest) (*pb.TransactionRequest, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	response, err := h.transactionService.WithDraw(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return response, nil
}	