package handlers

import (
	"organization-manager/pkg/database/mongodb/models"
)

func buildSignupResponse() models.SignupResponse {
	var response models.SignupResponse
	response.Message = "Sign Up Succeeded."
	return response
}

func buildSigninResponse(usr *models.User) models.SigninResponse {
	var response models.SigninResponse
	response.Message = "Sign In Succeeded."
	response.AccessToken = usr.AccessToken
	response.RefreshToken = usr.RefreshToken
	return response
}

func buildRefreshResponse(usr *models.User) models.RefreshResponse {
	var response models.RefreshResponse
	response.Message = "Refreshing Tokens Succeeded."
	response.AccessToken = usr.AccessToken
	response.RefreshToken = usr.RefreshToken
	return response
}

func buildCreateOrganizationResponse(org models.Organization) models.CreateOrganizationResponse {
	var response models.CreateOrganizationResponse
	response.OrganizationID = org.ID
	return response
}

func buildReadOrganizationResponse(org models.Organization) models.Organization {
	return org
}

func buildReadAllOrganizationResponse(org []models.Organization) []models.Organization {
	return org
}

func buildUpdateOrganizationResponse(org models.Organization) models.Organization {
	org.Members = nil
	return org
}

func buildInviteUserResponse() models.InviteResponse {
	var response models.InviteResponse
	response.Message = "User Invited"
	return response
}

func buildDeleteOrganizationResponse() models.DeleteResponse {
	var response models.DeleteResponse
	response.Message = "Organization Deleted"
	return response
}
