package handlers

import (
	"context"
	"github.com/mowemcfc/discRSS/internal/response"
	"github.com/aws/aws-lambda-go/events"
	"github.com/gin-gonic/gin"
	"net/http"
)

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
