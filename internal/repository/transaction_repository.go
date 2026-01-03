package repository

import (
	"context"
	
	"github.com/paudelanil/grpc-crud/models"
	"gorm.io/gorm"
)

// ITransactionRepository defines the interface for transaction data operations
type ITransactionRepository interface {
	CreateJournalEntry(ctx context.Context, entry *models.JournalEntry) error

	CreateJournal(ctx context.Context, journal *models.Journal) error

	GetJournalByRefID(ctx context.Context, id string) (*models.Journal, error)
	IsJournalExists(ctx context.Context, refID string) (bool, error)
	BeginTransaction(ctx context.Context) (*gorm.DB, error)
}

// TransactionRepository implements ITransactionRepository interface
type TransactionRepository struct {
	db *gorm.DB
}

// NewTransactionRepository creates a new instance of TransactionRepository
func NewTransactionRepository(db *gorm.DB) ITransactionRepository {
	return &TransactionRepository{db: db}
}


// CreateJournalEntry creates a new journal entry in the database
func (r *TransactionRepository) CreateJournalEntry(ctx context.Context, entry *models.JournalEntry) error {
	result := r.db.WithContext(ctx).Create(entry)
	return result.Error
}

// CreateJournal creates a new journal in the database
func (r *TransactionRepository) CreateJournal(ctx context.Context, journal *models.Journal) error {
	result := r.db.WithContext(ctx).Create(journal)
	return result.Error
}

// GetJournalByRefID retrieves a journal by its reference ID
func (r *TransactionRepository) GetJournalByRefID(ctx context.Context, id string) (*models.Journal, error) {
	var journal models.Journal
	result := r.db.WithContext(ctx).Where("reference_id = ?", id).First(&journal)
	if result.Error != nil {
		return nil, result.Error
	}
	return &journal, nil
}

func (r *TransactionRepository) IsJournalExists(ctx context.Context, refID string) (bool, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&models.Journal{}).Where("reference_id = ?", refID).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

func (r *TransactionRepository) BeginTransaction(ctx context.Context) (*gorm.DB, error) {
	tx := r.db.WithContext(ctx).Begin()
	return tx, tx.Error
}


