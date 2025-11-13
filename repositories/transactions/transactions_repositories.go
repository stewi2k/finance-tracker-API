package transactions

import (
	"github.com/stevenwijaya/finance-tracker/database"
	models "github.com/stevenwijaya/finance-tracker/models/transactions"
)

func CreateTransaction(transactions *models.Transaction) error {
	result := database.DB.Create(transactions)
	return result.Error
}

func GetAllTransaction(userID uint) ([]models.Transaction, error) {
	var transactions []models.Transaction
	result := database.DB.Preload("User").Where("user_id = ?", userID).Find(&transactions)
	return transactions, result.Error
}

func GetTransactionById(transactionId uint, userId uint) (models.Transaction, error) {
	var transactions models.Transaction
	result := database.DB.Preload("User").Where("id = ? AND user_id = ?", transactionId, userId).First(&transactions)
	return transactions, result.Error
}

func UpdateTransaction(transaction *models.Transaction) error {
	result := database.DB.Save(transaction)
	return result.Error
}

func DeleteTransaction(transactionId uint, userId uint) error {
	result := database.DB.Where("id = ? AND user_id = ?", transactionId, userId).Delete(&models.Transaction{})
	return result.Error
}
