package models

import (
	"time"
)

type Transaction struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"not null" json:"user_id"`
	Type        string    `gorm:"not null" json:"type"` // "income" or "expense"
	Amount      float64   `gorm:"not null" json:"amount"`
	Category    string    `gorm:"not null" json:"category"`
	Description string    `json:"description"`
	Date        time.Time `gorm:"not null" json:"date"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}
