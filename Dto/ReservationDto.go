package Dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type ReservationDto struct {
	UserID    primitive.ObjectID   `json:"userID"`
	AssetIDs  []primitive.ObjectID `json:"assets"`
	StartDate primitive.DateTime   `json:"startDate"`
	EndDate   primitive.DateTime   `json:"endDate"`
}

type ReportDto struct {
	Description string `json:"description"`
}

type GradeDto struct {
	Grade       int64    `json:"grade"`
	Description []string `json:"description"`
}

type ReservationviewDto struct {
	UserID      primitive.ObjectID `json:"userID"`
	AssetName   string             `json:"assets"`
	CompanyName string
	StartDate   primitive.DateTime `json:"startDate"`
	EndDate     primitive.DateTime `json:"endDate"`
}
