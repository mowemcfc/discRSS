package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

  "github.com/mowemcfc/discRSS/internal/response"
  "github.com/mowemcfc/discRSS/internal/dynamodb"

  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockDynamoDBService struct {
	dynamodbiface.DynamoDBAPI
}

func TestAddFeedHandler(t *testing.T) {
	app := &App{
    Engine: gin.Default(),
		DdbSvc: dynamodb.NewMockDynamoDB(),
	}

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.POST("/users/:userId/feeds", app.AddFeedHandler)

	tests := []struct {
		name               string
		userId             int
		addFeedParams      AddFeedParams
		expectedStatusCode int
    expectedResponseBody   map[string]interface{}
	}{
		{
			name: "Valid request",
			userId: 0,
			addFeedParams: AddFeedParams{
				Title: "feed1",
				URL:   "https://feed1.com/rss",
			},
			expectedStatusCode: http.StatusOK,
      expectedResponseBody: map[string]interface{}{
        "title": "feed1",
        "url": "https://feed1.com/rss",
      },
		},
		{
			name: "Invalid URL",
			userId: 1,
			addFeedParams: AddFeedParams{
				Title: "Invalid Feed",
				URL:   "invalid-url",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Invalid UserId",
			userId: -1,
			addFeedParams: AddFeedParams{
				Title: "test",
        URL:   "https://www.samplefeed.com/rss",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Non-existent user",
			userId: 9999,
			addFeedParams: AddFeedParams{
				Title: "test",
        URL:   "https://www.samplefeed.com/rss",
			},
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(test.addFeedParams)
			req, _ := http.NewRequest("POST", fmt.Sprintf("/users/%d/feeds", test.userId), bytes.NewReader(jsonBody))
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, test.expectedStatusCode, resp.Code)

      var response response.ApiGatewayLambdaProxyResponse
			json.Unmarshal(resp.Body.Bytes(), &response)
      fmt.Printf("%v\n", response)

      responseBodyBytes, _ := json.Marshal(response.Body)
      var responseBody map[string]interface{}
      json.Unmarshal(responseBodyBytes, &responseBody)
      
			assert.Equal(t, test.expectedResponseBody["title"], responseBody["title"])
			assert.Equal(t, test.expectedResponseBody["url"], responseBody["url"])
		})
	}
}

func TestGetFeedHandler(t *testing.T) {
	app := &App{
    Engine: gin.Default(),
		DdbSvc: dynamodb.NewMockDynamoDB(),
	}

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.GET("/users/:userId/feeds/:feedId", app.GetFeedHandler)

	tests := []struct {
		name               string
		userId             int
		feedId             string
		expectedStatusCode int
    expectedResponseBody map[string]interface{}
	}{
		{
			name:               "Valid request",
			userId:             0,
			feedId:             "0",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Invalid UserID",
			userId:             -1,
			feedId:             "1",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Invalid feed ID",
			userId:             0,
			feedId:             "invalid_feed_id",
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", fmt.Sprintf("/users/%d/feeds/%s", test.userId, test.feedId), nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, test.expectedStatusCode, resp.Code)
		})
	}
}

func TestDeleteFeedHandler(t *testing.T) {
	app := &App{
    Engine: gin.Default(),
		DdbSvc: dynamodb.NewMockDynamoDB(),
	}

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.DELETE("/users/:userId/feeds/:feedId", app.DeleteFeedHandler)

	tests := []struct {
		name               string
		userId             int
		feedId             string
		expectedStatusCode int
	}{
		{
			name:               "Valid request",
			userId:             0,
			feedId:             "0",
			expectedStatusCode: http.StatusNoContent,
		},
		{
			name:               "Invalid feed ID",
			userId:             0,
			feedId:             "invalid_feed_id",
			expectedStatusCode: http.StatusNoContent,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest("DELETE", fmt.Sprintf("/users/%d/feeds/%s", test.userId, test.feedId), nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, test.expectedStatusCode, resp.Code)
		})
	}
}

