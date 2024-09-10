package Controller

import (
	"encoding/json"
	dto "main/Dto"
	service "main/Service"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CompanyController struct {
	Router *mux.Router
}

func NewCompanyController() *CompanyController {
	return &CompanyController{
		Router: mux.NewRouter(),
	}
}

func (uc *CompanyController) RegisterRoutes() {
	uc.Router.HandleFunc("/companies/save", uc.SaveCompany)
	uc.Router.HandleFunc("/companies/report", uc.ReportCompany)
	uc.Router.HandleFunc("/companies/grade", uc.GradeCompany)
	uc.Router.HandleFunc("/companies/get/all", uc.GetAllCompanies)
	uc.Router.HandleFunc("/companies/get/by/id", uc.GetCompanyById)
	uc.Router.HandleFunc("/companies/get/all/reports", uc.GetAllCompanyReports)
	uc.Router.HandleFunc("/companies/get/all/grades", uc.GetAllUserGrades)
	uc.Router.HandleFunc("/companies/get/report/by/id/{reportID}", uc.GetReportByID)
	uc.Router.HandleFunc("/companies/get/grade/by/id/{gradeID}", uc.GetGradeByID)
	uc.Router.HandleFunc("/companies/report/replay/{reportID}", uc.ReportReplay)
	uc.Router.HandleFunc("/companies/grade/update/{gradeID}", uc.UpdateGrade)

}

func (uc *CompanyController) SaveCompany(res http.ResponseWriter, req *http.Request) {
	var requestBody dto.CompanyDto
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

	done := service.SaveCompany(requestBody, res)
	if done {
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode("Company successfuly created")
	} else {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error while creating compnay")
	}
}

func (uc *CompanyController) ReportCompany(res http.ResponseWriter, req *http.Request) {
	companyID := req.Header.Get("companyID")
	companyObjectID, companyIDErr := primitive.ObjectIDFromHex(companyID)
	if companyIDErr != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error with company id")
		return
	}

	var requestBody dto.ReportDto
	err := json.NewDecoder(req.Body).Decode(&requestBody)
	if err != nil {
		http.Error(res, "Error decoding request body", http.StatusBadRequest)
		return
	}

	authHeader := req.Header.Get("Authorization")
	user, pointer := service.GetUserFromToken(res, req, authHeader)
	if pointer == nil || user.Role != "User" {
		res.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(res).Encode("Unauthorized")
		return
	}

	done := service.ReportCompany(user.ID, companyObjectID, requestBody, res)
	if done {
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode("Company successfuly reported")
	} else {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error while reporting compnay or you didn't have past reservation in this company")
	}
}

func (uc *CompanyController) GradeCompany(res http.ResponseWriter, req *http.Request) {
	companyID := req.Header.Get("companyID")
	companyObjectID, companyIDErr := primitive.ObjectIDFromHex(companyID)
	if companyIDErr != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error with company id")
		return
	}

	var requestBody dto.GradeDto
	err := json.NewDecoder(req.Body).Decode(&requestBody)
	if err != nil {
		http.Error(res, "Error decoding request body", http.StatusBadRequest)
		return
	}

	authHeader := req.Header.Get("Authorization")
	user, pointer := service.GetUserFromToken(res, req, authHeader)
	if pointer == nil || user.Role != "User" {
		res.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(res).Encode("Unauthorized")
		return
	}

	done := service.GradeCompany(user.ID, companyObjectID, requestBody, res)
	if done {
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode("Company successfuly graded")
	} else {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error while grading compnay or you didn't have past reservation in this company")
	}
}

func (uc *CompanyController) GetCompanyById(res http.ResponseWriter, req *http.Request) {
	queryParams := req.URL.Query()
	companyID := queryParams.Get("companyID")
	companyObjectID, companyIDErr := primitive.ObjectIDFromHex(companyID)
	if companyIDErr != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error with company id")
		return
	}

	company, done := service.GetCompanyById(companyObjectID, res)
	if done {
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(company)
	} else {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error while getting company")
	}
}

func (uc *CompanyController) GetAllCompanies(res http.ResponseWriter, req *http.Request) {
	companies, err := service.GetAllCompanies()
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error while finding companies")
	} else {
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(companies)
	}

}

func (uc *CompanyController) GetAllCompanyReports(res http.ResponseWriter, req *http.Request) {
	authHeader := req.Header.Get("Authorization")
	user, pointer := service.GetUserFromToken(res, req, authHeader)
	if pointer == nil || user.Role != "Admin" {
		res.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(res).Encode("Unauthorized")
		return
	}

	reports, err := service.GetAllCompanyReports(res)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error while finding reports")
	} else {
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(reports)
	}

}

func (uc *CompanyController) GetAllUserGrades(res http.ResponseWriter, req *http.Request) {
	authHeader := req.Header.Get("Authorization")
	user, pointer := service.GetUserFromToken(res, req, authHeader)
	if pointer == nil || user.Role != "User" {
		res.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(res).Encode("Unauthorized")
		return
	}

	grades, err := service.GetAllUserGrades(user.ID, res)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error while finding companies")
	} else {
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(grades)
	}

}

func (uc *CompanyController) GetReportByID(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	reportID := vars["reportID"]

	reportObjectID, userIDErr := primitive.ObjectIDFromHex(reportID)
	if userIDErr != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error with report id")
		return
	}

	report, found := service.GetReportByID(reportObjectID, res)
	if !found {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error while finding report")
	} else {
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(report)
	}

}

func (uc *CompanyController) GetGradeByID(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	gradeID := vars["gradeID"]

	gradeObjectID, gradeIDErr := primitive.ObjectIDFromHex(gradeID)
	if gradeIDErr != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error with report id")
		return
	}

	report, found := service.GetGradeByID(gradeObjectID, res)
	if !found {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error while finding report")
	} else {
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(report)
	}

}

func (uc *CompanyController) ReportReplay(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	reportID := vars["reportID"]

	var requestBody dto.ReportDto
	err := json.NewDecoder(req.Body).Decode(&requestBody)
	if err != nil {
		http.Error(res, "Error decoding request body", http.StatusBadRequest)
		return
	}

	reportObjectID, userIDErr := primitive.ObjectIDFromHex(reportID)
	if userIDErr != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error with report id")
		return
	}

	authHeader := req.Header.Get("Authorization")
	user, pointer := service.GetUserFromToken(res, req, authHeader)
	if pointer == nil || user.Role != "Admin" {
		res.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(res).Encode("Unauthorized")
		return
	}

	report, found := service.ReportReplay(reportObjectID, requestBody.Description, res)
	if !found {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error while finding report")
	} else {
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(report)
		service.SendReportReplay(report.ID, res)

	}
}

func (ac *CompanyController) UpdateGrade(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	gradeHexID := vars["gradeID"]

	authHeader := req.Header.Get("Authorization")
	admin, pointer := service.GetUserFromToken(res, req, authHeader)
	if pointer == nil || admin.Role != "User" {
		res.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(res).Encode("Unauthorized")
		return
	}

	gradeID, gradeIDErr := primitive.ObjectIDFromHex(gradeHexID)
	if gradeIDErr != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error with company ID")
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(req.Body).Decode(&updates); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Invalid request payload")
		return
	}

	done := service.UpdateGrade(gradeID, updates)
	if done {
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode("Successfully updated")
	} else {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Bad request")
	}
}
