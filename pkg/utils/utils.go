package utils

import (
	"organization-manager/pkg/database/mongodb/models"

	"github.com/google/uuid"
)

func UpdateSession(usr *models.User) (string, string) {
	accessToken := "Bearer " + uuid.New().String()
	refreshToken := uuid.New().String()
	return accessToken, refreshToken
}

func CreateUser(signupData models.SignupModel) models.User {
	var usr models.User
	usr.Name = signupData.Name
	usr.Email = signupData.Email
	usr.HashedPassword = signupData.Password
	usr.AccessToken, usr.RefreshToken = UpdateSession(&usr)
	return usr
}
