package services

import (
	"errors"

	"github.com/stevenwijaya/finance-tracker/models"
	"github.com/stevenwijaya/finance-tracker/pkg/log"
	"github.com/stevenwijaya/finance-tracker/pkg/utils"
	"github.com/stevenwijaya/finance-tracker/repositories"
)

func RegisterUser(user *models.User) error {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		log.Error("Error hashing password")
		return err
	}
	user.Password = hashedPassword
	if err := repositories.CreateUser(user); err != nil {
		log.Error("Error creating user")
		return err
	}

	return nil
}

func LoginUser(username, password string) (*models.User, error) {
	user, err := repositories.GetUserByUsername(username)
	if err != nil {
		log.Error("Error fetching user by username:", err)
		return nil, err
	}

	if user == nil {
		log.Error("User not found")
		return nil, errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(user.Password, password) {
		log.Error("Invalid password")
		return nil, errors.New("Invalid Credentials")
	}

	return user, nil
}
