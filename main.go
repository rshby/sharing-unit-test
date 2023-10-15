package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/opentracing/opentracing-go"
	"log"
	accountHandler "sharingunittest/account/handler"
	accRepository "sharingunittest/account/repository"
	accService "sharingunittest/account/service"
	database "sharingunittest/database/connection"
	"sharingunittest/router"
	"sharingunittest/tracing"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("error cant load env file")
	}

	log.Println("success load env file")
}

func main() {
	db := database.CreateConnectionDB()

	// register jaeger
	tracer, closer := tracing.ConnectJaeger("ut-service")
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

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
