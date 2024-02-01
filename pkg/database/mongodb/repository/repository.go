package repository

import (
	"context"
	"fmt"
	"os"
	"time"

	"organization-manager/pkg/database/mongodb/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var UsersCollection *mongo.Collection
var OrganizationsCollection *mongo.Collection
var ctx context.Context
var cancel context.CancelFunc

func InitDatabase() {
	fmt.Println("Starting the DB connection...")
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(os.Getenv("MONGO_URL")).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(err)
	}

	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(ctx, bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	databaseName := "organization_manager"
	UsersCollection = client.Database(databaseName).Collection("users")
	OrganizationsCollection = client.Database(databaseName).Collection("organizations")
}

func AddUser(user models.User) {
	UsersCollection.InsertOne(ctx, user)
}

func AddOrganization(organization models.Organization) {
	UsersCollection.InsertOne(ctx, organization)
}

func AddMemberToOrganization(org *models.Organization, usr models.User, accessLevel string) {
	var member models.OrganizationMember
	member.Email = usr.Email
	member.Name = usr.Name
	member.AccessLevel = accessLevel

	filter := bson.M{"organization_id": org.ID}
	update := bson.M{
		"$push": bson.M{"organization_members": member},
	}

	// Perform the update operation
	result, err := OrganizationsCollection.UpdateOne(ctx, filter, update)
	if err == nil {
		fmt.Println(result.ModifiedCount)
	}
}

func GetUserByEmail(email string) *models.User {
	cursor, err := UsersCollection.Find(ctx, bson.M{"email": email})
	if err != nil {
		fmt.Println(err.Error())
	} else {
		var user models.User
		cursor.Next(ctx)
		cursor.Decode(&user)
		return &user
	}
	return nil
}

func GetUserByRefreshToken(refreshToken string) *models.User {
	cursor, err := UsersCollection.Find(ctx, bson.M{"refresh_token": refreshToken})
	if err != nil {
		fmt.Println(err.Error())
	} else {
		var user models.User
		cursor.Next(ctx)
		cursor.Decode(&user)
		return &user
	}
	return nil
}

func GetUserByAccessToken(accessToken string) *models.User {
	cursor, err := UsersCollection.Find(ctx, bson.M{"access_token": accessToken})
	if err != nil {
		fmt.Println(err.Error())
	} else {
		var user models.User
		cursor.Next(ctx)
		cursor.Decode(&user)
		return &user
	}
	return nil
}

func GetOrganizationById(ID string) *models.Organization {
	cursor, err := OrganizationsCollection.Find(ctx, bson.M{"organization_id": ID})
	if err != nil {
		fmt.Println(err.Error())
	} else {
		var organization models.Organization
		cursor.Next(ctx)
		cursor.Decode(&organization)
		return &organization
	}
	return nil
}

func GetOrganizationByName(name string) *models.Organization {
	cursor, err := OrganizationsCollection.Find(ctx, bson.M{"name": name})
	if err != nil {
		fmt.Println(err.Error())
	} else {
		var organization models.Organization
		cursor.Next(ctx)
		cursor.Decode(&organization)
		return &organization
	}
	return nil
}

func GetUserOrgs(usr models.User) []models.Organization {
	filter := bson.M{
		"organization_members": bson.M{
			"$elemMatch": bson.M{
				"email": usr.Email,
			},
		},
	}
	cursor, err := OrganizationsCollection.Find(ctx, filter)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		var organizations []models.Organization
		err = cursor.All(ctx, &organizations)
		if err == nil {
			return organizations
		}
	}
	return nil
}

func DeleteOrganization(ID string) {
	OrganizationsCollection.DeleteOne(ctx, bson.M{"organization_id": ID})
}

func IsMember(usr models.User, org models.Organization) bool {
	for _, mem := range org.Members {
		if mem.Email == usr.Email {
			return true
		}
	}
	return false
}
