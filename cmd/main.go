package main

import (
	"organization-manager/pkg/api/routes"
	"organization-manager/pkg/database/mongodb/repository"
	//"organization-manager/pkg/utils"
)

func main() {
	repository.InitDatabase()
	//utils.IsValidOrganization("Org")
	routes.SetupRouter()
}
