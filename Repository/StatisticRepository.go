package Repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAverageGrade(companyID primitive.ObjectID) (float64, error) {
	gradesCollection := GetClient().Database("MedicinskaOprema").Collection("Grades")

	filter := bson.M{"companyid": companyID}

	cursor, err := gradesCollection.Find(context.TODO(), filter)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(context.TODO())

	var totalGrade float64
	var count int

	for cursor.Next(context.TODO()) {
		var gradeDoc struct {
			Grade int `bson:"grade"`
		}
		if err := cursor.Decode(&gradeDoc); err != nil {
			return 0, err
		}

		totalGrade += float64(gradeDoc.Grade)
		count++
	}

	if count == 0 {
		return 0, nil
	}

	averageGrade := totalGrade / float64(count)
	return averageGrade, nil
}

func GetReservationsByMonth(companyID primitive.ObjectID, year int) (map[string]int, error) {
	reservationsCollection := GetClient().Database("MedicinskaOprema").Collection("Reservations")

	filter := bson.M{
		"companyID": companyID,
		"startDate": bson.M{
			"$gte": time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC),
			"$lt":  time.Date(year+1, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	cursor, err := reservationsCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	reservationsByMonth := map[string]int{
		"January":   0,
		"February":  0,
		"March":     0,
		"April":     0,
		"May":       0,
		"June":      0,
		"July":      0,
		"August":    0,
		"September": 0,
		"October":   0,
		"November":  0,
		"December":  0,
	}

	for cursor.Next(context.TODO()) {
		var reservation struct {
			StartDate time.Time `bson:"startDate"`
		}

		if err := cursor.Decode(&reservation); err != nil {
			return nil, err
		}

		month := reservation.StartDate.Month()

		monthName := month.String()

		reservationsByMonth[monthName]++
	}

	return reservationsByMonth, nil
}

func GetReservationsByQuartals(companyID primitive.ObjectID, year int) (map[string]int, error) {
	reservationsCollection := GetClient().Database("MedicinskaOprema").Collection("Reservations")

	filter := bson.M{
		"companyID": companyID,
		"startDate": bson.M{
			"$gte": time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC),
			"$lt":  time.Date(year+1, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	cursor, err := reservationsCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	reservationsByQuartals := map[string]int{
		"Q1": 0, // Januar, Februar, Mart
		"Q2": 0, // April, Maj, Jun
		"Q3": 0, // Jul, Avgust, Septembar
		"Q4": 0, // Oktobar, Novembar, Decembar
	}

	for cursor.Next(context.TODO()) {
		var reservation struct {
			StartDate time.Time `bson:"startDate"`
		}

		if err := cursor.Decode(&reservation); err != nil {
			return nil, err
		}

		month := reservation.StartDate.Month()

		switch {
		case month >= time.January && month <= time.March:
			reservationsByQuartals["Q1"]++
		case month >= time.April && month <= time.June:
			reservationsByQuartals["Q2"]++
		case month >= time.July && month <= time.September:
			reservationsByQuartals["Q3"]++
		case month >= time.October && month <= time.December:
			reservationsByQuartals["Q4"]++
		}
	}

	return reservationsByQuartals, nil
}

func GetReservationsByYears(companyID primitive.ObjectID, startYear int, endYear int) (map[string]int, error) {
	reservationsCollection := GetClient().Database("MedicinskaOprema").Collection("Reservations")

	filter := bson.M{
		"companyID": companyID,
		"startDate": bson.M{
			"$gte": time.Date(startYear, time.January, 1, 0, 0, 0, 0, time.UTC),
			"$lt":  time.Date(endYear+1, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	cursor, err := reservationsCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	reservationsByYear := make(map[string]int)
	for year := startYear; year <= endYear; year++ {
		reservationsByYear[fmt.Sprintf("%d", year)] = 0
	}

	for cursor.Next(context.TODO()) {
		var reservation struct {
			StartDate time.Time `bson:"startDate"`
		}

		if err := cursor.Decode(&reservation); err != nil {
			return nil, err
		}

		year := reservation.StartDate.Year()

		yearStr := fmt.Sprintf("%d", year)
		if _, exists := reservationsByYear[yearStr]; exists {
			reservationsByYear[yearStr]++
		}
	}

	return reservationsByYear, nil
}

func GetAppointmentsByMonth(companyID primitive.ObjectID, year int) (map[string]int, error) {
	appointmentsCollection := GetClient().Database("MedicinskaOprema").Collection("Appointments")

	filter := bson.M{
		"companyID": companyID,
		"startDate": bson.M{
			"$gte": time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC),
			"$lt":  time.Date(year+1, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	cursor, err := appointmentsCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	appointmentsByMonth := map[string]int{
		"January":   0,
		"February":  0,
		"March":     0,
		"April":     0,
		"May":       0,
		"June":      0,
		"July":      0,
		"August":    0,
		"September": 0,
		"October":   0,
		"November":  0,
		"December":  0,
	}

	for cursor.Next(context.TODO()) {
		var appointment struct {
			StartDate time.Time `bson:"startDate"`
		}

		if err := cursor.Decode(&appointment); err != nil {
			return nil, err
		}

		month := appointment.StartDate.Month()

		monthName := month.String()

		appointmentsByMonth[monthName]++
	}

	return appointmentsByMonth, nil
}

func GetAppointmentsByQuartals(companyID primitive.ObjectID, year int) (map[string]int, error) {
	appointmentsCollection := GetClient().Database("MedicinskaOprema").Collection("Appointments")

	filter := bson.M{
		"companyID": companyID,
		"startDate": bson.M{
			"$gte": time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC),
			"$lt":  time.Date(year+1, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	cursor, err := appointmentsCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	appointmentsByQuartals := map[string]int{
		"Q1": 0, // Januar, Februar, Mart
		"Q2": 0, // April, Maj, Jun
		"Q3": 0, // Jul, Avgust, Septembar
		"Q4": 0, // Oktobar, Novembar, Decembar
	}

	for cursor.Next(context.TODO()) {
		var appointment struct {
			StartDate time.Time `bson:"startDate"`
		}

		if err := cursor.Decode(&appointment); err != nil {
			return nil, err
		}

		month := appointment.StartDate.Month()

		switch {
		case month >= time.January && month <= time.March:
			appointmentsByQuartals["Q1"]++
		case month >= time.April && month <= time.June:
			appointmentsByQuartals["Q2"]++
		case month >= time.July && month <= time.September:
			appointmentsByQuartals["Q3"]++
		case month >= time.October && month <= time.December:
			appointmentsByQuartals["Q4"]++
		}
	}

	return appointmentsByQuartals, nil
}

func GetAppointmentsByYears(companyID primitive.ObjectID, startYear int, endYear int) (map[string]int, error) {
	appointmentsCollection := GetClient().Database("MedicinskaOprema").Collection("Appointments")

	filter := bson.M{
		"companyID": companyID,
		"startDate": bson.M{
			"$gte": time.Date(startYear, time.January, 1, 0, 0, 0, 0, time.UTC),
			"$lt":  time.Date(endYear+1, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	cursor, err := appointmentsCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	appointmentsByYear := make(map[string]int)
	for year := startYear; year <= endYear; year++ {
		appointmentsByYear[fmt.Sprintf("%d", year)] = 0
	}

	for cursor.Next(context.TODO()) {
		var appointment struct {
			StartDate time.Time `bson:"startDate"`
		}

		if err := cursor.Decode(&appointment); err != nil {
			return nil, err
		}

		year := appointment.StartDate.Year()

		yearStr := fmt.Sprintf("%d", year)
		if _, exists := appointmentsByYear[yearStr]; exists {
			appointmentsByYear[yearStr]++
		}
	}

	return appointmentsByYear, nil
}
