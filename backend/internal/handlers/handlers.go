package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/mowemcfc/discRSS/internal/response"
	"github.com/mowemcfc/discRSS/internal/sessions"
	user "github.com/mowemcfc/discRSS/internal/user/http"
	userDynamoDbRepo "github.com/mowemcfc/discRSS/internal/user/repository/dynamodb"
	"github.com/mowemcfc/discRSS/internal/user/usecase"
	"github.com/mowemcfc/discRSS/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	ginLambdaAdapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

type App struct {
  *gin.Engine
  GinLambda *ginLambdaAdapter.GinLambda
  AwsSession *session.Session
  DdbSvc dynamodbiface.DynamoDBAPI
  SecretsManagerSvc *secretsmanager.SecretsManager
  AppConfig *models.AppConfig
  UserHandler user.UserHandler
  IsLocal bool
}

func NewApp() (*App, error) {
  app := &App{
    Engine: gin.Default(),
    IsLocal: os.Getenv("LAMBDA_TASK_ROOT") == "",
  }

  awsSession, err := sessions.GetAWSSession(app.IsLocal)
	if err != nil {
    return nil, fmt.Errorf("error opening AWS session: %s", err)
	}
  app.AwsSession = awsSession
	app.DdbSvc = dynamodb.New(app.AwsSession)
	app.SecretsManagerSvc = secretsmanager.New(app.AwsSession)
  userRepo := userDynamoDbRepo.NewDynamoDBUserRepository(app.DdbSvc)
  userUsecase := usecase.NewUserUsecase(userRepo)
  app.UserHandler = user.NewUserHandler(app.Engine, userUsecase)

  return app, nil
}


func (app *App) HelloWorldHandler(c *gin.Context) {
	appG := response.Gin{C: c}
	appG.Response(http.StatusOK, "Hello, World!")
}

func (app *App) NotFoundHandler(c *gin.Context) {
	appG := response.Gin{C: c}
	appG.Response(http.StatusNotFound, "Resource not found.")
}


func (app *App) LambdaHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return app.GinLambda.ProxyWithContext(ctx, request)
}
