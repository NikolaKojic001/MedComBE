package Service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"unicode"

	dto "main/Dto"
	model "main/Model"
	repository "main/Repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ApplicationRegister(email string, firstname string, lastname string, phonenumber string, password string, passwordAgain string, city string, country string, company primitive.ObjectID, role string, profession string, res http.ResponseWriter) {
	fmt.Println(email, firstname, lastname, phonenumber, password, passwordAgain, city, country, company, role, profession)
	if email == "" || city == "" || password == "" || country == "" || phonenumber == "" || firstname == "" || lastname == "" {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Some required parameters are missing.")
		return
	}
	if password != passwordAgain {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Passwords not match.")
		return
	}
	if role != "Admin" && role != "CompanyAdmin" && role != "User" {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Role not exist.")
		return
	}
	Emailmatch, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, email)
	if !Emailmatch {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Invalid email format.")
		return
	}
	if len(password) < 6 {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Password must contain at least 6 charracters.")
		return
	}
	if !containsSpecialCharacters(password) {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Password must contain a special charracter!")
		return
	}
	HasUpper := false
	HasLower := true
	for _, r := range password {
		if unicode.IsUpper(r) {
			HasUpper = true
		}
		if unicode.IsLower(r) {
			HasLower = true
		}
	}
	if !(HasUpper && HasLower) {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Password must contain uppercase and lowercase letters!")
		return
	}
	//Save user
	if repository.FindUserEmail(email) {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Email in use")
		return
	}

	hashedPass, hashError := repository.HashPassword(password)
	if hashError != nil {
		log.Panic(hashError)
	}
	if role == "User" {
		id := repository.SaveUserApplication(email, firstname, lastname, phonenumber, hashedPass, city, country)
		SendConfirmationMail(id.Hex(), email)
	}
	if role == "CompanyAdmin" {
		id := repository.SaveCompanyAdminApplication(email, firstname, lastname, phonenumber, hashedPass, city, country, company, role, profession)
		SendConfirmationMail(id.Hex(), email)
	}
	if role == "Admin" {
		id := repository.SaveAdminApplication(email, firstname, lastname, phonenumber, hashedPass, city, country, role)
		SendConfirmationMail(id.Hex(), email)
	}

}
func VerifyUser(userID primitive.ObjectID) bool {
	return repository.VerifyUser(userID)
}

func containsSpecialCharacters(input string) bool {
	pattern := regexp.MustCompile(`[!@#$%^&*()?><,./|\}{=-_+]`)
	if match := pattern.FindStringIndex(input); match != nil {
		return true
	}

	return false
}

func GetAllUserReports(userID primitive.ObjectID, res http.ResponseWriter) ([]model.Report, bool) {
	return repository.GetAllUserReports(userID, res)
}

func SaveLoyaltyProgram(requestBody dto.LoyaltyProgram, res http.ResponseWriter) bool {
	if requestBody.Name == "" || requestBody.Discount <= 0 || requestBody.UpValue <= requestBody.DownValue {
		fmt.Println(requestBody)
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Some paramtears is not valid")
		return false
	}
	return repository.SaveLoyaltyProgram(requestBody.Name, requestBody.Discount, requestBody.UpValue, requestBody.DownValue, res)
}

func SetLoyaltyProgram(userObjectID primitive.ObjectID, res http.ResponseWriter) bool {
	return repository.SetLoyaltyProgram(userObjectID, res)
}
