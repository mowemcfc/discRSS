package app

import (
	"github.com/mowemcfc/discRSS/models"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	ginLambdaAdapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"

	"github.com/gin-gonic/gin"
)

type App struct {
  *gin.Engine
  ginLambda *ginLambdaAdapter.GinLambda
  awsSession *session.Session
  ddbSvc *dynamodb.DynamoDB
  secretsManagerSvc *secretsmanager.SecretsManager
  appConfig *models.AppConfig
  isLocal bool
}

