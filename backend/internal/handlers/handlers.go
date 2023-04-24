package handlers

import (
	"context"
	"net/http"
  "os"
  "fmt"

	"github.com/mowemcfc/discRSS/internal/response"
	"github.com/mowemcfc/discRSS/internal/sessions"
	"github.com/mowemcfc/discRSS/models"

	"github.com/aws/aws-sdk-go/service/dynamodb"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-lambda-go/events"
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
