package Service

import (
	"encoding/json"
	dto "main/Dto"
	model "main/Model"
	repository "main/Repository"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SaveReservation(reservation dto.ReservationDto, companyID primitive.ObjectID, res http.ResponseWriter) (model.Reservation, bool) {
	company, found := repository.FindCompanyById(companyID, res)
	if !found {
		res.WriteHeader(http.StatusNotFound)
		json.NewEncoder(res).Encode("Company not found")
	}
	return repository.SaveReservation(reservation, company.ID)
}

func CreateRequest(request dto.AssetsRequests, companyID primitive.ObjectID, res http.ResponseWriter) []dto.PosibleDate {
	company, found := repository.FindCompanyById(companyID, res)
	if !found {
		res.WriteHeader(http.StatusNotFound)
		json.NewEncoder(res).Encode("Company not found")
	}
	return repository.CreateRequest(request, company.ID)
}

func GetAllForUser(userID primitive.ObjectID, res http.ResponseWriter) ([]model.Reservation, bool) {
	user, found := repository.FindUserById(userID, res)
	if !found {
		res.WriteHeader(http.StatusNotFound)
		json.NewEncoder(res).Encode("User not found")
	}

	return repository.GetAllForUser(user.ID)

}

func CancelReservation(reservationID primitive.ObjectID, userID primitive.ObjectID) bool {
	return repository.CancelReservation(reservationID, userID)
}

func GetAllCanceled(res http.ResponseWriter) ([]model.Reservation, bool) {
	return repository.GetAllCanceled()
}

func TakeCanceledReservation(reservationID primitive.ObjectID, userID primitive.ObjectID) bool {
	return repository.TakeCanceledReservation(reservationID, userID)
}
