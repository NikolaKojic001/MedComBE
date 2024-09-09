package Repository

import (
	"context"
	"log"
	dto "main/Dto"
	model "main/Model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func SaveAsset(asset dto.AssetDto, company model.Company) bool {
	collection := GetClient().Database("MedicinskaOprema").Collection("Companies")

	newAsset := model.Asset{
		ID:   primitive.NewObjectID(),
		Name: asset.Name,
	}
	filter := bson.M{"_id": company.ID}

	update := bson.M{
		"$push": bson.M{
			"assets": newAsset,
		},
	}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Println("Error updating company with new asset:", err)
		return false
	}

	return true
}

func GetAllAssets(companyID primitive.ObjectID) ([]model.Asset, bool) {
	collection := GetClient().Database("MedicinskaOprema").Collection("Companies")
	filter := bson.M{"_id": companyID}
	var company model.Company

	err := collection.FindOne(context.Background(), filter).Decode(&company)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("No company found with the specified ID.")
		} else {
			log.Println("Error retrieving company:", err)
		}
		return nil, false
	}

	return company.Assets, true
}

func GetAssetById(companyID primitive.ObjectID, assetID primitive.ObjectID) (model.Asset, bool) {
	collection := GetClient().Database("MedicinskaOprema").Collection("Companies")
	filter := bson.M{"_id": companyID}
	var company model.Company

	err := collection.FindOne(context.Background(), filter).Decode(&company)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("No company found with the specified ID.")
		} else {
			log.Println("Error retrieving company:", err)
		}
		var assetNil model.Asset
		return assetNil, false
	}

	for _, asset := range company.Assets {
		if asset.ID == assetID {
			return asset, true
		}
	}

	var assetNil model.Asset
	return assetNil, false
}
