package repositories

import (
	"errors"

	"github.com/stevenwijaya/finance-tracker/database"
	"github.com/stevenwijaya/finance-tracker/models"
	"gorm.io/gorm"
)

func CreateUser(user *models.User) error {
	result := database.DB.Create(&user)
	return result.Error
}

func GetUserByUsername(username string) (*models.User, error) {
	var user models.User

	result := database.DB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}
	return &user, nil
}

func GetUserByID(id uint) (*models.User, error) {
	var user models.User
	result := database.DB.First(&user, &id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}
	return &user, nil
}
