package main

import (
	"log"
	"os"

	"github.com/mowemcfc/discRSS/internal/auth0"
	"github.com/mowemcfc/discRSS/internal/handlers"

	"github.com/aws/aws-lambda-go/lambda"
	ginLambdaAdapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	adapter "github.com/gwatts/gin-adapter"

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

	var jwtMiddleware gin.HandlerFunc
	if app.IsLocal {
		jwtMiddleware = func(c *gin.Context) {}
	} else {
		jwtMiddleware = adapter.Wrap(auth0.EnsureValidToken())
	}

	log.Println("Configuring API methods")
	app.Engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:9001", "http://localhost:3000"},
		AllowMethods:     []string{"POST", "PATCH", "PUT", "DELETE", "GET", "OPTIONS"},
		AllowHeaders:     []string{"*", "Authorization"},
		AllowCredentials: true,
	}))
	app.Engine.GET("/hello", app.HelloWorldHandler)

	userRoute := app.Engine.Group("/user")
	{
		userRoute.GET("/:userId", jwtMiddleware, app.GetUserHandler)
		userRoute.POST("/:userId", jwtMiddleware, app.AddUserHandler)
		userRoute.DELETE("/:userId", jwtMiddleware, app.DeleteUserHandler)

		userRoute.POST("/:userId/feeds", jwtMiddleware)

		userRoute.GET("/:userId/feed/:feedId", jwtMiddleware, app.GetFeedHandler)
		userRoute.POST("/:userId/feed", jwtMiddleware, app.AddFeedHandler)
		userRoute.DELETE("/:userId/feed/:feedId", jwtMiddleware, app.DeleteFeedHandler)
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

