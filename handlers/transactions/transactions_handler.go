package transactions

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	models "github.com/stevenwijaya/finance-tracker/models/transactions"
	"github.com/stevenwijaya/finance-tracker/pkg/log"
	"github.com/stevenwijaya/finance-tracker/pkg/utils"
	"github.com/stevenwijaya/finance-tracker/pkg/validator"
	services "github.com/stevenwijaya/finance-tracker/services/transactions"
)

// create transaction
func CreateTransaction(c *gin.Context) {
	userId := c.GetUint("user_id")

	var input models.Transaction
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		log.Error("Invalid transaction input : ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(&input); err != nil {
		log.Warn("Transaction validation failed : ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.UserID = &userId

	if err := services.CreateTransaction(&input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	log.Info("Transaction created for user : ", userId)
	c.JSON(http.StatusCreated, gin.H{"message": "Transaction Created", "transaction": input})
}

// get all transaction
func GetAllTransaction(c *gin.Context) {
	userId := c.GetUint("user_id")

	pagination := utils.GetPagination(c)
	filters := map[string]interface{}{
		"type":       c.Query("type"),
		"category":   c.Query("category"),
		"start_date": c.Query("start_date"),
		"end_date":   c.Query("end_date"),
	}

	transactions, err := services.GetAllTransaction(userId, filters, pagination)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"page":         pagination.Page,
		"limit":        pagination.Limit,
		"total":        len(transactions),
		"transactions": transactions,
	})
}

// get transaction by id
func GetTransactionById(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	userId := c.GetUint("user_id")

	transaction, err := services.GetTransactionById(uint(id), userId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transaction": transaction})
}

func GetTransactionSummary(c *gin.Context) {
	userId := c.GetUint("user_id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate != "" && endDate != "" && startDate > endDate {
		log.Error("Invalid Data Range")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data range"})
		return
	}

	summary, err := services.GetTransactionSummary(userId, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Transaction summary retrieved successfully",
		"summary": summary,
	})
}

func GetTransactionSummaryByCategory(c *gin.Context) {
	userId := c.GetUint("user_id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	tType := c.Query("type")

	if startDate != "" && endDate != "" && startDate > endDate {
		log.Error("Invalid Data Range")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data range"})
		return
	}

	if tType == "" {
		log.Error("Type can't be empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Type can't be empty"})
	}

	summary, err := services.GetTransactionSummaryByCategory(userId, startDate, endDate, tType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	log.Info("summary : ", summary)

	c.JSON(http.StatusOK, gin.H{
		"message": "Transaction summary retrieved successfully",
		"summary": summary,
	})
}

// update transaction
func UpdateTransaction(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
	}

	userId := c.GetUint("user_id")

	var input models.Transaction
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		log.Error("Invalid transaction input : ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(&input); err != nil {
		log.Error("Transaction validation failed : ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.UpdateTransaction(uint(id), userId, &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Info("Updated Successfully for user : ", userId)
	c.JSON(http.StatusOK, gin.H{"message": "Updated Successfully"})
}

func DeleteTransaction(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
	}

	userId := c.GetUint("user_id")

	if err := services.DeleteTransaction(uint(id), userId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Info("Deleted Successfully for user : ", userId)
	c.JSON(http.StatusOK, gin.H{"message": "Deleted Successfully"})
}
