package Controller

import (
	"encoding/json"
	"fmt"
	service "main/Service"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type StatisticController struct {
	Router *mux.Router
}

func NewStatisticController() *StatisticController {
	return &StatisticController{
		Router: mux.NewRouter(),
	}
}

func (uc *StatisticController) RegisterRoutes() {
	uc.Router.HandleFunc("/statistics/get/average/grade", uc.GetAverageGrade)
	uc.Router.HandleFunc("/statistics/get/reservations/by/month/{year}", uc.GetReservationsByMonth)
	uc.Router.HandleFunc("/statistics/get/reservations/by/queartals/{year}", uc.GetReservationsByQuartals)
	uc.Router.HandleFunc("/statistics/get/reservations/by/years/{startYear}/{endYear}", uc.GetReservationsByYears)
	uc.Router.HandleFunc("/statistics/get/appointments/by/month/{year}", uc.GetAppointmentsByMonth)
	uc.Router.HandleFunc("/statistics/get/appointments/by/queartals/{year}", uc.GetAppointmentsByQuartals)
	uc.Router.HandleFunc("/statistics/get/appointments/by/years/{startYear}/{endYear}", uc.GetAppointmentsByYears)
	uc.Router.HandleFunc("/statistics/fill/database", uc.FillDataSet)

}

func (uc *StatisticController) GetAverageGrade(res http.ResponseWriter, req *http.Request) {
	authHeader := req.Header.Get("Authorization")
	user, pointer := service.GetUserFromToken(res, req, authHeader)
	if pointer == nil || user.Role != "CompanyAdmin" {
		res.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(res).Encode("Unauthorized")
		return
	}

	fmt.Println(user)

	value, err := service.GetAverageGrade(user.CompmnayID)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error while calculating average grade of company")
		return
	} else {
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(value)
	}

}

func (uc *StatisticController) GetReservationsByMonth(res http.ResponseWriter, req *http.Request) {
	authHeader := req.Header.Get("Authorization")
	user, pointer := service.GetUserFromToken(res, req, authHeader)
	if pointer == nil || user.Role != "CompanyAdmin" {
		res.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(res).Encode("Unauthorized")
		return
	}

	vars := mux.Vars(req)
	year := vars["year"]

	yearInt, err := strconv.Atoi(year)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error with year")
		return
	}

	value, err := service.GetReservationsByMonth(user.CompmnayID, yearInt)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error while calculating monthley reservations")
		return
	} else {
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(value)
	}

}

func (uc *StatisticController) GetReservationsByQuartals(res http.ResponseWriter, req *http.Request) {
	authHeader := req.Header.Get("Authorization")
	user, pointer := service.GetUserFromToken(res, req, authHeader)
	if pointer == nil || user.Role != "CompanyAdmin" {
		res.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(res).Encode("Unauthorized")
		return
	}

	vars := mux.Vars(req)
	year := vars["year"]

	yearInt, err := strconv.Atoi(year)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error with year")
		return
	}

	value, err := service.GetReservationsByQuartals(user.CompmnayID, yearInt)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error while calculating monthley reservations")
		return
	} else {
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(value)
	}

}

func (uc *StatisticController) GetReservationsByYears(res http.ResponseWriter, req *http.Request) {
	authHeader := req.Header.Get("Authorization")
	user, pointer := service.GetUserFromToken(res, req, authHeader)
	if pointer == nil || user.Role != "CompanyAdmin" {
		res.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(res).Encode("Unauthorized")
		return
	}

	vars := mux.Vars(req)
	startYear := vars["startYear"]
	endYear := vars["endYear"]

	startYearInt, err := strconv.Atoi(startYear)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error with start year")
		return
	}
	endYearInt, err := strconv.Atoi(endYear)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error with end year")
		return
	}

	value, err := service.GetReservationsByYears(user.CompmnayID, startYearInt, endYearInt)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error while calculating monthley reservations")
		return
	} else {
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(value)
	}

}

func (uc *StatisticController) GetAppointmentsByMonth(res http.ResponseWriter, req *http.Request) {
	authHeader := req.Header.Get("Authorization")
	user, pointer := service.GetUserFromToken(res, req, authHeader)
	if pointer == nil || user.Role != "CompanyAdmin" {
		res.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(res).Encode("Unauthorized")
		return
	}

	vars := mux.Vars(req)
	year := vars["year"]

	yearInt, err := strconv.Atoi(year)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error with year")
		return
	}

	value, err := service.GetAppointmentsByMonth(user.CompmnayID, yearInt)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error while calculating monthley Appointments")
		return
	} else {
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(value)
	}

}

func (uc *StatisticController) GetAppointmentsByQuartals(res http.ResponseWriter, req *http.Request) {
	authHeader := req.Header.Get("Authorization")
	user, pointer := service.GetUserFromToken(res, req, authHeader)
	if pointer == nil || user.Role != "CompanyAdmin" {
		res.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(res).Encode("Unauthorized")
		return
	}

	vars := mux.Vars(req)
	year := vars["year"]

	yearInt, err := strconv.Atoi(year)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error with year")
		return
	}

	value, err := service.GetAppointmentsByQuartals(user.CompmnayID, yearInt)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error while calculating monthley Appointments")
		return
	} else {
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(value)
	}

}

func (uc *StatisticController) GetAppointmentsByYears(res http.ResponseWriter, req *http.Request) {
	authHeader := req.Header.Get("Authorization")
	user, pointer := service.GetUserFromToken(res, req, authHeader)
	if pointer == nil || user.Role != "CompanyAdmin" {
		res.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(res).Encode("Unauthorized")
		return
	}

	vars := mux.Vars(req)
	startYear := vars["startYear"]
	endYear := vars["endYear"]

	startYearInt, err := strconv.Atoi(startYear)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error with start year")
		return
	}
	endYearInt, err := strconv.Atoi(endYear)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error with end year")
		return
	}

	value, err := service.GetAppointmentsByYears(user.CompmnayID, startYearInt, endYearInt)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error while calculating monthley Appointments")
		return
	} else {
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(value)
	}

}

func (uc *StatisticController) FillDataSet(res http.ResponseWriter, req *http.Request) {
	err := service.FillDataSet()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode("Error while inserting data")
		return
	} else {
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode("Data successfully inserted")
	}
}
