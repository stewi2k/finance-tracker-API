package transactions

import (
	"errors"
	"time"

	models "github.com/stevenwijaya/finance-tracker/models/transactions"
	"github.com/stevenwijaya/finance-tracker/pkg/log"
	repositories "github.com/stevenwijaya/finance-tracker/repositories/transactions"
	repositoriesUser "github.com/stevenwijaya/finance-tracker/repositories/users"
)

func CreateTransaction(input *models.Transaction) error {

	user, err := repositoriesUser.GetUserByID(*input.UserID)
	if err != nil {
		log.Error("Error fetching user by username:", err)
		return err
	}

	if input.Date.Time.IsZero() {
		input.Date.Time = time.Now()
	}

	input.User = user

	if err := repositories.CreateTransaction(input); err != nil {
		log.Error("Error Creating Transaction")
		return err
	}

	return nil
}

func GetAllTransaction(userId uint) ([]models.Transaction, error) {
	transactions, err := repositories.GetAllTransaction(userId)
	if err != nil {
		log.Errorf("Transaction by id %v not found", userId)
		return nil, err
	}

	return transactions, nil
}

func GetTransactionById(transactionsId uint, userId uint) (models.Transaction, error) {
	transactions, err := repositories.GetTransactionById(transactionsId, userId)
	if err != nil {
		log.Errorf("Transaction with id %v by id %v not found", transactionsId, userId)
		return models.Transaction{}, err
	}

	return transactions, nil
}

func UpdateTransaction(transactionId uint, userId uint, input *models.Transaction) error {
	transaction, err := repositories.GetTransactionById(transactionId, userId)
	if err != nil {
		return err
	}

	if transaction.ID == 0 {
		log.Error("Transaction not found")
		return errors.New("Transaction not Found")
	}

	if input.Date.Time.IsZero() {
		input.Date.Time = time.Now()
	}

	// ✳️ update hanya field yang diubah
	transaction.Type = input.Type
	transaction.Amount = input.Amount
	transaction.Category = input.Category
	transaction.Description = input.Description
	transaction.Date = input.Date

	if err := repositories.UpdateTransaction(&transaction); err != nil {
		log.Error("Failed to update transaction")
		return err
	}

	return nil
}

func DeleteTransaction(transactionId uint, userId uint) error {
	if err := repositories.DeleteTransaction(transactionId, userId); err != nil {
		log.Error("Failed to delete transaction")
		return err
	}

	return nil
}
