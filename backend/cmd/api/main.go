package main

import (
	"log"
	"os"

	"github.com/mowemcfc/discRSS/internal/handlers"

	"github.com/aws/aws-lambda-go/lambda"
	ginLambdaAdapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"

	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
		os.Exit(1)
	}

  app, err := handlers.NewApp()
  if err != nil {
    log.Fatal("Error instantiating app object", err)
  }

	log.Println("Configuring API methods")
	app.Engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:9001", "http://localhost:3000"},
		AllowMethods:     []string{"POST", "PATCH", "PUT", "DELETE", "GET", "OPTIONS"},
		AllowHeaders:     []string{"*", "Authorization"},
		AllowCredentials: true,
	}))

	if app.IsLocal {
		log.Println("Inside LOCAL environment, using default router")
		app.Engine.Run("0.0.0.0:9001")
	} else {
		log.Println("Inside REMOTE lambda environment, using ginLambda router")
		app.GinLambda = ginLambdaAdapter.New(app.Engine)
		lambda.Start(app.LambdaHandler)
	}
}

