package Controller

import (
	"encoding/json"
	dto "main/Dto"
	service "main/Service"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AssetController struct {
	Router *mux.Router
}

func NewAssetController() *AssetController {
	return &AssetController{
		Router: mux.NewRouter(),
	}
}

func (uc *AssetController) RegisterRoutes() {
	uc.Router.HandleFunc("/assets/save", uc.SaveAsset)
	uc.Router.HandleFunc("/assets/get/all", uc.GetAllAssets)
	uc.Router.HandleFunc("/assets/get/by/id", uc.GetAssetById)

}

func (uc *AssetController) SaveAsset(res http.ResponseWriter, req *http.Request) {
	var requestBody dto.AssetDto
	err := json.NewDecoder(req.Body).Decode(&requestBody)
	if err != nil {
		http.Error(res, "Error decoding request body", http.StatusBadRequest)
		return
	}

	authHeader := req.Header.Get("Authorization")
	user, pointer := service.GetUserFromToken(res, req, authHeader)
	if pointer == nil || user.Role != "Admin" {
		res.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(res).Encode("Unauthorized")
		return
	}

	companyID := req.Header.Get("companyID")
	companyObjectID, companyIDErr := primitive.ObjectIDFromHex(companyID)
	if companyIDErr != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error with company id")
		return
	}
	done := service.SaveAsset(requestBody, companyObjectID, res)
	if done {
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode("Asset successfuly created")
	} else {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error while creating asset")
	}
}

func (uc *AssetController) GetAllAssets(res http.ResponseWriter, req *http.Request) {
	companyID := req.Header.Get("companyID")
	companyObjectID, companyIDErr := primitive.ObjectIDFromHex(companyID)
	if companyIDErr != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error with company id")
		return
	}
	assets, done := service.GetAllAssets(companyObjectID, res)
	if done {
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(assets)
	} else {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error while finding asset")
	}
}

func (uc *AssetController) GetAssetById(res http.ResponseWriter, req *http.Request) {
	companyID := req.Header.Get("companyID")
	companyObjectID, companyIDErr := primitive.ObjectIDFromHex(companyID)
	if companyIDErr != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error with company id")
		return
	}
	queryParams := req.URL.Query()
	assetID := queryParams.Get("assetID")
	assetObjectID, assetIDErr := primitive.ObjectIDFromHex(assetID)
	if assetIDErr != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error with company id")
		return
	}
	assets, done := service.GetAssetById(companyObjectID, assetObjectID, res)
	if done {
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(assets)
	} else {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error while finding asset")
	}
}
