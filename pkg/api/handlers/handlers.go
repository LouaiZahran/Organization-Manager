package handlers

import (
	"net/http"
	"organization-manager/pkg/database/mongodb/models"
	"organization-manager/pkg/database/mongodb/repository"
	"organization-manager/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SignupHandler(c *gin.Context) {
	var signupData models.SignupModel
	err := c.BindJSON(&signupData)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if !utils.IsValidSignup(signupData) {
		c.IndentedJSON(http.StatusConflict, gin.H{"message": "The Username or password are not valid."})
		return
	}

	usr := utils.CreateUser(signupData)
	repository.AddUser(usr)
	c.IndentedJSON(http.StatusOK, buildSignupResponse())
}

func SigninHandler(c *gin.Context) {
	var signinData models.SigninModel
	err := c.BindJSON(&signinData)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if !utils.IsValidSignin(signinData) {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "The Username or password are not correct."})
		return
	}

	usr := repository.GetUserByEmail(signinData.Email)
	accessToken, refreshToken := utils.UpdateSession(usr)
	usr.AccessToken = accessToken
	usr.RefreshToken = refreshToken
	repository.UpdateUser(usr)
	c.IndentedJSON(http.StatusOK, buildSigninResponse(usr))
}

func RefreshHandler(c *gin.Context) {
	var refreshData models.RefreshModel
	err := c.BindJSON(&refreshData)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if !utils.IsValidRefreshToken(refreshData.RefreshToken) {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "The refresh token is not valid."})
		return
	}

	usr := repository.GetUserByRefreshToken(refreshData.RefreshToken)
	accessToken, refreshToken := utils.UpdateSession(usr)
	usr.AccessToken = accessToken
	usr.RefreshToken = refreshToken
	repository.UpdateUser(usr)
	c.IndentedJSON(http.StatusOK, buildRefreshResponse(usr))
}

func CreateOrganizationHandler(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if !utils.IsAuthorized(authHeader) {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var createOrganizationData models.CreateOrganizationModel
	err := c.BindJSON(&createOrganizationData)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if !utils.IsValidOrganization(createOrganizationData.Name) {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "The organization name is not valid"})
		return
	}

	var org models.Organization
	org.Name = createOrganizationData.Name
	org.Description = createOrganizationData.Description
	org.ID = uuid.New().String()
	repository.AddOrganization(org)
	repository.AddMemberToOrganization(&org, *repository.GetUserByAccessToken(authHeader), "Admin")
	c.IndentedJSON(http.StatusOK, buildCreateOrganizationResponse(org))
}

func ReadOrganizationHandler(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if !utils.IsAuthorized(authHeader) {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	ID := c.Param("organization_id")
	if !utils.IsValidOrganizationID(ID) {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "The organization id is not valid"})
		return
	}

	org := repository.GetOrganizationById(ID)

	c.IndentedJSON(http.StatusOK, buildReadOrganizationResponse(*org))
}

func ReadAllOrganizationHandler(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if !utils.IsAuthorized(authHeader) {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	usr := repository.GetUserByAccessToken(authHeader)
	orgs := repository.GetUserOrgs(*usr)
	c.IndentedJSON(http.StatusOK, buildReadAllOrganizationResponse(orgs))
}

func UpdateOrganizationHandler(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if !utils.IsAuthorized(authHeader) {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var updateOrganizationData models.CreateOrganizationModel
	err := c.BindJSON(&updateOrganizationData)
	ID := c.Param("organization_id")
	if err != nil || !utils.IsValidOrganizationID(ID) {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "The organization id is not valid"})
		return
	}

	org := repository.GetOrganizationById(ID)
	usr := repository.GetUserByAccessToken(authHeader)
	if !repository.IsMember(*usr, *org) {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "You are not a member of this organization"})
		return
	}

	org.Name = updateOrganizationData.Name
	org.Description = updateOrganizationData.Description
	repository.UpdateOrganization(org)
	c.IndentedJSON(http.StatusOK, buildUpdateOrganizationResponse(*org))
}

func InviteUserHandler(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if !utils.IsAuthorized(authHeader) {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var inviteData models.InviteModel
	err := c.BindJSON(&inviteData)
	ID := c.Param("organization_id")
	if err != nil || !utils.IsValidOrganizationID(ID) {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "The organization id is not valid"})
		return
	}

	org := repository.GetOrganizationById(ID)
	usr := repository.GetUserByAccessToken(authHeader)
	if !repository.IsMember(*usr, *org) {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "You are not a member of this organization"})
		return
	}

	invitedUser := repository.GetUserByEmail(inviteData.Email)
	if invitedUser == nil || invitedUser.Email == "" {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "There is no user with such email"})
		return
	}

	if repository.IsMember(*invitedUser, *org) {
		c.IndentedJSON(http.StatusAlreadyReported, gin.H{"message": "The user is already a member"})
		return
	}

	repository.AddMemberToOrganization(org, *invitedUser, "Invited")
	c.IndentedJSON(http.StatusOK, buildInviteUserResponse())
}

func DeleteOrganizationHandler(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if !utils.IsAuthorized(authHeader) {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	ID := c.Param("organization_id")
	if !utils.IsValidOrganizationID(ID) {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "The organization id is not valid"})
		return
	}

	org := repository.GetOrganizationById(ID)
	usr := repository.GetUserByAccessToken(authHeader)
	if !repository.IsMember(*usr, *org) {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "You are not a member of this organization"})
		return
	}

	repository.DeleteOrganization(ID)
	c.IndentedJSON(http.StatusOK, buildDeleteOrganizationResponse())
}
