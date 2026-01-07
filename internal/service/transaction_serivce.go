package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/paudelanil/grpc-crud/internal/repository"
	"github.com/paudelanil/grpc-crud/pb"
	"github.com/paudelanil/grpc-crud/models"
	"time"
	"strconv"
) 


type ITransactionService interface {
	Deposit(ctx context.Context, req *pb.TransactionRequest) (*pb.TransactionResponse, error)
	WithDraw(ctx context.Context, req *pb.TransactionRequest) (*pb.TransactionResponse, error)
	TransferMoney(ctx context.Context, req *pb.TransferMoneyRequest) (*pb.TransferMoneyResponse, error)
	GetStatement(ctx context.Context, req *pb.GetStatementRequest) (*pb.GetStatementResponse, error)
	GenerateStatement(ctx context.Context, req *pb.GenerateStatementRequest) (*pb.GenerateStatementResponse, error)
	ListStatements(ctx context.Context, req *pb.ListStatementsRequest) (*pb.ListStatementsResponse, error)
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

func (s *TransactionService) Deposit(ctx context.Context, req *pb.TransactionRequest) (*pb.TransactionResponse, error) {

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


	return &pb.TransactionResponse{
		JournalId: journal.ID,
		Message: "Deposit Success",
		NewBalance: newBalance,
		


	}, nil
}

func (s *TransactionService) WithDraw(ctx context.Context, req *pb.TransactionRequest) (*pb.TransactionResponse, error){
	// Implement withdraw logic here


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
	
	accountResponse, err := s.accountRepo.FindByID(ctx, req.AccountId); 
	
	if err != nil {
		tx.Rollback()
		return nil, errors.New("account not found")
	}

	// Validate amount
	if req.Amount <= 0 || req.Amount >accountResponse.Balance {
		tx.Rollback()
		return nil, errors.New("Not Enough balance or invalid amount")
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
		EntryType: "DEBIT",
		Amount:    req.Amount,
		CreatedAt: time.Now(),
	}


	if err := s.transactionRepo.CreateJournalEntry(ctx, &creditEntry); err != nil {
		tx.Rollback()
		return nil, err
	}
		

	// Update account balance
	newBalance := accountResponse.Balance - req.Amount
	
	if err := s.accountRepo.UpdateBalance(ctx,accountResponse.ID ,newBalance); err != nil {
		tx.Rollback()
		return nil, err
	}
	// Commit transaction
	if err := tx.Commit().Error;  err != nil {
		return nil, err
	}


	return &pb.TransactionResponse{
		JournalId: journal.ID,
		Message: "Withdraw Success",
		NewBalance: newBalance,
		


	}, nil

}


func (s *TransactionService) TransferMoney(ctx context.Context, req *pb.TransferMoneyRequest) (*pb.TransferMoneyResponse, error){
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
	if req.Amount <= 0  {
		tx.Rollback()
		return nil, errors.New("amount must be greater than zero")
	}

	fromAccount, err := s.accountRepo.FindByID(ctx, req.FromAccountId); 
	toAccount, err := s.accountRepo.FindByID(ctx, req.ToAccountId);

	if err != nil {
		tx.Rollback()
		return nil, errors.New("account not found")
	}

	if fromAccount.Balance < req.Amount {
		tx.Rollback()
		return nil, errors.New("Not Enough balance in from account")
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
		AccountID: toAccount.ID,
		EntryType: "CREDIT",
		Amount:    req.Amount,
		CreatedAt: time.Now(),
	}


	debitEntry := models.JournalEntry{
		ID:        uuid.New().String(),
		JournalID: journal.ID,
		AccountID: fromAccount.ID,
		EntryType: "DEBIT",
		Amount:    req.Amount,
		CreatedAt: time.Now(),
	}


	if err := s.transactionRepo.CreateJournalEntry(ctx, &creditEntry); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := s.transactionRepo.CreateJournalEntry(ctx, &debitEntry); err != nil {
		tx.Rollback()
		return nil, err
	}
		

	// Update account balance
	fromAccountBalance := fromAccount.Balance - req.Amount
	toAccountBalance := toAccount.Balance + req.Amount
	
	if err := s.accountRepo.UpdateBalance(ctx,fromAccount.ID ,fromAccountBalance); err != nil {
		tx.Rollback()
		return nil, err
	}
	
	if err := s.accountRepo.UpdateBalance(ctx,toAccount.ID, toAccountBalance); err != nil {
		tx.Rollback()
		return nil, err
	}
	// Commit transaction
	if err := tx.Commit().Error;  err != nil {
		return nil, err
	}


	return &pb.TransferMoneyResponse{
		JournalId: journal.ID,
		Message: "Transfer Success",
		FromAccountBalance: fromAccountBalance,
		ToAccountBalance: toAccountBalance,
		


	}, nil
}
// GenerateStatement generates a statement for an account for a specific period
func (s *TransactionService) GenerateStatement(ctx context.Context, req *pb.GenerateStatementRequest) (*pb.GenerateStatementResponse, error) {
	// Parse dates
	periodStart, err := time.Parse(time.RFC3339, req.StartDate)
	if err != nil {
		return nil, errors.New("invalid period_start format:" + err.Error())
	}

	periodEnd, err := time.Parse(time.RFC3339, req.EndDate)
	if err != nil {
		return nil, errors.New("invalid period_end format:"+ err.Error())
	}

	// Fetch account

	account, err := s.accountRepo.FindByID(ctx, req.AccountId)

	// Calculate opening balance (balance at period start)
	journalEntriesBeforePeriod, err := s.transactionRepo.GetJournalEntriesBeforePeriod(ctx, req.AccountId, req.StartDate)
	if err != nil {
		return nil, errors.New("failed to fetch opening balance: "+ err.Error())
	}

	openingBalance := 0.0
	for _, entry := range journalEntriesBeforePeriod {
		if entry.EntryType == "CREDIT" {
			openingBalance += entry.Amount
		} else if entry.EntryType == "DEBIT" {
			openingBalance -= entry.Amount
		}
	}

	// Fetch transactions for the period	
	journalEntriesBetweenPeriod, err := s.transactionRepo.GetJournalEntriesBetweenPeriod(ctx, req.AccountId, req.StartDate, req.EndDate)
	if err != nil {
		return nil, errors.New("failed to fetch journal entries: "+ err.Error())
	}

	// Calculate closing balance
	closingBalance := openingBalance
	var transactions []*pb.StatementTransaction
	for _, entry := range journalEntriesBetweenPeriod {
		if entry.EntryType == "CREDIT" {
			closingBalance += entry.Amount
		} else if entry.EntryType == "DEBIT" {
			closingBalance -= entry.Amount
		}

		transactions = append(transactions, &pb.StatementTransaction{
			TransactionId: entry.ID,
			Date:          entry.Journal.ValueDate.Format(time.RFC3339),
			Description:   entry.Journal.Narration,
			EntryType:     entry.EntryType,
			Amount:        entry.Amount,
			Balance:       closingBalance,
		})
	}

	// Create statement record
	statement := models.Statement{
		ID:             uuid.New().String(),
		AccountID:      account.ID, 
		FromDate:    periodStart,
		ToDate:      periodEnd,
		OpeningBalance: openingBalance,
		ClosingBalance: closingBalance,
		GeneratedAt:    time.Now(),
	}
	

	if err := s.transactionRepo.CreateStatement(ctx,&statement); err != nil {
		return nil, errors.New("failed to create statement:" + err.Error())
	}

	return &pb.GenerateStatementResponse{
		StatementId:    statement.ID,
		AccountId:      req.AccountId,
		StartDate:    statement.FromDate.String(),
		EndDate:      statement.ToDate.String(),
		OpeningBalance: strconv.FormatFloat(statement.OpeningBalance, 'f', 2, 64),
		ClosingBalance: strconv.FormatFloat(statement.ClosingBalance, 'f', 2, 64),
		Transactions:   transactions,
		GeneratedAt:    statement.GeneratedAt.Format(time.RFC3339),
	}, nil
}

// GetStatement retrieves a previously generated statement
func (s *TransactionService) GetStatement(ctx context.Context, req *pb.GetStatementRequest) (*pb.GetStatementResponse, error) {
	

	return &pb.GetStatementResponse{}, nil
}
func (s *TransactionService) ListStatements(ctx context.Context, req *pb.ListStatementsRequest) (*pb.ListStatementsResponse, error){
	// Implement list statements logic here
	return &pb.ListStatementsResponse{}, nil
}