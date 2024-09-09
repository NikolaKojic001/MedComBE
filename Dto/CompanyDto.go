package Dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type CompanyDto struct {
	Role          string             `json:"role"`
	Email         string             `json:"email"`
	Password      string             `json:"password"`
	PasswordAgain string             `json:"passwordAgain"`
	Firstname     string             `json:"firstname"`
	Lastname      string             `json:"lastname"`
	LocationAdmin Location           `json:"locationAdmin"`
	Phonenumber   string             `json:"phonenumber"`
	Profession    string             `json:"profession"`
	Compmnay      primitive.ObjectID `json:"company"`

	Name            string   `json:"name"`
	LocationCompany Location `json:"locationCompany"`
}
