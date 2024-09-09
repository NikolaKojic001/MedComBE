package Repository

import (
	"context"
	"encoding/json"
	"fmt"
	dto "main/Dto"
	model "main/Model"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func SaveCompany(company dto.CompanyDto) (bool, primitive.ObjectID) {
	UsersCollection := GetClient().Database("MedicinskaOprema").Collection("Companies")
	location := model.Location{
		City:    company.LocationCompany.City,
		Country: company.LocationCompany.Country,
	}
	id := primitive.NewObjectID()
	newCompany := model.Company{
		ID:           id,
		Name:         company.Name,
		Location:     location,
		Assets:       make([]model.Asset, 0),
		Reservations: make([]primitive.ObjectID, 0),
	}

	insertResult, err := UsersCollection.InsertOne(context.Background(), newCompany)
	if err != nil {
		return false, primitive.NilObjectID
	}

	fmt.Println("Added new company with ID:", insertResult.InsertedID)
	return true, (insertResult.InsertedID).(primitive.ObjectID)

}

func FindCompanyByName(companyName string) (model.Company, bool) {
	UsersCollection := GetClient().Database("MedicinskaOprema").Collection("Companies")
	filter := bson.M{"name": companyName}
	var result model.Company
	err = UsersCollection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			var company model.Company
			return company, false
		}
		fmt.Println(err)
	}
	return result, true
}

func GetAllCompanies() ([]model.Company, error) {
	CompaniesCollection := GetClient().Database("MedicinskaOprema").Collection("Companies")

	// Use context for better control over resource management
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Fetch all documents without projection
	cursor, err := CompaniesCollection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []model.Company
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func FindCompanyById(id primitive.ObjectID, res http.ResponseWriter) (model.Company, bool) {
	collection := GetClient().Database("MedicinskaOprema").Collection("Companies")
	filter := bson.M{"_id": id}
	var result model.Company
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			json.NewEncoder(res).Encode("Didnt find company!")
			return result, false
		}
	}
	return result, true
}

func ReportCompany(userID primitive.ObjectID, companyID primitive.ObjectID, requestBody dto.ReportDto) bool {
	reservationsCollection := GetClient().Database("MedicinskaOprema").Collection("Reservations")
	reportsCollection := GetClient().Database("MedicinskaOprema").Collection("Reports")

	filter := bson.M{
		"userID":    userID,
		"companyID": companyID,
		"canceled":  false,
	}
	count, err := reservationsCollection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return false
	}

	if count > 0 {
		report := model.Report{
			ID:          primitive.NewObjectID(),
			UserID:      userID,
			CompanyID:   companyID,
			Description: requestBody.Description,
			Replay:      "",
		}

		_, err := reportsCollection.InsertOne(context.TODO(), report)
		return err == nil
	} else {
		return false
	}

}

func GradeCompany(userID primitive.ObjectID, companyID primitive.ObjectID, requestBody dto.GradeDto) bool {
	reservationsCollection := GetClient().Database("MedicinskaOprema").Collection("Reservations")
	gradesCollection := GetClient().Database("MedicinskaOprema").Collection("Grades")

	now := time.Now()

	filter := bson.M{
		"userID":    userID,
		"companyID": companyID,
		"endDate":   bson.M{"$lt": now},
		"canceled":  false,
	}
	count, err := reservationsCollection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return false
	}

	if requestBody.Grade < 1 || requestBody.Grade > 10 {
		return false
	}

	if count > 0 {
		grade := model.Grade{
			ID:          primitive.NewObjectID(),
			UserID:      userID,
			CompanyID:   companyID,
			Description: requestBody.Description,
			Grade:       requestBody.Grade,
		}

		_, err := gradesCollection.InsertOne(context.TODO(), grade)
		return err == nil
	} else {
		return false
	}
}

func GetAllCompanyReports(companyID primitive.ObjectID, res http.ResponseWriter) ([]model.Report, error) {
	reportsCollection := GetClient().Database("MedicinskaOprema").Collection("Reports")

	filter := bson.M{
		"companyid": companyID,
		"replay":    "",
	}

	cur, err := reportsCollection.Find(context.TODO(), filter)
	if err != nil {
		http.Error(res, "Error fetching reports from database", http.StatusInternalServerError)
		return nil, err
	}
	defer cur.Close(context.TODO())

	var reports []model.Report

	for cur.Next(context.TODO()) {
		var report model.Report
		err := cur.Decode(&report)
		if err != nil {
			http.Error(res, "Error decoding report data", http.StatusInternalServerError)
			return nil, err
		}
		reports = append(reports, report)
	}

	if err := cur.Err(); err != nil {
		http.Error(res, "Error during cursor iteration", http.StatusInternalServerError)
		return nil, err
	}

	return reports, nil
}

func GetAllUserGrades(userID primitive.ObjectID, res http.ResponseWriter) ([]model.Grade, error) {
	reportsCollection := GetClient().Database("MedicinskaOprema").Collection("Grades")

	filter := bson.M{
		"userid": userID,
	}

	cur, err := reportsCollection.Find(context.TODO(), filter)
	if err != nil {
		http.Error(res, "Error fetching reports from database", http.StatusInternalServerError)
		return nil, err
	}
	defer cur.Close(context.TODO())

	var reports []model.Grade

	for cur.Next(context.TODO()) {
		var report model.Grade
		err := cur.Decode(&report)
		if err != nil {
			http.Error(res, "Error decoding report data", http.StatusInternalServerError)
			return nil, err
		}
		reports = append(reports, report)
	}

	if err := cur.Err(); err != nil {
		http.Error(res, "Error during cursor iteration", http.StatusInternalServerError)
		return nil, err
	}

	return reports, nil
}

func GetReportByID(id primitive.ObjectID, res http.ResponseWriter) (model.Report, bool) {
	collection := GetClient().Database("MedicinskaOprema").Collection("Reports")
	filter := bson.M{"_id": id}
	var result model.Report
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			json.NewEncoder(res).Encode("Didnt find user!")
			return result, false
		}
	}
	return result, true
}
func GetGradeByID(id primitive.ObjectID, res http.ResponseWriter) (model.Grade, bool) {
	collection := GetClient().Database("MedicinskaOprema").Collection("Grades")
	filter := bson.M{"_id": id}
	var result model.Grade
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			json.NewEncoder(res).Encode("Didnt find grade!")
			return result, false
		}
	}
	return result, true
}

func GetPastGrade(userID primitive.ObjectID, companyID primitive.ObjectID, res http.ResponseWriter) (model.Grade, bool) {
	collection := GetClient().Database("MedicinskaOprema").Collection("Grades")
	fmt.Println("nista")
	filter := bson.M{"userid": userID, "companyid": companyID}
	var result model.Grade
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return result, false
		}
	}
	return result, true
}

func ReportReplay(id primitive.ObjectID, description string, res http.ResponseWriter) (model.Report, bool) {
	reportsCollection := GetClient().Database("MedicinskaOprema").Collection("Reports")

	filter := bson.M{"_id": id}

	update := bson.M{"$set": bson.M{"replay": description}}

	result, err := reportsCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		http.Error(res, "Error updating report in database", http.StatusInternalServerError)
		fmt.Println("Error updating report:", err)
		return model.Report{}, false
	}

	if result.MatchedCount == 0 {
		http.Error(res, "No report found with the provided ID", http.StatusNotFound)
		fmt.Println("No report found with the provided ID")
		return model.Report{}, false
	}

	var updatedReport model.Report
	err = reportsCollection.FindOne(context.TODO(), filter).Decode(&updatedReport)
	if err != nil {
		http.Error(res, "Error fetching updated report from database", http.StatusInternalServerError)
		fmt.Println("Error fetching updated report:", err)
		return model.Report{}, false
	}

	return updatedReport, true
}

func UpdateGrade(gradeID primitive.ObjectID, updates map[string]interface{}) bool {
	gradesCollection := GetClient().Database("MedicinskaOprema").Collection("Grades")

	updateDoc := bson.M{"$set": updates}

	result, err := gradesCollection.UpdateOne(
		context.TODO(),
		bson.M{"_id": gradeID},
		updateDoc,
	)

	if err != nil {
		return false
	}

	if result.MatchedCount == 0 {
		return false
	}

	return true
}
