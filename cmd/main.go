package main

import (
	"fmt"

	"organization-manager/pkg/database/mongodb/repository"
)

func main() {
	//routes.SetupRouter()

	repository.InitDatabase()
	fmt.Println(repository.GetUserByEmail("louai0nasr@gmail.com").Name)
	fmt.Println(repository.GetUserByRefreshToken("abc").Name)
	fmt.Println(repository.GetUserByAccessToken("abc").Name)
	fmt.Println(repository.GetOrganizationById("100").Name)
	fmt.Println(repository.GetUserOrgs(*repository.GetUserByEmail("louai0nasr@gmail.com"))[0].Name)
	repository.AddMemberToOrganization(repository.GetOrganizationById("100"), *repository.GetUserByEmail("louai0nasr@gmail.com"), "Invited")
	repository.DeleteOrganization("100")
}
