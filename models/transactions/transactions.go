package transactions

import (
	"time"

	"github.com/stevenwijaya/finance-tracker/models/users"
	"github.com/stevenwijaya/finance-tracker/pkg/utils"
)

type Transaction struct {
	ID          uint             `gorm:"primaryKey" json:"id"`
	UserID      *uint            `gorm:"not null;index" json:"-"` // Foreign key ke tabel users
	User        *users.User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:UserID;references:ID" json:"user,omitempty"`
	Type        string           `gorm:"not null" json:"type" validate:"required"` // "income" or "expense"
	Amount      float64          `gorm:"not null" json:"amount" validate:"required"`
	Category    string           `gorm:"not null" json:"category" validate:"required"`
	Description string           `json:"description"`
	Date        utils.CustomDate `gorm:"not null" json:"date" validate:"required"`
	CreatedAt   time.Time        `gorm:"autoCreateTime" json:"-"`
	UpdatedAt   time.Time        `gorm:"autoUpdateTime" json:"-"`
}
