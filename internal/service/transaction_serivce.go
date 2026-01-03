package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/paudelanil/grpc-crud/internal/repository"
	"github.com/paudelanil/grpc-crud/pb"
	"github.com/paudelanil/grpc-crud/models"
	"time"
)

type ITransactionService interface {
	Deposit(ctx context.Context, req *pb.TransactionRequest) (*pb.TransactinoResponse, error)
	WithDraw(ctx context.Context, req *pb.TransactionRequest) (*pb.TransactionRequest, error)
}

type  TransactionService struct {
	transactionRepo  repository.ITransactionRepository
	accountRepo repository.IAccountRepository
}

func NewTransactionService(transactionRepo repository.ITransactionRepository, accountRepo repository.IAccountRepository) ITransactionService { 

	return &TransactionService{
		transactionRepo: transactionRepo,
		accountRepo: accountRepo,
	}
}

func (s *TransactionService) Deposit(ctx context.Context, req *pb.TransactionRequest) (*pb.TransactinoResponse, error) {

	refID := req.ReferenceId
	if refID == "" {
		refID = uuid.New().String()
	}

	// begin transaxtion for atomicity

	tx, err := s.transactionRepo.BeginTransaction(ctx)

	if err != nil {
		return nil, err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Validate amount
	if req.Amount <= 0 {
		tx.Rollback()
		return nil, errors.New("amount must be greater than zero")
	}

	accountResponse, err := s.accountRepo.FindByID(ctx, req.AccountId); 
	
	if err != nil {
		tx.Rollback()
		return nil, errors.New("account not found")
	}

	// Check for duplicate reference_id (idempotency)
	exists, err := s.transactionRepo.IsJournalExists(ctx, refID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if exists {
		tx.Rollback()
		return nil, errors.New("duplicate reference_id")
	}

	// Create journal entry
	journal := models.Journal{
		ID:          uuid.New().String(),
		ReferenceID: refID,
		Narration:   req.Narration,
		ValueDate:   time.Now(),
		Status:      "posted",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.transactionRepo.CreateJournal(ctx, &journal); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create journal entries (debit and credit)


	creditEntry := models.JournalEntry{
		ID:        uuid.New().String(),
		JournalID: journal.ID,
		AccountID: accountResponse.ID,
		EntryType: "CREDIT",
		Amount:    req.Amount,
		CreatedAt: time.Now(),
	}


	if err := s.transactionRepo.CreateJournalEntry(ctx, &creditEntry); err != nil {
		tx.Rollback()
		return nil, err
	}
		

	// Update account balance
	newBalance := accountResponse.Balance + req.Amount
	
	if err := s.accountRepo.UpdateBalance(ctx,accountResponse.ID ,newBalance); err != nil {
		tx.Rollback()
		return nil, err
	}
	// Commit transaction
	if err := tx.Commit().Error;  err != nil {
		return nil, err
	}


	return &pb.TransactinoResponse{
		JournalId: journal.ID,
		Message: "Deposit Success",
		NewBalance: newBalance,
		


	}, nil
}

func (s *TransactionService) WithDraw(ctx context.Context, req *pb.TransactionRequest) (*pb.TransactionRequest, error){
	// Implement withdraw logic here
	return &pb.TransactionRequest{}, nil
}