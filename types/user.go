package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost         = 12
	minFirstNameLength = 2
	minLastNameLength  = 2
	minPasswordLength  = 7
)

type CreateUserRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (params CreateUserRequest) Validate() []string {
	errors := []string{}
	if len(params.FirstName) < minFirstNameLength {
		errors = append(errors, fmt.Sprintf("first name must be at least %d characters", minFirstNameLength))
	}
	if len(params.LastName) < minLastNameLength {
		errors = append(errors, fmt.Sprintf("last name must be at least %d characters", minLastNameLength))
	}
	if !isEmailValid(params.Email) {
		errors = append(errors, fmt.Sprintf("email is invalid"))
	}
	if len(params.Password) < minPasswordLength {
		errors = append(errors, fmt.Sprintf("password must be at least %d characters", minPasswordLength))
	}
	return errors
}

func isEmailValid(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)
}

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"` //for marshalling and unmarshalling
	FirstName         string             `bson:"firstName" json:"firstName"`
	LastName          string             `bson:"lastName" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"encryptedPassword" json:"-"`
}

func CreateUserRequestToUser(req *CreateUserRequest) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcryptCost)
	if err != nil {
		return nil, err
	}

	return &User{
		FirstName:         req.FirstName,
		LastName:          req.LastName,
		Email:             req.Email,
		EncryptedPassword: string(encpw),
	}, nil
}
