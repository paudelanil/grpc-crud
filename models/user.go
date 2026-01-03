package models

import (
	"time"

	"gorm.io/gorm"
)

type Customer struct {
	ID        string `gorm:"primaryKey;column:customer_id"`
	FirstName string `gorm:"not null"`
	LastName  string `gorm:"not null"`
	Address   string
	Email     string `gorm:"not null;uniqueIndex"`
	Phone     string `gorm:"not null;uniqueIndex"`

	Accounts []Account `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (Customer) TableName() string {
	return "customers"
}

type Account struct {
	ID            string    `gorm:"primaryKey;column:account_id"`
	AccountNumber string    `gorm:"uniqueIndex;not null"`
	Status        string    `gorm:"type:varchar(20);not null"` // active, frozen, closed
	Balance       float64   `gorm:"type:numeric(18,2);not null;default:0"`
	OpenedAt      time.Time `gorm:"not null"`

	CustomerID  string   `gorm:"not null"`
	Customer    Customer `gorm:"foreignKey:CustomerID"`
	Currency    string   `gorm:"type:varchar(3);not null;default:'NPR'"`
	AccountType string   `gorm:"type:varchar(20);not null;default:'savings'"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (Account) TableName() string {
	return "accounts"
}

// User represents the authentication user
type User struct {
	ID        string `gorm:"primaryKey;column:user_id"`
	Username  string `gorm:"uniqueIndex;not null"`
	Password  string `gorm:"not null"` // hashed password
	Email     string `gorm:"uniqueIndex;not null"`
	IsActive  bool   `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (User) TableName() string {
	return "users"
}
