package Repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	model "main/Model"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func FindUserById(id primitive.ObjectID, res http.ResponseWriter) (model.User, bool) {
	collection := GetClient().Database("MedicinskaOprema").Collection("Users")
	filter := bson.M{"_id": id}
	var result model.User
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
		if err == mongo.ErrNoDocuments {
			json.NewEncoder(res).Encode("Didnt find user!")
			return result, false
		}
	}
	return result, true
}

func GetUserData(email string) (model.User, error) {
	UsersCollection := GetClient().Database("MedicinskaOprema").Collection("Users")
	filter := bson.M{"email": email}
	var result model.User
	err = UsersCollection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		var res http.ResponseWriter
		json.NewEncoder(res).Encode("Didnt find user!")
		return model.User{}, err
	}

	return result, nil
}

func FindUserEmail(email string) bool {
	UsersCollection := GetClient().Database("MedicinskaOprema").Collection("Users")
	filter := bson.M{"email": email}
	var result model.User
	err = UsersCollection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		}
		fmt.Println(err)
	}
	return true
}

func ValidUser(email string, password string) bool {
	UsersCollection := GetClient().Database("MedicinskaOprema").Collection("Users")
	filter := bson.M{"email": email}
	var result model.User
	err = UsersCollection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		}
		log.Fatal(err)
	}
	if CheckPasswordHash(password, result.Password) {
		return true
	}
	return false
}

func SaveUserApplication(email string, firstname string, lastname string, phonenumber string, password string, city string, country string) primitive.ObjectID {
	UsersCollection := GetClient().Database("MedicinskaOprema").Collection("Users")
	location := model.Location{
		City:    city,
		Country: country,
	}
	id := primitive.NewObjectID()
	user := model.User{
		ID:             id,
		Email:          email,
		Role:           "User",
		Firstname:      firstname,
		Lastname:       lastname,
		Phonenumber:    phonenumber,
		Password:       password,
		Location:       location,
		Verified:       false,
		CompmnayID:     primitive.NilObjectID,
		Profession:     "",
		Penalies:       0,
		Points:         0,
		LoyaltyProgram: primitive.NilObjectID,
	}

	insertResult, err := UsersCollection.InsertOne(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Added new user with ID:", insertResult.InsertedID)
	return id

}

func SaveCompanyAdminApplication(email string, firstname string, lastname string, phonenumber string, password string, city string, country string, companyID primitive.ObjectID, role string, profession string) primitive.ObjectID {
	UsersCollection := GetClient().Database("MedicinskaOprema").Collection("Users")
	location := model.Location{
		City:    city,
		Country: country,
	}
	id := primitive.NewObjectID()
	user := model.User{
		ID:             id,
		Email:          email,
		Role:           role,
		Firstname:      firstname,
		Lastname:       lastname,
		Phonenumber:    phonenumber,
		Password:       password,
		Location:       location,
		Verified:       false,
		CompmnayID:     companyID,
		Profession:     profession,
		Penalies:       0,
		Points:         0,
		LoyaltyProgram: primitive.NilObjectID,
	}

	insertResult, err := UsersCollection.InsertOne(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Added new user with ID:", insertResult.InsertedID)
	return id

}

func SaveAdminApplication(email string, firstname string, lastname string, phonenumber string, password string, city string, country string, role string) primitive.ObjectID {
	UsersCollection := GetClient().Database("MedicinskaOprema").Collection("Users")
	location := model.Location{
		City:    city,
		Country: country,
	}
	id := primitive.NewObjectID()
	user := model.User{
		ID:             id,
		Email:          email,
		Role:           role,
		Firstname:      firstname,
		Lastname:       lastname,
		Phonenumber:    phonenumber,
		Password:       password,
		Location:       location,
		Verified:       false,
		CompmnayID:     primitive.NilObjectID,
		Profession:     "",
		Penalies:       0,
		Points:         0,
		LoyaltyProgram: primitive.NilObjectID,
	}

	// Adding user to the database
	insertResult, err := UsersCollection.InsertOne(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Added new user with ID:", insertResult.InsertedID)
	return id

}

func VerifyUser(userID primitive.ObjectID) bool {
	UsersCollection := GetClient().Database("MedicinskaOprema").Collection("Users")
	filter := bson.M{"_id": userID}

	update := bson.M{"$set": bson.M{"verified": true}}

	result, err := UsersCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return false
	}

	if result.MatchedCount == 0 {
		return false
	}

	return true
}

func GetAllUserReports(userID primitive.ObjectID, res http.ResponseWriter) ([]model.Report, bool) {
	reportsCollection := GetClient().Database("MedicinskaOprema").Collection("Reports")

	filter := bson.M{
		"userid": userID,
	}

	cur, err := reportsCollection.Find(context.TODO(), filter)
	if err != nil {
		http.Error(res, "Error fetching reports from database", http.StatusInternalServerError)
		return nil, false
	}
	defer cur.Close(context.TODO())

	var reports []model.Report

	for cur.Next(context.TODO()) {
		var report model.Report
		err := cur.Decode(&report)
		if err != nil {
			http.Error(res, "Error decoding report data", http.StatusInternalServerError)
			return nil, false
		}
		reports = append(reports, report)
	}

	if err := cur.Err(); err != nil {
		http.Error(res, "Error during cursor iteration", http.StatusInternalServerError)
		return nil, false
	}

	return reports, true

}

func SaveLoyaltyProgram(name string, discount float64, upValue int32, downValue int32, res http.ResponseWriter) bool {
	LPCollection := GetClient().Database("MedicinskaOprema").Collection("LoyaltyPrograms")
	id := primitive.NewObjectID()
	lp := model.LoyaltyProgram{
		ID:        id,
		Name:      name,
		Discount:  discount,
		UpValue:   upValue,
		DownValue: downValue,
	}

	insertResult, err := LPCollection.InsertOne(context.Background(), lp)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Added new loyalty program with ID:", insertResult.InsertedID)
	return true

}

func SetLoyaltyProgram(userObjectID primitive.ObjectID, res http.ResponseWriter) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userCollection := GetClient().Database("MedicinskaOprema").Collection("Users")
	loyaltyProgramCollection := GetClient().Database("MedicinskaOprema").Collection("LoyaltyPrograms")

	var user bson.M
	err := userCollection.FindOne(ctx, bson.M{"_id": userObjectID}).Decode(&user)
	if err != nil {
		fmt.Println("Error finding user:", err)
		http.Error(res, "User not found", http.StatusNotFound)
		return false
	}

	points, okPoints := user["points"].(int32)
	penalties, okPenalties := user["penalties"].(int32)
	if !okPoints || !okPenalties {
		fmt.Println("Invalid points or penalties format in user")
		http.Error(res, "Invalid user data", http.StatusInternalServerError)
		return false
	}

	netPoints := points - penalties

	cursor, err := loyaltyProgramCollection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println("Error finding loyalty programs:", err)
		http.Error(res, "Loyalty programs not found", http.StatusInternalServerError)
		return false
	}
	defer cursor.Close(ctx)

	var matchingLoyaltyProgram primitive.ObjectID = primitive.NilObjectID
	for cursor.Next(ctx) {
		var loyaltyProgram bson.M
		if err := cursor.Decode(&loyaltyProgram); err != nil {
			fmt.Println("Error decoding loyalty program:", err)
			continue
		}

		upValue, okUpValue := loyaltyProgram["upValue"].(int32)
		downValue, okDownValue := loyaltyProgram["downValue"].(int32)
		if !okUpValue || !okDownValue {
			fmt.Println("Invalid upValue or downValue format in loyalty program")
			continue
		}

		if netPoints >= downValue && netPoints <= upValue {
			matchingLoyaltyProgram = loyaltyProgram["_id"].(primitive.ObjectID)
			break
		}
	}

	_, err = userCollection.UpdateOne(
		ctx,
		bson.M{"_id": userObjectID},
		bson.M{"$set": bson.M{"loyaltyProgram": matchingLoyaltyProgram}},
	)
	if err != nil {
		fmt.Println("Error updating user loyalty program:", err)
		http.Error(res, "Error updating user loyalty program", http.StatusInternalServerError)
		return false
	}

	res.WriteHeader(http.StatusOK)
	res.Write([]byte("Loyalty program updated successfully"))
	return true
}
