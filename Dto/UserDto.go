package Dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Role          string             `json:"role"`
	Email         string             `json:"email"`
	Password      string             `json:"password"`
	PasswordAgain string             `json:"passwordAgain"`
	Firstname     string             `json:"firstname"`
	Lastname      string             `json:"lastname"`
	Location      Location           `json:"location"`
	Phonenumber   string             `json:"phonenumber"`
	Profession    string             `json:"profession"`
	Compmnay      primitive.ObjectID `json:"company"`
}

type Location struct {
	City    string `json:"city"`
	Country string `json:"country"`
}

type LoginDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoyaltyProgram struct {
	Name      string  `json:"name"`
	Discount  float64 `json:"discount"`
	UpValue   int32   `json:"upValue"`
	DownValue int32   `json:"downValue"`
}
