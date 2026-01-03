package repository

import (
	"context"
	"errors"

	"github.com/paudelanil/grpc-crud/models"
	"gorm.io/gorm"
)

// ICustomerRepository defines the interface for customer data operations
type ICustomerRepository interface {
	Create(ctx context.Context, customer *models.Customer) error
	FindByID(ctx context.Context, id string) (*models.Customer, error)
	FindAll(ctx context.Context, limit, offset int) ([]*models.Customer, error)
	Update(ctx context.Context, customer *models.Customer) error
	Delete(ctx context.Context, id string) error
	IsEmailTaken(ctx context.Context, email string) (bool, error)
	IsPhoneTaken(ctx context.Context, phone string) (bool, error)
}

// CustomerRepository implements ICustomerRepository interface
type CustomerRepository struct {
	db *gorm.DB
}

// NewCustomerRepository creates a new instance of CustomerRepository
func NewCustomerRepository(db *gorm.DB) ICustomerRepository {
	return &CustomerRepository{db: db}
}

// Create creates a new customer in the database
func (r *CustomerRepository) Create(ctx context.Context, customer *models.Customer) error {
	result := r.db.WithContext(ctx).Create(customer)
	return result.Error
}

// FindByID finds a customer by ID
func (r *CustomerRepository) FindByID(ctx context.Context, id string) (*models.Customer, error) {
	var customer models.Customer
	result := r.db.WithContext(ctx).Where("customer_id = ?", id).First(&customer)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("customer not found")
		}
		return nil, result.Error
	}
	return &customer, nil
}

// FindAll retrieves all customers with pagination
func (r *CustomerRepository) FindAll(ctx context.Context, limit, offset int) ([]*models.Customer, error) {
	var customers []*models.Customer
	result := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&customers)
	if result.Error != nil {
		return nil, result.Error
	}
	return customers, nil
}

// Update updates a customer in the database
func (r *CustomerRepository) Update(ctx context.Context, customer *models.Customer) error {
	result := r.db.WithContext(ctx).Save(customer)
	return result.Error
}

// Delete soft deletes a customer by ID
func (r *CustomerRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("customer_id = ?", id).Delete(&models.Customer{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("customer not found")
	}
	return nil
}

// IsEmailTaken checks if an email is already taken
func (r *CustomerRepository) IsEmailTaken(ctx context.Context, email string) (bool, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&models.Customer{}).Where("email = ?", email).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

// IsPhoneTaken checks if a phone number is already taken
func (r *CustomerRepository) IsPhoneTaken(ctx context.Context, phone string) (bool, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&models.Customer{}).Where("phone = ?", phone).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}
