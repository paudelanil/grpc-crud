package models

import (
	"time"
)
// Journal represents the header for a double-entry transaction
type Journal struct {
	ID          string    `gorm:"primaryKey;column:journal_id"`
	ReferenceID string    `gorm:"uniqueIndex;not null"` // idempotency key
	Narration   string    `gorm:"not null"`
	ValueDate   time.Time `gorm:"not null"`
	Status      string    `gorm:"type:varchar(20);not null;default:'posted'"` // posted, reversed

	JournalEntries []JournalEntry `gorm:"foreignKey:JournalID;constraint:OnDelete:CASCADE"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Journal) TableName() string {
	return "journals"
}

// JournalEntry represents individual debit/credit entries in double-entry bookkeeping
type JournalEntry struct {
	ID        string  `gorm:"primaryKey;column:journal_entry_id"`
	JournalID string  `gorm:"not null;index"`
	Journal   Journal `gorm:"foreignKey:JournalID"`

	AccountID string  `gorm:"not null;index"`
	Account   Account `gorm:"foreignKey:AccountID"`

	EntryType string  `gorm:"type:varchar(6);not null"`    // DEBIT or CREDIT
	Amount    float64 `gorm:"type:numeric(18,2);not null"` // amount in decimal

	CreatedAt time.Time
}

func (JournalEntry) TableName() string {
	return "journal_entries"
}
