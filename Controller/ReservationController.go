package Controller

import (
	"encoding/json"
	dto "main/Dto"
	repository "main/Repository"
	service "main/Service"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReservationController struct {
	Router *mux.Router
}

func NewReservationController() *ReservationController {
	return &ReservationController{
		Router: mux.NewRouter(),
	}
}

func (uc *ReservationController) RegisterRoutes() {
	uc.Router.HandleFunc("/reservations/request", uc.RequestForReservation)
	uc.Router.HandleFunc("/reservations/create", uc.CreateReservation)
	uc.Router.HandleFunc("/reservations/get/all/for/{userID}", uc.GetAllForUser)
	uc.Router.HandleFunc("/reservations/cancel/{reservationID}", uc.Cancel)
	uc.Router.HandleFunc("/reservations/get/all/canceled", uc.GetAllCanceled)
	uc.Router.HandleFunc("/reservations/take/{reservationID}", uc.TakeCanceled)

}

func (uc *ReservationController) RequestForReservation(res http.ResponseWriter, req *http.Request) {
	companyID := req.Header.Get("companyID")
	companyObjectID, companyIDErr := primitive.ObjectIDFromHex(companyID)
	if companyIDErr != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error with company id")
		return
	}
	var request dto.AssetsRequests
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		http.Error(res, "Error decoding request body", http.StatusBadRequest)
		return
	}

	dates := service.CreateRequest(request, companyObjectID, res)

	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(dates)

}

func (uc *ReservationController) CreateReservation(res http.ResponseWriter, req *http.Request) {
	var requestBody dto.ReservationDto
	err := json.NewDecoder(req.Body).Decode(&requestBody)
	if err != nil {
		http.Error(res, "Error decoding request body", http.StatusBadRequest)
		return
	}

	companyID := req.Header.Get("companyID")
	companyObjectID, companyIDErr := primitive.ObjectIDFromHex(companyID)
	if companyIDErr != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error with company id")
		return
	}

	reservation, done := service.SaveReservation(requestBody, companyObjectID, res)
	if done {
		filePath := reservation.ID.Hex() + ".png"
		err = service.GenerateQRCode(reservation, filePath)
		if err != nil {
			http.Error(res, "Error generating QR code", http.StatusInternalServerError)
			return
		}

		user, _ := repository.FindUserById(reservation.UserID, res)
		service.SendReservationMail(user.Email, filePath)

		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode("Reservation successfully created with QR code")
	} else {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error while creating reservation")
	}

}

func (uc *ReservationController) GetAllForUser(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	userID := vars["userID"]

	userObjectID, userIDErr := primitive.ObjectIDFromHex(userID)
	if userIDErr != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error with user id")
		return
	}

	requsts, done := service.GetAllForUser(userObjectID, res)
	if !done {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error while find requests")
		return
	} else {
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(requsts)
	}

}

func (uc *ReservationController) Cancel(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	reservationID := vars["reservationID"]

	reservationObjectID, reservationIDErr := primitive.ObjectIDFromHex(reservationID)
	if reservationIDErr != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error with reservation id")
		return
	}
	authHeader := req.Header.Get("Authorization")
	user, pointer := service.GetUserFromToken(res, req, authHeader)
	if pointer == nil || user.Role != "User" {
		res.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(res).Encode("Unauthorized")
		return
	}
	done := service.CancelReservation(reservationObjectID, user.ID)
	if done {
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode("Reservation deleted successfuly")
	} else {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error while deleting")
	}
}

func (uc *ReservationController) GetAllCanceled(res http.ResponseWriter, req *http.Request) {
	requsts, done := service.GetAllCanceled(res)
	if !done {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error while find requests")
		return
	} else {
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(requsts)
	}

}

func (uc *ReservationController) TakeCanceled(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	reservationID := vars["reservationID"]

	reservationObjectID, reservationIDErr := primitive.ObjectIDFromHex(reservationID)
	if reservationIDErr != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error with reservation id")
		return
	}
	authHeader := req.Header.Get("Authorization")
	user, pointer := service.GetUserFromToken(res, req, authHeader)
	if pointer == nil || user.Role != "User" {
		res.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(res).Encode("Unauthorized")
		return
	}
	done := service.TakeCanceledReservation(reservationObjectID, user.ID)
	if done {
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode("Reservation taken successfully")
	} else {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode("Error while taking")
	}

}
