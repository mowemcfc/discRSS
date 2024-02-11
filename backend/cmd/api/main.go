package main

import (
	"log"
	"os"

	"github.com/mowemcfc/discRSS/internal/app"

	"github.com/aws/aws-lambda-go/lambda"
	ginLambdaAdapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
		os.Exit(1)
	}

  app, err := app.NewApp()
  if err != nil {
    log.Fatal("Error instantiating app object", err)
  }

	if app.IsLocal {
		log.Println("Inside LOCAL environment, using default router")
		app.Engine.Run("0.0.0.0:9001")
	} else {
		log.Println("Inside REMOTE lambda environment, using ginLambda router")
		app.GinLambda = ginLambdaAdapter.New(app.Engine)
		lambda.Start(app.LambdaHandler)
	}
}


