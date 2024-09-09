package Model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Reservation struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty"`
	AssetIDs  []primitive.ObjectID `bson:"assetID"`
	UserID    primitive.ObjectID   `bson:"userID"`
	CompanyID primitive.ObjectID   `bson:"companyID"`
	StartDate primitive.DateTime   `bson:"startDate"`
	EndDate   primitive.DateTime   `bson:"endDate"`
	Canceled  bool                 `bson:"canceled"`
}

type Appointment struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty"`
	AssetIDs  []primitive.ObjectID `bson:"assetID"`
	UserID    primitive.ObjectID   `bson:"userID"`
	CompanyID primitive.ObjectID   `bson:"companyID"`
	Date      primitive.DateTime   `bson:"startDate"`
}

type Report struct {
	ID          primitive.ObjectID `bson:"_id"`
	UserID      primitive.ObjectID `bson:"userid"`
	CompanyID   primitive.ObjectID `bson:"companyid"`
	Description string             `bson:"description"`
	Replay      string             `bson:"replay"`
}

type Grade struct {
	ID          primitive.ObjectID `bson:"_id"`
	Grade       int64              `bson:"grade"`
	UserID      primitive.ObjectID `bson:"userid"`
	CompanyID   primitive.ObjectID `bson:"companyid"`
	Description []string           `bson:"description"`
}
