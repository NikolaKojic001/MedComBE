package Service

import (
	"encoding/json"
	dto "main/Dto"
	model "main/Model"
	repository "main/Repository"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SaveCompany(company dto.CompanyDto, res http.ResponseWriter) bool {
	if company.Name == "" || company.LocationCompany.City == "" || company.LocationCompany.Country == "" {
		json.NewEncoder(res).Encode("Some required parameters are missing.")
		return false
	}

	_, found := repository.FindCompanyByName(company.Name)
	if found {
		json.NewEncoder(res).Encode("Company name in use")
		return false
	}

	done, companyID := repository.SaveCompany(company)
	if !done {
		return false
	} else {
		ApplicationRegister(company.Email, company.Firstname, company.Lastname, company.Phonenumber, company.Password, company.PasswordAgain, company.LocationAdmin.City, company.LocationAdmin.Country, companyID, "CompanyAdmin", company.Profession, res)
		return true
	}
}

func GetAllCompanies() ([]model.Company, error) {
	return repository.GetAllCompanies()
}

func ReportCompany(userID primitive.ObjectID, companyID primitive.ObjectID, requestBody dto.ReportDto, res http.ResponseWriter) bool {
	company, found := repository.FindCompanyById(companyID, res)
	if !found {
		res.WriteHeader(http.StatusNotFound)
		json.NewEncoder(res).Encode("Company not found")
		return false
	}
	return repository.ReportCompany(userID, company.ID, requestBody)
}

func GradeCompany(userID primitive.ObjectID, companyID primitive.ObjectID, requestBody dto.GradeDto, res http.ResponseWriter) bool {
	company, found := repository.FindCompanyById(companyID, res)
	if !found {
		res.WriteHeader(http.StatusNotFound)
		json.NewEncoder(res).Encode("Company not found")
		return false
	}

	_, found = repository.GetPastGrade(userID, companyID, res)
	if found {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("You allready grade this company")
		return false
	}

	return repository.GradeCompany(userID, company.ID, requestBody)
}

func GetAllCompanyReports(companyID primitive.ObjectID, res http.ResponseWriter) ([]model.Report, error) {
	company, found := repository.FindCompanyById(companyID, res)
	if !found {
		res.WriteHeader(http.StatusNotFound)
		json.NewEncoder(res).Encode("Company not found")
	}
	return repository.GetAllCompanyReports(company.ID, res)
}

func GetCompanyById(companyID primitive.ObjectID, res http.ResponseWriter) (model.Company, bool) {
	company, found := repository.FindCompanyById(companyID, res)
	if !found {
		res.WriteHeader(http.StatusNotFound)
		json.NewEncoder(res).Encode("Company not found")
	}
	return company, found
}

func GetAllUserGrades(userID primitive.ObjectID, res http.ResponseWriter) ([]model.Grade, error) {
	return repository.GetAllUserGrades(userID, res)
}

func GetReportByID(reportID primitive.ObjectID, res http.ResponseWriter) (model.Report, bool) {
	return repository.GetReportByID(reportID, res)
}

func GetGradeByID(gradeID primitive.ObjectID, res http.ResponseWriter) (model.Grade, bool) {
	return repository.GetGradeByID(gradeID, res)
}
func ReportReplay(reportID primitive.ObjectID, description string, res http.ResponseWriter) (model.Report, bool) {
	return repository.ReportReplay(reportID, description, res)
}

func UpdateGrade(gradeID primitive.ObjectID, updates map[string]interface{}) bool {
	return repository.UpdateGrade(gradeID, updates)
}
