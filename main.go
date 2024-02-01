package main

import (
	"organization-manager/pkg/api/routes"
	"organization-manager/pkg/database/mongodb/repository"
)

func main() {
	repository.InitDatabase()
	routes.SetupRouter()
}
