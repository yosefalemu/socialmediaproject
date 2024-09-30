package controllers

import (
	"api/api/middlewares"
	"api/models"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	DB     *gorm.DB
	Router *gin.Engine
}

var errList = make(map[string]string)

func (server *Server) Intialize(DbDriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {
	var err error
	if DbDriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)
		server.DB, err = gorm.Open(mysql.Open(DBURL), &gorm.Config{})
		if err != nil {
			fmt.Printf("Cannot connetct to %s database", DbDriver)
			log.Fatal("This is the error", err)
		} else {
			fmt.Printf("We are connected to the %s database", DbDriver)
		}

	} else if DbDriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
		server.DB, err = gorm.Open(postgres.Open(DBURL), &gorm.Config{})
		if err != nil {
			fmt.Printf("Cannot connect to %s database", DbDriver)
			log.Fatal("This is the error", err)
		} else {
			fmt.Printf("We are connected to the %s database", DbDriver)
		}
	} else {
		fmt.Println("UNKOWN DRIVER")
	}

	//DATABASE MIGRATIONS
	server.DB.Debug().AutoMigrate(
		&models.User{},
	)
	server.Router = gin.Default()
	server.Router.Use(middlewares.CORSMiddleware())
	server.initializeRoutes()
}

func (server *Server) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
