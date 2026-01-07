package repository

import (
	"context"
	"errors"

	"github.com/paudelanil/grpc-crud/models"
	"gorm.io/gorm"
)

// IAccountRepository defines the interface for account data operations
type IAccountRepository interface {
	Create(ctx context.Context, account *models.Account) error
	FindByID(ctx context.Context, id string) (*models.Account, error)
	FindByCustomerID(ctx context.Context, customerID string) ([]*models.Account, error)
	FindAll(ctx context.Context, limit, offset int) ([]*models.Account, error)
	Update(ctx context.Context, account *models.Account) error
	Delete(ctx context.Context, id string) error
	IsAccountNumberTaken(ctx context.Context, accountNumber string) (bool, error)
	UpdateBalance(ctx context.Context, accountID string, newBalance float64) error
	GetBalanceByID(ctx context.Context, id string) (float64, error)
}

// AccountRepository implements IAccountRepository interface
type AccountRepository struct {
	db *gorm.DB
}

// NewAccountRepository creates a new instance of AccountRepository
func NewAccountRepository(db *gorm.DB) IAccountRepository {
	return &AccountRepository{db: db}
}

// Create creates a new account in the database
func (r *AccountRepository) Create(ctx context.Context, account *models.Account) error {
	result := r.db.WithContext(ctx).Create(account)
	return result.Error
}

// FindByID finds an account by ID
func (r *AccountRepository) FindByID(ctx context.Context, id string) (*models.Account, error) {
	var account models.Account
	result := r.db.WithContext(ctx).Where("account_id = ?", id).First(&account)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("account not found")
		}
		return nil, result.Error
	}
	return &account, nil
}

// FindByCustomerID finds all accounts for a specific customer
func (r *AccountRepository) FindByCustomerID(ctx context.Context, customerID string) ([]*models.Account, error) {
	var accounts []*models.Account
	result := r.db.WithContext(ctx).Where("customer_id = ?", customerID).Find(&accounts)
	if result.Error != nil {
		return nil, result.Error
	}
	return accounts, nil
}

// FindAll retrieves all accounts with pagination
func (r *AccountRepository) FindAll(ctx context.Context, limit, offset int) ([]*models.Account, error) {
	var accounts []*models.Account
	result := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&accounts)
	if result.Error != nil {
		return nil, result.Error
	}
	return accounts, nil
}

// Update updates an account in the database
func (r *AccountRepository) Update(ctx context.Context, account *models.Account) error {
	result := r.db.WithContext(ctx).Save(account)
	return result.Error
}

// Delete soft deletes an account by ID
func (r *AccountRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("account_id = ?", id).Delete(&models.Account{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("account not found")
	}
	return nil
}

// IsAccountNumberTaken checks if an account number is already taken
func (r *AccountRepository) IsAccountNumberTaken(ctx context.Context, accountNumber string) (bool, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&models.Account{}).Where("account_number = ?", accountNumber).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}


//UpdateBalance updates the balance and update time for given accountID

func (r *AccountRepository) UpdateBalance(ctx context.Context, accountID string, newBalance float64) error {
	result := r.db.WithContext(ctx).Model(&models.Account{}).Where("account_id = ?", accountID).Updates(map[string]interface{}{
		"balance":     newBalance,
		"updated_at":  gorm.Expr("NOW()"),
	})
	return result.Error
}

// GetBalanceByID retrieves the balance of an account by ID
func (r *AccountRepository) GetBalanceByID(ctx context.Context, id string) (float64, error) {
	var account models.Account
	result := r.db.WithContext(ctx).Select("balance").Where("account_id = ?", id).First(&account)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, errors.New("account not found")
		}
		return 0, result.Error
	}
	return account.Balance, nil
}