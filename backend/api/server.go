package api

import (
	"api/api/controllers"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var server = controllers.Server{}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("sad .env file found")
	}
}

func Run() {
	var err error
	err = godotenv.Load()
	fmt.Println("REACH RUN FUNCTION")
	if err != nil {
		log.Fatalf("Error getting env, %v", err)
	} else {
		fmt.Println("We are getting values")
	}
	server.Intialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))
	airPort := fmt.Sprintf(":%s", os.Getenv("API_PORT"))
	fmt.Printf("Listening to port %s", airPort)
	server.Run(airPort)
}
