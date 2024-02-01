package models

type OrganizationMember struct {
	Name        string `bson:"name" json:"name"`
	Email       string `bson:"email" json:"email"`
	AccessLevel string `bson:"access_level" json:"access_level"`
}

type Organization struct {
	ID          string               `bson:"organization_id" json:"organization_id"`
	Name        string               `bson:"name" json:"name"`
	Description string               `bson:"description" json:"description"`
	Members     []OrganizationMember `bson:"organization_members,omitempty" json:"organization_members,omitempty"`
}

type User struct {
	Name           string   `bson:"name" json:"name"`
	Email          string   `bson:"email" json:"email"`
	HashedPassword string   `bson:"password" json:"password"`
	AccessToken    string   `bson:"access_token" json:"access_token"`
	RefreshToken   string   `bson:"refresh_token" json:"refresh_token"`
	Organizations  []string `bson:"organizations" json:"organizations"`
}

type SignupModel struct {
	Name     string `bson:"name" json:"name"`
	Email    string `bson:"email" json:"email"`
	Password string `bson:"password" json:"password"`
}

type SignupResponse struct {
	Message string `bson:"message" json:"message"`
}

type SigninModel struct {
	Email    string `bson:"email" json:"email"`
	Password string `bson:"password" json:"password"`
}

type SigninResponse struct {
	Message      string `bson:"message" json:"message"`
	AccessToken  string `bson:"access_token" json:"access_token"`
	RefreshToken string `bson:"refresh_token" json:"refresh_token"`
}

type RefreshModel struct {
	RefreshToken string `bson:"refresh_token" json:"refresh_token"`
}

type RefreshResponse struct {
	Message      string `bson:"message" json:"message"`
	AccessToken  string `bson:"access_token" json:"access_token"`
	RefreshToken string `bson:"refresh_token" json:"refresh_token"`
}

type CreateOrganizationModel struct {
	Name        string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
}

type CreateOrganizationResponse struct {
	OrganizationID string `bson:"organization_id" json:"organization_id"`
}

type InviteModel struct {
	Email string `bson:"user_email" json:"user_email"`
}

type InviteResponse struct {
	Message string `bson:"message" json:"message"`
}
