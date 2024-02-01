package utils

import (
	"fmt"
	"organization-manager/pkg/database/mongodb/models"
	"organization-manager/pkg/database/mongodb/repository"
)

func IsValidSignup(signupData models.SignupModel) bool {
	for _, usr := range repository.Users {
		if usr.Email == signupData.Email {
			return false
		}
	}
	return true
}

func IsValidSignin(signinData models.SigninModel) bool {
	usr := repository.GetUserByEmail(signinData.Email)
	if usr == nil || usr.HashedPassword != signinData.Password {
		return false
	}
	return true
}

func IsValidRefreshToken(refreshToken string) bool {
	if refreshToken == "" {
		return false
	}

	for _, usr := range repository.Users {
		if usr.RefreshToken == refreshToken {
			return true
		}
	}
	return false
}

func IsAuthorized(authHeader string) bool {
	for _, usr := range repository.Users {
		fmt.Println(usr.AccessToken)
		if usr.AccessToken == authHeader {
			return true
		}
	}
	return false
}

func IsValidOrganization(name string) bool {
	if name == "" {
		return false
	}

	for _, org := range repository.Organizations {
		if org.Name == name {
			return false
		}
	}
	return true
}

func IsValidOrganizationID(ID string) bool {
	if ID == "" {
		return false
	}

	for _, org := range repository.Organizations {
		if org.ID == ID {
			return true
		}
	}

	return false
}
