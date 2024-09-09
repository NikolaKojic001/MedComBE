package Repository

import (
	"context"
	"fmt"
	model "main/Model"
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

type Appointment struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	AssetID   primitive.ObjectID `bson:"assetID"`
	UserID    primitive.ObjectID `bson:"userID"`
	CompanyID primitive.ObjectID `bson:"companyID"`
	StartDate time.Time          `bson:"startDate"`
}

func StringToObject(id string) primitive.ObjectID {
	ID, _ := primitive.ObjectIDFromHex(id)
	return ID
}

func FillDataSet() error {

	db := GetClient().Database("Proba")

	commands := []bson.D{}
	collectionNames := []string{}
	commands = append(commands, bson.D{{Key: "create", Value: "Appointments"}})
	collectionNames = append(collectionNames, "Appointments")

	commands = append(commands, bson.D{{Key: "create", Value: "Companies"}})
	collectionNames = append(collectionNames, "Companies")

	commands = append(commands, bson.D{{Key: "create", Value: "Grades"}})
	collectionNames = append(collectionNames, "Grades")

	commands = append(commands, bson.D{{Key: "create", Value: "Users"}})
	collectionNames = append(collectionNames, "Users")

	commands = append(commands, bson.D{{Key: "create", Value: "LoyaltyPrograms"}})
	collectionNames = append(collectionNames, "LoyaltyPrograms")

	commands = append(commands, bson.D{{Key: "create", Value: "Reports"}})
	collectionNames = append(collectionNames, "Reports")

	commands = append(commands, bson.D{{Key: "create", Value: "Reservations"}})
	collectionNames = append(collectionNames, "Reservations")

	var collectionResult bson.M
	for i, command := range commands {
		if collectionErr := db.RunCommand(context.TODO(), command).Decode(&collectionResult); collectionErr != nil {
			for j := 0; j < i; j++ {
				dropCmd := bson.D{{Key: "drop", Value: collectionNames[j]}}
				db.RunCommand(context.TODO(), dropCmd)
			}
			return collectionErr
		}
	}

	appointmentsCollection := db.Collection("Appointments")
	companiesCollection := db.Collection("Companies")
	gradesCollection := db.Collection("Grades")
	loyaltyProgramsCollection := db.Collection("LoyaltyPrograms")
	reportsCollection := db.Collection("Reports")
	reservationsCollection := db.Collection("Reservations")
	usersCollection := db.Collection("Users")

	appointments := []interface{}{
		Appointment{
			ID:        StringToObject("64a7f8d3e138f1e5b7d6f118"),
			AssetID:   StringToObject("66c9d4c1dee03ccc5d094440"),
			UserID:    StringToObject("64b8c68f1350dff265912861"),
			CompanyID: StringToObject("66c9c5f4a78345b32a3d8ad2"),
			StartDate: time.Date(2024, 7, 17, 17, 0, 0, 0, time.UTC),
		},
		Appointment{
			ID:        StringToObject("64a7f8d3e138f1e5b7d6f115"),
			AssetID:   StringToObject("66c9d4c1dee03ccc5d094437"),
			UserID:    StringToObject("64b8c68f1350dff26591285e"),
			CompanyID: StringToObject("66c9c5f4a78345b32a3d8ad2"),
			StartDate: time.Date(2024, 3, 14, 15, 0, 0, 0, time.UTC),
		},
		Appointment{
			ID:        StringToObject("64a7f8d3e138f1e5b7d6f113"),
			AssetID:   StringToObject("66c9d4c1dee03ccc5d094435"),
			UserID:    StringToObject("64b8c68f1350dff26591285c"),
			CompanyID: StringToObject("66c9c5f4a78345b32a3d8ad2"),
			StartDate: time.Date(2024, 1, 12, 9, 0, 0, 0, time.UTC),
		},
		Appointment{
			ID:        StringToObject("64a7f8d3e138f1e5b7d6f120"),
			AssetID:   StringToObject("66c9d4c1dee03ccc5d094442"),
			UserID:    StringToObject("64b8c68f1350dff265912863"),
			CompanyID: StringToObject("66c9c5f4a78345b32a3d8ad2"),
			StartDate: time.Date(2024, 9, 19, 10, 0, 0, 0, time.UTC),
		},
		Appointment{
			ID:        StringToObject("64a7f8d3e138f1e5b7d6f117"),
			AssetID:   StringToObject("66c9d4c1dee03ccc5d094439"),
			UserID:    StringToObject("64b8c68f1350dff265912860"),
			CompanyID: StringToObject("66c9c5f4a78345b32a3d8ad2"),
			StartDate: time.Date(2024, 9, 16, 8, 0, 0, 0, time.UTC),
		},
		Appointment{
			ID:        StringToObject("64a7f8d3e138f1e5b7d6f114"),
			AssetID:   StringToObject("66c9d4c1dee03ccc5d094436"),
			UserID:    StringToObject("64b8c68f1350dff26591285d"),
			CompanyID: StringToObject("66c9c5f4a78345b32a3d8ad2"),
			StartDate: time.Date(2024, 3, 13, 13, 0, 0, 0, time.UTC),
		},
		Appointment{
			ID:        StringToObject("64a7f8d3e138f1e5b7d6f111"),
			AssetID:   StringToObject("66c9d4c1dee03ccc5d094433"),
			UserID:    StringToObject("64b8c68f1350dff26591285a"),
			CompanyID: StringToObject("66c9c5f4a78345b32a3d8ad2"),
			StartDate: time.Date(2023, 9, 10, 10, 0, 0, 0, time.UTC),
		},
		Appointment{
			ID:        StringToObject("64a7f8d3e138f1e5b7d6f119"),
			AssetID:   StringToObject("66c9d4c1dee03ccc5d094441"),
			UserID:    StringToObject("64b8c68f1350dff265912862"),
			CompanyID: StringToObject("66c9c5f4a78345b32a3d8ad2"),
			StartDate: time.Date(2024, 9, 18, 12, 0, 0, 0, time.UTC),
		},
		Appointment{
			ID:        StringToObject("64a7f8d3e138f1e5b7d6f112"),
			AssetID:   StringToObject("66c9d4c1dee03ccc5d094434"),
			UserID:    StringToObject("64b8c68f1350dff26591285b"),
			CompanyID: StringToObject("66c9c5f4a78345b32a3d8ad2"),
			StartDate: time.Date(2023, 9, 11, 14, 0, 0, 0, time.UTC),
		},
		Appointment{
			ID:        StringToObject("64a7f8d3e138f1e5b7d6f121"),
			AssetID:   StringToObject("66c9d4c1dee03ccc5d094443"),
			UserID:    StringToObject("64b8c68f1350dff265912864"),
			CompanyID: StringToObject("66c9c5f4a78345b32a3d8ad2"),
			StartDate: time.Date(2024, 9, 20, 9, 0, 0, 0, time.UTC),
		},
		Appointment{
			ID:        StringToObject("64a7f8d3e138f1e5b7d6f116"),
			AssetID:   StringToObject("66c9d4c1dee03ccc5d094438"),
			UserID:    StringToObject("64b8c68f1350dff26591285f"),
			CompanyID: StringToObject("66c9c5f4a78345b32a3d8ad2"),
			StartDate: time.Date(2024, 9, 15, 11, 0, 0, 0, time.UTC),
		},
	}

	_, err = appointmentsCollection.InsertMany(context.TODO(), appointments)
	if err != nil {
		fmt.Println("Error inserting appointments:", err)
		return err
	}

	fmt.Println("Appointments inserted successfully!")

	companies := []interface{}{
		model.Company{
			ID:       StringToObject("66c9c5f4a78345b32a3d8ad2"),
			AdminID:  StringToObject("66c9c581de1c903f3d1b0647"),
			Name:     "KojicPlast",
			Location: model.Location{City: "Novi Sad", Country: "Serbia"},
			Assets: []model.Asset{
				{ID: StringToObject("66c9d4c1dee03ccc5d094433"), Name: "toplomer"},
				{ID: StringToObject("66c9d53d2d1b1814e2876824"), Name: "stetoskop"},
				{ID: StringToObject("66c9d8cc86b559875213bd2e"), Name: "lenjir"},
				{ID: StringToObject("66c9d8d186b559875213bd2f"), Name: "makaze"},
			},
			Reservations: []primitive.ObjectID{
				StringToObject("66d75d2f8ff1a013e4bf3d07"),
				StringToObject("66d76609958d22bb1bc886e9"),
				StringToObject("66d9d3322e66149f425af333"),
			},
		},
	}

	_, err = companiesCollection.InsertMany(context.TODO(), companies)
	if err != nil {
		fmt.Println("Error inserting companies:", err)
		return err
	}

	fmt.Println("Companies inserted successfully!")

	grades := []interface{}{
		model.Grade{
			ID:          StringToObject("66d9d64b2e66149f425af334"),
			Grade:       4,
			UserID:      StringToObject("66c8c68f1250bdf26591285a"),
			CompanyID:   StringToObject("66c9c5f4a78345b32a3d8ad2"),
			Description: []string{"cistoca", "ljubaznost"},
		},
		model.Grade{
			ID:          StringToObject("66db22f829635607ede47c89"),
			Grade:       2,
			UserID:      StringToObject("66c8c68f1250bdf26591285a"),
			CompanyID:   StringToObject("66c9c5f4a78345b32a3d8ad2"),
			Description: []string{"cistoca", "ljubaznost"},
		},
	}

	// Ubaci dokumente u kolekciju
	_, err = gradesCollection.InsertMany(context.TODO(), grades)
	if err != nil {
		fmt.Println("Error inserting grades:", err)
		return err
	}

	fmt.Println("Grades inserted successfully!")

	loyalties := []interface{}{
		model.LoyaltyProgram{
			ID:        StringToObject("66d9c77a2e66149f425af330"),
			Name:      "Bronze",
			Discount:  5,
			UpValue:   10,
			DownValue: 5,
		},
		model.LoyaltyProgram{
			ID:        StringToObject("66d9c78e2e66149f425af331"),
			Name:      "Silver",
			Discount:  10,
			UpValue:   15,
			DownValue: 10,
		},
		model.LoyaltyProgram{
			ID:        StringToObject("66d9c7ab2e66149f425af332"),
			Name:      "Gold",
			Discount:  20,
			UpValue:   100,
			DownValue: 15,
		},
	}

	// Ubaci dokumente u kolekciju
	_, err = loyaltyProgramsCollection.InsertMany(context.TODO(), loyalties)
	if err != nil {
		fmt.Println("Error inserting loyalty programs:", err)
		return err
	}

	fmt.Println("Loyalty programs inserted successfully!")

	reports := []interface{}{
		model.Report{
			ID:          StringToObject("66cdb9f8ea33165e4ec0f901"),
			UserID:      StringToObject("66c8c68f1250bdf26591285a"),
			CompanyID:   StringToObject("66c9c5f4a78345b32a3d8ad2"),
			Description: "Oprema nije u novom stanju",
			Replay:      "Izvinite odmah cemo poslati ispravnu opremu",
		},
		model.Report{
			ID:          StringToObject("66d89f167f0854626eaade83"),
			UserID:      StringToObject("66c8c68f1250bdf26591285a"),
			CompanyID:   StringToObject("66c9c5f4a78345b32a3d8ad2"),
			Description: "Nista ne valja",
			Replay:      "bas me briga",
		},
		model.Report{
			ID:          StringToObject("66d8b20a5dbece87c3ed7c12"),
			UserID:      StringToObject("66c8c68f1250bdf26591285a"),
			CompanyID:   StringToObject("66c9c5f4a78345b32a3d8ad2"),
			Description: "Nista ne valja",
			Replay:      "",
		},
		model.Report{
			ID:          StringToObject("66d8b20c5dbece87c3ed7c13"),
			UserID:      StringToObject("66c8c68f1250bdf26591285a"),
			CompanyID:   StringToObject("66c9c5f4a78345b32a3d8ad2"),
			Description: "Nista ne valja",
			Replay:      "",
		},
	}

	// Ubaci dokumente u kolekciju
	_, err = reportsCollection.InsertMany(context.TODO(), reports)
	if err != nil {
		fmt.Println("Error inserting reports:", err)
		return err
	}

	fmt.Println("Reports inserted successfully!")

	reservations := []interface{}{
		model.Reservation{
			ID:        StringToObject("66d75d2f8ff1a013e4bf3d07"),
			AssetIDs:  []primitive.ObjectID{StringToObject("66c9d4c1dee03ccc5d094433")},
			UserID:    StringToObject("000000000000000000000000"),
			CompanyID: StringToObject("66c9c5f4a78345b32a3d8ad2"),
			StartDate: primitive.NewDateTimeFromTime(time.Date(2024, 9, 3, 22, 0, 1, 1, time.UTC)),
			EndDate:   primitive.NewDateTimeFromTime(time.Date(2024, 9, 5, 22, 0, 1, 1, time.UTC)),
			Canceled:  true,
		},
		model.Reservation{
			ID:        StringToObject("66d76609958d22bb1bc886e9"),
			AssetIDs:  []primitive.ObjectID{StringToObject("66c9d4c1dee03ccc5d094433")},
			UserID:    StringToObject("66c8c68f1250bdf26591285a"),
			CompanyID: StringToObject("66c9c5f4a78345b32a3d8ad2"),
			StartDate: primitive.NewDateTimeFromTime(time.Date(2024, 9, 5, 22, 0, 1, 1, time.UTC)),
			EndDate:   primitive.NewDateTimeFromTime(time.Date(2024, 9, 9, 22, 0, 1, 1, time.UTC)),
			Canceled:  false,
		},
		model.Reservation{
			ID:        StringToObject("66d9d3322e66149f425af333"),
			AssetIDs:  []primitive.ObjectID{StringToObject("66c9d4c1dee03ccc5d094433")},
			UserID:    StringToObject("66c8c68f1250bdf26591285a"),
			CompanyID: StringToObject("66c9c5f4a78345b32a3d8ad2"),
			StartDate: primitive.NewDateTimeFromTime(time.Date(2023, 9, 12, 22, 0, 1, 1, time.UTC)),
			EndDate:   primitive.NewDateTimeFromTime(time.Date(2023, 9, 16, 22, 0, 1, 1, time.UTC)),
			Canceled:  false,
		},
	}

	// Ubaci dokumente u kolekciju
	_, err = reservationsCollection.InsertMany(context.TODO(), reservations)
	if err != nil {
		fmt.Println("Error inserting reservations:", err)
		return err
	}

	fmt.Println("Reservations inserted successfully!")

	users := []interface{}{
		model.User{
			ID:             StringToObject("66c8c68f1250bdf26591285a"),
			Email:          "kojicn7@gmail.com",
			Password:       "$2a$10$VfC2ySe6aUxwzFCXrkEAXOWbKcODu4fkYPMYTM0n8jAqrGQ8ATOhy",
			Role:           "User",
			Firstname:      "John",
			Lastname:       "Doe",
			Location:       model.Location{City: "Belgrade", Country: "Serbia"},
			Phonenumber:    "+381601234567",
			Profession:     "Software Developer",
			CompmnayID:     primitive.NilObjectID,
			Verified:       true,
			Penalies:       0,
			Points:         6,
			LoyaltyProgram: StringToObject("66d9c77a2e66149f425af330"),
		},
		model.User{
			ID:          StringToObject("66c9c581de1c903f3d1b0647"),
			Email:       "kojicn116@gmail.com",
			Password:    "$2a$10$MGvMUuXDfi7c4JldQoMk..uR46NIvutZTo.nArGd.mBihkLrLHU4K",
			Role:        "CompanyAdmin",
			Firstname:   "John",
			Lastname:    "Doe",
			Location:    model.Location{City: "Belgrade", Country: "Serbia"},
			Phonenumber: "+381601234567",
			Profession:  "Software Developer",
			CompmnayID:  StringToObject("66c9c5f4a78345b32a3d8ad2"),
			Penalies:    0,
			Points:      0,
			Verified:    true,
		},
		model.User{
			ID:          StringToObject("66cc891359d3fa8ba9efbda7"),
			Email:       "kojicn1@gmail.com",
			Password:    "$2a$10$hcKq3bl8EAtgceyaVYXvBeP7XoJeo1uCu7/Z23.6B2ZrQd9OQqQC.",
			Role:        "Admin",
			Firstname:   "John",
			Lastname:    "Doe",
			Location:    model.Location{City: "Belgrade", Country: "Serbia"},
			Phonenumber: "+381601234567",
			Profession:  "Software Developer",
			Verified:    true,
			Penalies:    0,
			Points:      0,
			CompmnayID:  primitive.NilObjectID,
		},
	}

	// Ubaci dokumente u kolekciju
	_, err = usersCollection.InsertMany(context.TODO(), users)
	if err != nil {
		fmt.Println("Error inserting users:", err)
		return err
	}

	fmt.Println("Users inserted successfully!")

	return nil
}
