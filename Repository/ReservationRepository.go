package Repository

import (
	"context"
	"fmt"
	"log"
	dto "main/Dto"
	model "main/Model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func SaveReservation(reservation dto.ReservationDto, companyID primitive.ObjectID) (model.Reservation, bool) {
	ReservationsCollection := GetClient().Database("MedicinskaOprema").Collection("Reservations")

	id := primitive.NewObjectID()
	newReservation := model.Reservation{
		ID:        id,
		UserID:    reservation.UserID,
		AssetIDs:  reservation.AssetIDs,
		CompanyID: companyID,
		StartDate: reservation.StartDate,
		EndDate:   reservation.EndDate,
		Canceled:  false,
	}

	insertResult, err := ReservationsCollection.InsertOne(context.Background(), newReservation)
	if err != nil {
		log.Println("Error inserting new reservation:", err)
		var empReservation model.Reservation
		return empReservation, false
	}

	fmt.Println("Added new reservation with ID:", insertResult.InsertedID)

	CompaniesCollection := GetClient().Database("MedicinskaOprema").Collection("Companies")

	filter := bson.M{"_id": companyID}

	update := bson.M{
		"$push": bson.M{
			"reservations": id,
		},
	}

	_, err = CompaniesCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Println("Error updating company with new reservation ID:", err)
		var empReservation model.Reservation
		return empReservation, false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userCollection := GetClient().Database("MedicinskaOprema").Collection("Users")
	_, err = userCollection.UpdateOne(
		ctx,
		bson.M{"_id": newReservation.UserID},
		bson.M{"$inc": bson.M{"points": 1}},
	)
	if err != nil {
		fmt.Println("Error updating user penalties:", err)
		var empReservation model.Reservation
		return empReservation, false
	}

	return newReservation, true
}

func CreateRequest(request dto.AssetsRequests, companyID primitive.ObjectID) []dto.PosibleDate {
	companiesCollection := GetClient().Database("MedicinskaOprema").Collection("Companies")
	reservationsCollection := GetClient().Database("MedicinskaOprema").Collection("Reservations")

	var company model.Company
	companyFilter := bson.M{"_id": companyID}

	err := companiesCollection.FindOne(context.Background(), companyFilter).Decode(&company)
	if err != nil {
		log.Println("Error fetching company:", err)
		return []dto.PosibleDate{}
	}

	reservationFilter := bson.M{
		"assetID": bson.M{"$in": request.AssetIDs},
	}

	cursor, err := reservationsCollection.Find(context.Background(), reservationFilter)
	if err != nil {
		log.Println("Error fetching reservations:", err)
		return []dto.PosibleDate{}
	}
	defer cursor.Close(context.Background())

	var existingReservations []model.Reservation
	if err := cursor.All(context.Background(), &existingReservations); err != nil {
		log.Println("Error decoding reservations:", err)
		return []dto.PosibleDate{}
	}

	var validReservations []model.Reservation
	for _, reservation := range existingReservations {
		if !reservation.Canceled {
			validReservations = append(validReservations, reservation)
		}
	}

	possibleDates := findAvailableDates(validReservations, request.StartDate, request.EndDate, request.NumberOfDays)

	if len(possibleDates) > 0 {
		if len(possibleDates) > 5 {
			possibleDates = possibleDates[:5]
		}
		return possibleDates
	}

	return []dto.PosibleDate{}
}

func findAvailableDates(reservations []model.Reservation, requestStart, requestEnd primitive.DateTime, numberOfDays int64) []dto.PosibleDate {
	start := requestStart.Time()
	end := requestEnd.Time()
	var possibleDates []dto.PosibleDate

	for d := start; d.Before(end); d = d.Add(24 * time.Hour) {
		periodEnd := d.Add(time.Duration(numberOfDays) * 24 * time.Hour)
		if periodEnd.After(end) {
			break
		}
		if isAvailableForPeriod(reservations, d, periodEnd) {
			possibleDates = append(possibleDates, dto.PosibleDate{
				StartDate: primitive.NewDateTimeFromTime(d),
				EndDate:   primitive.NewDateTimeFromTime(periodEnd),
			})
		}
	}

	return possibleDates
}

func isAvailableForPeriod(reservations []model.Reservation, periodStart, periodEnd time.Time) bool {
	for _, reservation := range reservations {
		resStart := reservation.StartDate.Time()
		resEnd := reservation.EndDate.Time()

		if (periodStart.Before(resEnd) && periodEnd.After(resStart)) ||
			(periodStart.Equal(resStart) || periodEnd.Equal(resEnd)) {
			return false
		}
	}
	return true
}

func GetAllForUser(userID primitive.ObjectID) ([]model.Reservation, bool) {
	reservationsCollection := client.Database("MedicinskaOprema").Collection("Reservations")

	filter := bson.M{"userID": userID}

	cursor, err := reservationsCollection.Find(context.Background(), filter)
	if err != nil {
		log.Println("Error fetching reservations:", err)
		return nil, false
	}
	defer cursor.Close(context.Background())

	var reservations []model.Reservation
	if err := cursor.All(context.Background(), &reservations); err != nil {
		log.Println("Error decoding reservations:", err)
		return nil, false
	}

	if len(reservations) == 0 {
		return []model.Reservation{}, true
	}

	return reservations, true
}

func CancelReservation(reservationID primitive.ObjectID, userID primitive.ObjectID) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	reservationCollection := GetClient().Database("MedicinskaOprema").Collection("Reservations")
	userCollection := GetClient().Database("MedicinskaOprema").Collection("Users")

	var reservation bson.M
	err := reservationCollection.FindOne(ctx, bson.M{"_id": reservationID}).Decode(&reservation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("Reservation not found")
		} else {
			fmt.Println("Error finding reservation:", err)
		}
		return false
	}

	_, err = reservationCollection.UpdateOne(ctx, bson.M{"_id": reservationID}, bson.M{"$set": bson.M{"canceled": true, "userID": primitive.NilObjectID}})
	if err != nil {
		fmt.Println("Error updating reservation:", err)
		return false
	}

	endDate, ok := reservation["endDate"].(primitive.DateTime)
	if !ok {
		fmt.Println("Invalid end date format in reservation")
		return false
	}

	timeDiff := time.Until(endDate.Time())

	penaltyIncrement := 1
	if timeDiff.Hours() < 24 {
		penaltyIncrement = 2
	}

	_, err = userCollection.UpdateOne(
		ctx,
		bson.M{"_id": userID},
		bson.M{"$inc": bson.M{"penalties": penaltyIncrement}},
	)
	if err != nil {
		fmt.Println("Error updating user penalties:", err)
		return false
	}

	return true
}

func GetAllCanceled() ([]model.Reservation, bool) {
	reservationsCollection := client.Database("MedicinskaOprema").Collection("Reservations")

	filter := bson.M{"canceled": true}

	cursor, err := reservationsCollection.Find(context.Background(), filter)
	if err != nil {
		log.Println("Error fetching reservations:", err)
		return nil, false
	}
	defer cursor.Close(context.Background())

	var reservations []model.Reservation
	if err := cursor.All(context.Background(), &reservations); err != nil {
		log.Println("Error decoding reservations:", err)
		return nil, false
	}

	if len(reservations) == 0 {
		return []model.Reservation{}, true
	}

	return reservations, true
}

func TakeCanceledReservation(reservationID primitive.ObjectID, userID primitive.ObjectID) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	reservationCollection := GetClient().Database("MedicinskaOprema").Collection("Reservations")

	var reservation bson.M
	err := reservationCollection.FindOne(ctx, bson.M{"_id": reservationID}).Decode(&reservation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("Reservation not found")
		} else {
			fmt.Println("Error finding reservation:", err)
		}
		return false
	}

	canceled, ok := reservation["canceled"].(bool)
	if !ok || !canceled {
		fmt.Println("Reservation is not canceled or invalid format")
		return false
	}

	_, err = reservationCollection.UpdateOne(
		ctx,
		bson.M{"_id": reservationID},
		bson.M{
			"$set": bson.M{
				"canceled": false,
				"userID":   userID,
			},
		},
	)
	if err != nil {
		fmt.Println("Error updating reservation:", err)
		return false
	}

	return true
}
