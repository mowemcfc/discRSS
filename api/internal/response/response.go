package response

import (
	"github.com/gin-gonic/gin"
)

/*
We cannot define new methods on non-local type *gin.Context, so
instead create a wrapper around it
*/
type Gin struct {
	C *gin.Context
}

/*
In lambda proxy integrations, API Gateway requires responses in the following JSON format.
Not all fields are required.

	{
			"isBase64Encoded": true|false,
			"statusCode": httpStatusCode,
			"headers": { "headerName": "headerValue", ... },
			"multiValueHeaders": { "headerName": ["headerValue", "headerValue2", ...], ... },
			"body": "..."
	}
*/
type ApiGatewayLambdaProxyResponse struct {
	IsBase64Encoded   bool
	StatusCode        int
	Headers           map[string]string
	MultiValueHeaders map[string][]string
	Body              interface{}
}

func (g *Gin) Response(httpCode int, data interface{}) {
	g.C.JSON(httpCode, ApiGatewayLambdaProxyResponse{
		IsBase64Encoded:   false,
		StatusCode:        httpCode,
		Headers:           map[string]string{"Content-Type": "application/json"},
		MultiValueHeaders: nil,
		Body:              data,
	})
}
