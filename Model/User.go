package Model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Email          string             `bson:"email"`
	Password       string             `bson:"password"`
	Role           string             `bson:"role"`
	Firstname      string             `bson:"firstname"`
	Lastname       string             `bson:"lastname"`
	Location       Location           `bson:"location"`
	Phonenumber    string             `bson:"phonenumber"`
	Profession     string             `bson:"profession"`
	CompmnayID     primitive.ObjectID `bson:"company"`
	Verified       bool               `bson:"verified"`
	Penalies       int32              `bson:"penalties"`
	Points         int32              `bson:"points"`
	LoyaltyProgram primitive.ObjectID `bson:"loyaltyProgram"`
}

type Location struct {
	City    string `bson:"city"`
	Country string `bson:"country"`
}

type LoyaltyProgram struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Discount  float64            `bson:"discount"`
	UpValue   int32              `bson:"upValue"`
	DownValue int32              `bson:"downValue"`
}
