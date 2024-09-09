package Service

import (
	"encoding/json"
	dto "main/Dto"
	model "main/Model"
	repository "main/Repository"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SaveAsset(asset dto.AssetDto, companyID primitive.ObjectID, res http.ResponseWriter) bool {
	company, found := repository.FindCompanyById(companyID, res)
	if !found {
		res.WriteHeader(http.StatusNotFound)
		json.NewEncoder(res).Encode("Company not found")
	}
	return repository.SaveAsset(asset, company)
}

func GetAllAssets(companyID primitive.ObjectID, res http.ResponseWriter) ([]model.Asset, bool) {
	company, found := repository.FindCompanyById(companyID, res)
	if !found {
		res.WriteHeader(http.StatusNotFound)
		json.NewEncoder(res).Encode("Company not found")
	}
	return repository.GetAllAssets(company.ID)
}

func GetAssetById(companyID primitive.ObjectID, assetID primitive.ObjectID, res http.ResponseWriter) (model.Asset, bool) {
	company, found := repository.FindCompanyById(companyID, res)
	if !found {
		res.WriteHeader(http.StatusNotFound)
		json.NewEncoder(res).Encode("Company not found")
	}
	return repository.GetAssetById(company.ID, assetID)
}
