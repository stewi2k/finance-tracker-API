package transactions

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	models "github.com/stevenwijaya/finance-tracker/models/transactions"
	"github.com/stevenwijaya/finance-tracker/pkg/log"
	"github.com/stevenwijaya/finance-tracker/pkg/response"
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
		response.Error(c, http.StatusBadRequest, "Invalid transaction input : "+err.Error())
		return
	}

	if err := validator.Validate.Struct(&input); err != nil {
		log.Warn("Transaction validation failed : ", err)
		response.Error(c, http.StatusBadRequest, "Transaction validation failed : "+err.Error())
		return
	}

	input.UserID = &userId

	if err := services.CreateTransaction(&input); err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to created transaction : "+err.Error())
		return
	}

	log.Info("Transaction created for user : ", userId)
	response.Success(c, http.StatusOK, "Transaction created successfully", input)
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
		response.Error(c, http.StatusNotFound, "Failed to get transaction : "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Get transaction succesfully", gin.H{
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
		response.Error(c, http.StatusBadRequest, "Invalid input transaction id : "+err.Error())
		return
	}

	userId := c.GetUint("user_id")

	transaction, err := services.GetTransactionById(uint(id), userId)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Failed to get transaction by id "+idParam+": "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "get transaction by id "+idParam, transaction)
}

func GetTransactionSummary(c *gin.Context) {
	userId := c.GetUint("user_id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate != "" && endDate != "" && startDate > endDate {
		log.Error("Invalid data range")
		response.Error(c, http.StatusBadRequest, "Invalid data range")
		return
	}

	summary, err := services.GetTransactionSummary(userId, startDate, endDate)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get data transaction summary : "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Successfully retrieved transaction summary", gin.H{
		"summary": summary,
	})
}

func GetTransactionSummaryByCategory(c *gin.Context) {
	userId := c.GetUint("user_id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	tType := c.Query("type")

	if startDate != "" && endDate != "" && startDate > endDate {
		log.Error("Invalid data range")
		response.Error(c, http.StatusBadRequest, "Invalid data range")
		return
	}

	if tType == "" {
		log.Error("Type can't be empty")
		response.Error(c, http.StatusBadRequest, "Type can't be empty")
		return
	}

	summary, err := services.GetTransactionSummaryByCategory(userId, startDate, endDate, tType)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get data summary by category : "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Successfully retrieve transaction summary by category", gin.H{
		"summary": summary,
	})
}

// update transaction
func UpdateTransaction(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input transaction id : "+err.Error())
		return
	}

	userId := c.GetUint("user_id")

	var input models.Transaction
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		log.Error("Invalid transaction input : ", err)
		response.Error(c, http.StatusBadRequest, "Invalid transaction input : "+err.Error())
		return
	}

	if err := validator.Validate.Struct(&input); err != nil {
		log.Error("Transaction validation failed : ", err)
		response.Error(c, http.StatusBadRequest, "Transaction validation failed : "+err.Error())
		return
	}

	if err := services.UpdateTransaction(uint(id), userId, &input); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to udpate data transaction : "+err.Error())
		return
	}

	log.Info("Updated successfully for user : ", userId)
	response.Success(c, http.StatusOK, "Updated successfully", nil)
}

func DeleteTransaction(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input transaction id : "+err.Error())
		return
	}

	userId := c.GetUint("user_id")

	if err := services.DeleteTransaction(uint(id), userId); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete data transaction : "+err.Error())
		return
	}

	log.Info("Deleted successfully for user : ", userId)
	response.Success(c, http.StatusOK, "Deleted successfully", nil)
}
