package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	controllers "main/Controller"
	conn "main/Repository"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}

}

func main() {
	conn.InitConnection()
	router := mux.NewRouter()

	userController := controllers.NewUserController()
	userController.RegisterRoutes()
	companyController := controllers.NewCompanyController()
	companyController.RegisterRoutes()
	assetController := controllers.NewAssetController()
	assetController.RegisterRoutes()
	reservationController := controllers.NewReservationController()
	reservationController.RegisterRoutes()
	statisticController := controllers.NewStatisticController()
	statisticController.RegisterRoutes()

	router.PathPrefix("/users").Handler(userController.Router)
	router.PathPrefix("/companies").Handler(companyController.Router)
	router.PathPrefix("/assets").Handler(assetController.Router)
	router.PathPrefix("/reservations").Handler(reservationController.Router)
	router.PathPrefix("/statistics").Handler(statisticController.Router)

	fmt.Println("[ UP ON PORT 8080 ]")
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}
	err := http.ListenAndServe(fmt.Sprintf(":%s", PORT), router)
	log.Fatal(err)
	defer func() {
		if err := conn.GetClient().Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

}
