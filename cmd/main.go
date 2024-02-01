package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

var collection *mongo.Collection
var users []user
var organizations []organization

type organizationMember struct {
	Name        string `bson:"name" json:"name"`
	Email       string `bson:"email" json:"email"`
	AccessLevel string `bson:"access_level" json:"access_level"`
}

type organization struct {
	ID          string               `bson:"organization_id" json:"organization_id"`
	Name        string               `bson:"name" json:"name"`
	Description string               `bson:"description" json:"description"`
	Members     []organizationMember `bson:"organization_members,omitempty" json:"organization_members,omitempty"`
}

type user struct {
	Name           string   `bson:"name" json:"name"`
	Email          string   `bson:"email" json:"email"`
	HashedPassword string   `bson:"password" json:"password"`
	AccessToken    string   `bson:"access_token" json:"access_token"`
	RefreshToken   string   `bson:"refresh_token" json:"refresh_token"`
	Organizations  []string `bson:"organizations" json:"organizations"`
}

type signupModel struct {
	Name     string `bson:"name" json:"name"`
	Email    string `bson:"email" json:"email"`
	Password string `bson:"password" json:"password"`
}

type signupResponse struct {
	Message string `bson:"message" json:"message"`
}

func signupHandler(c *gin.Context) {
	var signupData signupModel
	err := c.BindJSON(&signupData)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if !isValidSignup(signupData) {
		c.IndentedJSON(http.StatusConflict, gin.H{"message": "The Username or password are not valid."})
		return
	}

	usr := createUser(signupData)
	users = append(users, usr)
	c.IndentedJSON(http.StatusOK, buildSignupResponse())
}

func buildSignupResponse() signupResponse {
	var response signupResponse
	response.Message = "Sign Up Succeeded."
	return response
}

func isValidSignup(signupData signupModel) bool {
	for _, usr := range users {
		if usr.Email == signupData.Email {
			return false
		}
	}
	return true
}

// ///////
type signinModel struct {
	Email    string `bson:"email" json:"email"`
	Password string `bson:"password" json:"password"`
}

type signinResponse struct {
	Message      string `bson:"message" json:"message"`
	AccessToken  string `bson:"access_token" json:"access_token"`
	RefreshToken string `bson:"refresh_token" json:"refresh_token"`
}

func signinHandler(c *gin.Context) {
	var signinData signinModel
	err := c.BindJSON(&signinData)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if !isValidSignin(signinData) {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "The Username or password are not correct."})
		return
	}

	usr := getUserByEmail(signinData.Email)
	accessToken, refreshToken := updateSession(usr)
	usr.AccessToken = accessToken
	usr.RefreshToken = refreshToken
	c.IndentedJSON(http.StatusOK, buildSigninResponse(usr))
}

func updateSession(usr *user) (string, string) {
	accessToken := "Bearer " + uuid.New().String()
	refreshToken := uuid.New().String()
	return accessToken, refreshToken
}

func getUserByEmail(email string) *user {
	for i, usr := range users {
		if usr.Email == email {
			return &users[i]
		}
	}
	return nil
}

func buildSigninResponse(usr *user) signinResponse {
	var response signinResponse
	response.Message = "Sign In Succeeded."
	response.AccessToken = usr.AccessToken
	response.RefreshToken = usr.RefreshToken
	return response
}
func createUser(signupData signupModel) user {
	var usr user
	usr.Name = signupData.Name
	usr.Email = signupData.Email
	usr.HashedPassword = signupData.Password
	usr.AccessToken, usr.RefreshToken = updateSession(&usr)
	return usr
}

func isValidSignin(signinData signinModel) bool {
	usr := getUserByEmail(signinData.Email)
	if usr == nil || usr.HashedPassword != signinData.Password {
		return false
	}
	return true
}

// ////
type refreshModel struct {
	RefreshToken string `bson:"refresh_token" json:"refresh_token"`
}

type refreshResponse struct {
	Message      string `bson:"message" json:"message"`
	AccessToken  string `bson:"access_token" json:"access_token"`
	RefreshToken string `bson:"refresh_token" json:"refresh_token"`
}

func refreshHandler(c *gin.Context) {
	var refreshData refreshModel
	err := c.BindJSON(&refreshData)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if !isValidRefreshToken(refreshData.RefreshToken) {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "The refresh token is not valid."})
		return
	}

	usr := getUserByRefreshToken(refreshData.RefreshToken)
	accessToken, refreshToken := updateSession(usr)
	usr.AccessToken = accessToken
	usr.RefreshToken = refreshToken
	c.IndentedJSON(http.StatusOK, buildRefreshResponse(usr))
}

func isValidRefreshToken(refreshToken string) bool {
	if refreshToken == "" {
		return false
	}

	for _, usr := range users {
		if usr.RefreshToken == refreshToken {
			return true
		}
	}
	return false
}

func buildRefreshResponse(usr *user) refreshResponse {
	var response refreshResponse
	response.Message = "Refreshing Tokens Succeeded."
	response.AccessToken = usr.AccessToken
	response.RefreshToken = usr.RefreshToken
	return response
}

func getUserByRefreshToken(refreshToken string) *user {
	for _, usr := range users {
		if usr.RefreshToken == refreshToken {
			return &usr
		}
	}
	return nil
}

//////

type createOrganizationModel struct {
	Name        string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
}

type createOrganizationResponse struct {
	OrganizationID string `bson:"organization_id" json:"organization_id"`
}

func createOrganizationHandler(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if !authorized(authHeader) {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var createOrganizationData createOrganizationModel
	err := c.BindJSON(&createOrganizationData)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if !isValidOrganization(createOrganizationData.Name) {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "The organization name is not valid"})
		return
	}

	var org organization
	org.Name = createOrganizationData.Name
	org.Description = createOrganizationData.Description
	org.ID = strconv.Itoa(len(organizations))
	addMemberToOrganization(&org, *getUserByAccessToken(authHeader), "Admin")
	organizations = append(organizations, org)
	c.IndentedJSON(http.StatusOK, buildCreateOrganizationResponse(org))
}

func getUserByAccessToken(accessToken string) *user {
	for i, usr := range users {
		if usr.AccessToken == accessToken {
			return &users[i]
		}
	}
	return nil
}

func authorized(authHeader string) bool {
	for _, usr := range users {
		fmt.Println(usr.AccessToken)
		if usr.AccessToken == authHeader {
			return true
		}
	}
	return false
}

func isValidOrganization(name string) bool {
	if name == "" {
		return false
	}

	for _, org := range organizations {
		if org.Name == name {
			return false
		}
	}
	return true
}

func buildCreateOrganizationResponse(org organization) createOrganizationResponse {
	var response createOrganizationResponse
	response.OrganizationID = org.ID
	return response
}

// ////
func readOrganizationHandler(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if !authorized(authHeader) {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	ID := c.Param("organization_id")
	if !isValidOrganizationID(ID) {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "The organization id is not valid"})
		return
	}

	org := getOrganizationById(ID)

	c.IndentedJSON(http.StatusOK, buildReadOrganizationResponse(*org))
}

func getOrganizationById(ID string) *organization {
	for _, org := range organizations {
		if org.ID == ID {
			return &org
		}
	}
	return nil
}

func isValidOrganizationID(ID string) bool {
	if ID == "" {
		return false
	}

	for _, org := range organizations {
		if org.ID == ID {
			return true
		}
	}

	return false
}

func buildReadOrganizationResponse(org organization) organization {
	return org
}

func addMemberToOrganization(org *organization, usr user, accessLevel string) {
	var member organizationMember
	member.Email = usr.Email
	member.Name = usr.Name
	member.AccessLevel = accessLevel
	org.Members = append(org.Members, member)
}

///////

func readAllOrganizationHandler(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if !authorized(authHeader) {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	usr := getUserByAccessToken(authHeader)
	orgs := getUserOrgs(*usr)
	c.IndentedJSON(http.StatusOK, buildReadAllOrganizationResponse(orgs))
}

func getUserOrgs(usr user) []organization {
	var orgs []organization
	for _, org := range organizations {
		if isMember(usr, org) {
			orgs = append(orgs, org)
		}
	}
	return orgs
}

func isMember(usr user, org organization) bool {
	for _, mem := range org.Members {
		if mem.Email == usr.Email {
			return true
		}
	}
	return false
}

func buildReadAllOrganizationResponse(org []organization) []organization {
	return org
}

// /////
func updateOrganizationHandler(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if !authorized(authHeader) {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var updateOrganizationData createOrganizationModel
	err := c.BindJSON(&updateOrganizationData)
	ID := c.Param("organization_id")
	if err != nil || !isValidOrganizationID(ID) {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "The organization id is not valid"})
		return
	}

	org := getOrganizationById(ID)
	usr := getUserByAccessToken(authHeader)
	if !isMember(*usr, *org) {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "You are not a member of this organization"})
		return
	}

	org.Name = updateOrganizationData.Name
	org.Description = updateOrganizationData.Description
	c.IndentedJSON(http.StatusOK, buildUpdateOrganizationResponse(*org))
}

func buildUpdateOrganizationResponse(org organization) organization {
	org.Members = nil
	return org
}

// ////
type inviteModel struct {
	Email string `bson:"user_email" json:"user_email"`
}

type inviteResponse struct {
	Message string `bson:"message" json:"message"`
}

func inviteUserHandler(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if !authorized(authHeader) {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var inviteData inviteModel
	err := c.BindJSON(&inviteData)
	ID := c.Param("organization_id")
	if err != nil || !isValidOrganizationID(ID) {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "The organization id is not valid"})
		return
	}

	org := getOrganizationById(ID)
	usr := getUserByAccessToken(authHeader)
	if !isMember(*usr, *org) {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "You are not a member of this organization"})
		return
	}

	invitedUser := getUserByEmail(inviteData.Email)
	if invitedUser == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "There is no user with such email"})
		return
	}

	if isMember(*invitedUser, *org) {
		c.IndentedJSON(http.StatusAlreadyReported, gin.H{"message": "The user is already a member"})
		return
	}

	var member organizationMember
	member.Name = invitedUser.Name
	member.Email = invitedUser.Email
	member.AccessLevel = "Invited"
	org.Members = append(org.Members, member)
	c.IndentedJSON(http.StatusOK, buildInviteUserResponse())
}

func buildInviteUserResponse() inviteResponse {
	var response inviteResponse
	response.Message = "User Invited"
	return response
}

// ////
func main() {

	router := gin.Default()
	router.POST("/signup", signupHandler)
	router.POST("/signin", signinHandler)
	router.POST("/refresh-token", refreshHandler)
	router.POST("/organization", createOrganizationHandler)
	router.GET("/organization/:organization_id", readOrganizationHandler)
	router.GET("/organization", readAllOrganizationHandler)
	router.PUT("/organization/:organization_id", updateOrganizationHandler)
	// router.DELETE("/organization/:organization_id", deleteOrganizationHandler)
	router.POST("/organization/:organization_id/invite", inviteUserHandler)

	router.Run("localhost:8080")
}
