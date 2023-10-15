package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	accountHandler "sharingunittest/account/handler"
	accRepository "sharingunittest/account/repository"
	accService "sharingunittest/account/service"
	database "sharingunittest/database/connection"
	"sharingunittest/router"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("error cant load env file")
	}

	log.Println("success load env file")
}

func main() {
	db := database.CreateConnectionDB()

	// register repository
	accRepo := accRepository.NewAccountRepository(db)

	// register service
	accService := accService.NewAccountService(accRepo)

	// handler
	accHandler := accountHandler.NewAccountHandler(accService)

	r := gin.Default()

	rv1 := r.Group("/api/v1")
	router.AccountRoutes(rv1, accHandler)

	r.Run(":5000")
}
