package transactions

import (
	"github.com/stevenwijaya/finance-tracker/database"
	models "github.com/stevenwijaya/finance-tracker/models/transactions"
	"github.com/stevenwijaya/finance-tracker/pkg/utils"
)

func CreateTransaction(transactions *models.Transaction) error {
	result := database.DB.Create(transactions)
	return result.Error
}

func GetAllTransaction(userID uint, filters map[string]interface{}, pagination utils.Pagination) ([]models.Transaction, error) {
	var transactions []models.Transaction
	query := database.DB.Preload("User").Where("user_id = ?", userID).Find(&transactions)

	if val, ok := filters["type"]; ok && val != "" {
		query = query.Where("type = ?", val)
	}

	if val, ok := filters["category"]; ok && val != "" {
		query = query.Where("category = ?", val)
	}

	if startDate, ok := filters["start_date"]; ok && startDate != "" {
		if endDate, ok := filters["end_date"]; ok && endDate != "" {
			query = query.Where("date BETWEEN ? AND ?", startDate, endDate)
		}
	}

	result := query.
		Order("date DESC").
		Offset(pagination.Offset()).
		Limit(pagination.Limit).
		Find(&transactions)

	return transactions, result.Error
}

func GetTransactionById(transactionId uint, userId uint) (models.Transaction, error) {
	var transactions models.Transaction
	result := database.DB.Preload("User").Where("id = ? AND user_id = ?", transactionId, userId).First(&transactions)
	return transactions, result.Error
}

func GetTransactionSummary(userId uint, startDate, endDate string) (map[string]float64, error) {
	type Summary struct {
		Income  float64
		Expense float64
	}

	var summary Summary
	query := database.DB.Table("transactions").Select(`
    	COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) AS income,
    	COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) AS expense
	`).Where("user_id = ?", userId)

	if startDate != "" && endDate != "" {
		// filter antara dua tanggal
		query = query.Where("date BETWEEN ? AND ?", startDate, endDate)
	} else if startDate != "" {
		// filter dari startDate ke atas
		query = query.Where("date >= ?", startDate)
	} else if endDate != "" {
		// filter sampai endDate ke bawah
		query = query.Where("date <= ?", endDate)
	}

	if err := query.Scan(&summary).Error; err != nil {
		return nil, err
	}

	result := map[string]float64{
		"income":  summary.Income,
		"expense": summary.Expense,
		"balance": summary.Income - summary.Expense,
	}

	return result, nil
}

func GetTransactionSummaryByCategory(userId uint, startDate, endDate, tType string) ([]map[string]interface{}, error) {
	var summaries []map[string]interface{}

	db := database.DB.Table("transactions").
		Select("category, COALESCE(SUM(amount), 0) as total_amount").
		Where("user_id", userId)

	if tType != "" {
		db = db.Where("type = ?", tType)
	}

	if startDate != "" && endDate != "" {
		db = db.Where("date BETWEEN ? AND ?", startDate, endDate)
	} else if startDate != "" {
		db = db.Where("date >= ?", startDate)
	} else if endDate != "" {
		db = db.Where("date <= ?", endDate)
	}

	db = db.Group("category").Order("category ASC")

	if err := db.Find(&summaries).Error; err != nil {
		return nil, err
	}

	return summaries, nil
}

func UpdateTransaction(transaction *models.Transaction) error {
	result := database.DB.Save(transaction)
	return result.Error
}

func DeleteTransaction(transactionId uint, userId uint) error {
	result := database.DB.Where("id = ? AND user_id = ?", transactionId, userId).Delete(&models.Transaction{})
	return result.Error
}
