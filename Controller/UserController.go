package Controller

import (
	"encoding/json"
	"net/http"

	dto "main/Dto"
	service "main/Service"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController struct {
	Router *mux.Router
}

func NewUserController() *UserController {
	return &UserController{
		Router: mux.NewRouter(),
	}
}

func (uc *UserController) RegisterRoutes() {
	uc.Router.HandleFunc("/users/save", uc.SaveUser)
	uc.Router.HandleFunc("/users/login", uc.Login)
	uc.Router.HandleFunc("/users/role", uc.GetRole)
	uc.Router.HandleFunc("/users/get/data", uc.GetData)
	uc.Router.HandleFunc("/users/verify/{userID}", uc.Verify)
	uc.Router.HandleFunc("/users/get/all/reports", uc.GetAllReports)
	uc.Router.HandleFunc("/users/save/loyalty/program", uc.SaveLoyaltyProgram)
	uc.Router.HandleFunc("/users/set/loyalty/program/{userID}", uc.SetLoyaltyProgram)

}

func (uc *UserController) SaveUser(res http.ResponseWriter, req *http.Request) {
	var requestBody dto.User
	err := json.NewDecoder(req.Body).Decode(&requestBody)
	if err != nil {
		http.Error(res, "Error decoding request body", http.StatusBadRequest)
		return
	}

	if requestBody.Role == "User" {
		service.ApplicationRegister(requestBody.Email, requestBody.Firstname, requestBody.Lastname, requestBody.Phonenumber, requestBody.Password, requestBody.PasswordAgain, requestBody.Location.City, requestBody.Location.Country, requestBody.Compmnay, requestBody.Role, requestBody.Profession, res)
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode("You successfully registrated")
	} else if requestBody.Role == "Admin" || requestBody.Role == "CompanyAdmin" {
		authHeader := req.Header.Get("Authorization")
		user, pointer := service.GetUserFromToken(res, req, authHeader)
		if pointer == nil || user.Role != "Admin" {
			res.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(res).Encode("Unauthorized")
			return
		}
		service.ApplicationRegister(requestBody.Email, requestBody.Firstname, requestBody.Lastname, requestBody.Phonenumber, requestBody.Password, requestBody.PasswordAgain, requestBody.Location.City, requestBody.Location.Country, requestBody.Compmnay, requestBody.Role, requestBody.Profession, res)
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode("You successfully registrated" + " " + requestBody.Role)
	}
}
func (uc *UserController) Login(res http.ResponseWriter, req *http.Request) {
	var requestBody dto.LoginDto
	err := json.NewDecoder(req.Body).Decode(&requestBody)
	if err != nil {
		http.Error(res, "Error decoding request body", http.StatusBadRequest)
		return
	}
	service.TokenAppLoginLogic(res, req, requestBody.Email, requestBody.Password)
}

func (uc *UserController) Verify(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	userID := vars["userID"]

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error with company id")
		return
	}

	done := service.VerifyUser(userObjectID)
	if done {
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode("User successfuly verified")
		return
	} else {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("User not verified")
	}
}

func (uc *UserController) GetAllReports(res http.ResponseWriter, req *http.Request) {
	authHeader := req.Header.Get("Authorization")
	user, pointer := service.GetUserFromToken(res, req, authHeader)
	if pointer == nil || user.Role != "User" {
		res.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(res).Encode("Unauthorized")
		return
	}

	reports, found := service.GetAllUserReports(user.ID, res)
	if !found {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error while finding companies")
	} else {
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(reports)
	}

}

func (uc *UserController) SaveLoyaltyProgram(res http.ResponseWriter, req *http.Request) {
	authHeader := req.Header.Get("Authorization")
	user, pointer := service.GetUserFromToken(res, req, authHeader)
	if pointer == nil || user.Role != "Admin" {
		res.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(res).Encode("Unauthorized")
		return
	}

	var requestBody dto.LoyaltyProgram
	err := json.NewDecoder(req.Body).Decode(&requestBody)
	if err != nil {
		http.Error(res, "Error decoding request body", http.StatusBadRequest)
		return
	}

	done := service.SaveLoyaltyProgram(requestBody, res)
	if !done {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error while cerating loyalty program")
	} else {
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode("Loyalty program created")
	}

}

func (uc *UserController) SetLoyaltyProgram(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	userID := vars["userID"]

	userObjectID, userIDErr := primitive.ObjectIDFromHex(userID)
	if userIDErr != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error with report id")
		return
	}

	done := service.SetLoyaltyProgram(userObjectID, res)
	if !done {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error while calculate loyalty program")
	} else {
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode("User updated")
	}

}

func (uc *UserController) GetRole(res http.ResponseWriter, req *http.Request) {
	authHeader := req.Header.Get("Authorization")
	user, pointer := service.GetUserFromToken(res, req, authHeader)
	if pointer == nil {
		res.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(res).Encode("Unauthorized")
		return
	}
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(user.Role)
}

func (uc *UserController) GetData(res http.ResponseWriter, req *http.Request) {
	authHeader := req.Header.Get("Authorization")
	user, pointer := service.GetUserFromToken(res, req, authHeader)
	if pointer == nil {
		res.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(res).Encode("Unauthorized")
		return
	}
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(user)
}
