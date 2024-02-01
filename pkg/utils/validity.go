package utils

import (
	"organization-manager/pkg/database/mongodb/models"
	"organization-manager/pkg/database/mongodb/repository"
)

func IsValidSignup(signupData models.SignupModel) bool {
	user := repository.GetUserByEmail(signupData.Email)
	return user == nil
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

	user := repository.GetUserByRefreshToken(refreshToken)
	return user != nil
}

func IsAuthorized(authHeader string) bool {
	user := repository.GetUserByAccessToken(authHeader)
	return user != nil
}

func IsValidOrganization(name string) bool {
	if name == "" {
		return false
	}

	organization := repository.GetOrganizationByName(name)
	return organization == nil
}

func IsValidOrganizationID(ID string) bool {
	if ID == "" {
		return false
	}

	organization := repository.GetOrganizationById(ID)

	return organization == nil
}
