package Model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Company struct {
	ID           primitive.ObjectID   `bson:"_id,omitempty"`
	AdminID      primitive.ObjectID   `bson:"adminID"`
	Name         string               `bson:"name"`
	Location     Location             `bson:"location"`
	Assets       []Asset              `bson:"assets"`
	Reservations []primitive.ObjectID `bson:"reservations"`
}
