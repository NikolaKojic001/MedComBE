package Dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type AssetDto struct {
	Name string `json:"name"`
}

type AssetsRequests struct {
	AssetIDs     []primitive.ObjectID `json:"assets"`
	NumberOfDays int64                `json:"days"`
	StartDate    primitive.DateTime   `json:"startDate"`
	EndDate      primitive.DateTime   `json:"endDate"`
}

type PosibleDate struct {
	StartDate primitive.DateTime `json:"startDate"`
	EndDate   primitive.DateTime `json:"endDate"`
}
