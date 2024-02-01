package routes

import (
	"organization-manager/pkg/api/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() {
	router := gin.Default()
	router.POST("/signup", handlers.SignupHandler)
	router.POST("/signin", handlers.SigninHandler)
	router.POST("/refresh-token", handlers.RefreshHandler)
	router.POST("/organization", handlers.CreateOrganizationHandler)
	router.GET("/organization/:organization_id", handlers.ReadOrganizationHandler)
	router.GET("/organization", handlers.ReadAllOrganizationHandler)
	router.PUT("/organization/:organization_id", handlers.UpdateOrganizationHandler)
	// router.DELETE("/organization/:organization_id", deleteOrganizationHandler)
	router.POST("/organization/:organization_id/invite", handlers.InviteUserHandler)

	router.Run("localhost:8080")
}
