package Service

import (
	repository "main/Repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAverageGrade(companyID primitive.ObjectID) (float64, error) {
	return repository.GetAverageGrade(companyID)
}

func GetReservationsByMonth(companyID primitive.ObjectID, year int) (map[string]int, error) {
	return repository.GetReservationsByMonth(companyID, year)
}

func GetReservationsByQuartals(companyID primitive.ObjectID, year int) (map[string]int, error) {
	return repository.GetReservationsByQuartals(companyID, year)
}

func GetReservationsByYears(companyID primitive.ObjectID, startYear int, endYear int) (map[string]int, error) {
	return repository.GetReservationsByYears(companyID, startYear, endYear)
}

func GetAppointmentsByMonth(companyID primitive.ObjectID, year int) (map[string]int, error) {
	return repository.GetAppointmentsByMonth(companyID, year)
}

func GetAppointmentsByQuartals(companyID primitive.ObjectID, year int) (map[string]int, error) {
	return repository.GetAppointmentsByQuartals(companyID, year)
}

func GetAppointmentsByYears(companyID primitive.ObjectID, startYear int, endYear int) (map[string]int, error) {
	return repository.GetAppointmentsByYears(companyID, startYear, endYear)
}

func FillDataSet() error {
	return repository.FillDataSet()
}
