package repository

import (
	"organization-manager/pkg/database/mongodb/models"
)

var Users []models.User
var Organizations []models.Organization

func GetUserByEmail(email string) *models.User {
	for i, usr := range Users {
		if usr.Email == email {
			return &Users[i]
		}
	}
	return nil
}

func GetUserByRefreshToken(refreshToken string) *models.User {
	for _, usr := range Users {
		if usr.RefreshToken == refreshToken {
			return &usr
		}
	}
	return nil
}

func GetUserByAccessToken(accessToken string) *models.User {
	for i, usr := range Users {
		if usr.AccessToken == accessToken {
			return &Users[i]
		}
	}
	return nil
}

func GetOrganizationById(ID string) *models.Organization {
	for _, org := range Organizations {
		if org.ID == ID {
			return &org
		}
	}
	return nil
}

func GetUserOrgs(usr models.User) []models.Organization {
	var orgs []models.Organization
	for _, org := range Organizations {
		if IsMember(usr, org) {
			orgs = append(orgs, org)
		}
	}
	return orgs
}

func AddMemberToOrganization(org *models.Organization, usr models.User, accessLevel string) {
	var member models.OrganizationMember
	member.Email = usr.Email
	member.Name = usr.Name
	member.AccessLevel = accessLevel
	org.Members = append(org.Members, member)
}

func IsMember(usr models.User, org models.Organization) bool {
	for _, mem := range org.Members {
		if mem.Email == usr.Email {
			return true
		}
	}
	return false
}
