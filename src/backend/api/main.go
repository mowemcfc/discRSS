package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	ginLambdaAdapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	adapter "github.com/gwatts/gin-adapter"
	"github.com/joho/godotenv"
)

type Feed struct {
	FeedID     json.Number `json:"feedId" dynamodbav:"feedId"`
	Title      string      `json:"title" dynamodbav:"title"`
	Url        string      `json:"url" dynamodbav:"url"`
	TimeFormat string      `json:"timeFormat" dynamodbav:"timeFormat"`
}

type UserAccount struct {
	UserID      json.Number      `json:"userId" dynamodbav:"userId"`
	Username    string           `json:"username" dynamodbav:"username"`
	FeedList    []Feed           `json:"feedList" dynamodbav:"feedList"`
	ChannelList []DiscordChannel `json:"channelList" dynamodbav:"channelList"`
}

type DiscordChannel struct {
	ChannelName string `json:"channelName" dynamodbav:"channelName"`
	ServerName  string `json:"serverName" dynamodbav:"serverName"`
	ChannelID   int    `json:"channelID" dynamodbav:"channelID"`
}

type AppConfig struct {
	AppName               string `json:"appName" dynamodbav:"appName"`
	LastCheckedTime       string `json:"lastCheckedTime" dynamodbav:"lastCheckedTime"`
	LastCheckedTimeFormat string `json:"lastCheckedTimeFormat" dynamodbav:"lastCheckedTimeFormat"`
}

var discRssConfig *AppConfig

var ginLambda *ginLambdaAdapter.GinLambda

var awsSession *session.Session
var secretsmanagerSvc *secretsmanager.SecretsManager
var ddbSvc *dynamodb.DynamoDB

var isLocal bool

const APP_NAME string = "discRSS"
const APP_CONFIG_TABLE_NAME string = "discRSS-AppConfigs"
const USER_TABLE_NAME string = "discRSS-UserRecords"
const BOT_TOKEN_SECRET_NAME string = "discRSS/discord-bot-secret"

func fetchUser(userID int) (*UserAccount, error) {

	getUserInput := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"userId": {
				N: aws.String(strconv.Itoa(userID)),
			},
		},
		TableName: aws.String(USER_TABLE_NAME),
	}

	user, err := ddbSvc.GetItem(getUserInput)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				return nil, fmt.Errorf("%s %s", dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				return nil, fmt.Errorf("%s %s", dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeRequestLimitExceeded:
				return nil, fmt.Errorf("%s %s", dynamodb.ErrCodeRequestLimitExceeded, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				return nil, fmt.Errorf("%s %s", dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				return nil, fmt.Errorf("%s", aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			return nil, fmt.Errorf(err.Error())
		}
	}

	unmarshalled := UserAccount{}
	err = dynamodbattribute.UnmarshalMap(user.Item, &unmarshalled)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling returned user item: %s", err)
	}

	return &unmarshalled, nil
}

func putUser(user *UserAccount) error {
	marshalledUser, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return fmt.Errorf("error marshalling user into ddb Item: %s", err.Error())
	}

	input := &dynamodb.PutItemInput{
		Item:                   marshalledUser,
		ReturnConsumedCapacity: aws.String("TOTAL"),
		TableName:              aws.String(USER_TABLE_NAME),
	}

	_, err = ddbSvc.PutItem(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				return fmt.Errorf(dynamodb.ErrCodeConditionalCheckFailedException, aerr.Error())
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				return fmt.Errorf(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				return fmt.Errorf(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeItemCollectionSizeLimitExceededException:
				return fmt.Errorf(dynamodb.ErrCodeItemCollectionSizeLimitExceededException, aerr.Error())
			case dynamodb.ErrCodeTransactionConflictException:
				return fmt.Errorf(dynamodb.ErrCodeTransactionConflictException, aerr.Error())
			case dynamodb.ErrCodeRequestLimitExceeded:
				return fmt.Errorf(dynamodb.ErrCodeRequestLimitExceeded, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				return fmt.Errorf(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				return fmt.Errorf(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			return fmt.Errorf(err.Error())
		}
	}

	return nil
}

func getUserHandler(c *gin.Context) {
	requestUserID, err := strconv.Atoi(c.Request.URL.Query().Get("userId"))
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("userId: %d\n", requestUserID)

	user, err := fetchUser(requestUserID)
	if err != nil {
		log.Println("error fetching user from DDB", err)
		return
	}

	log.Printf("user %d channels: %+v", user.UserID, user.ChannelList)
	log.Printf("user %d feeds: %+v", user.UserID, user.FeedList)

	marshalledUser, err := json.Marshal(user)
	log.Printf("\nmarshalled: %s\n", string(marshalledUser))
	if err != nil {
		log.Println("error marshalling user to JSON object", err)
		return
	}

	c.JSON(http.StatusOK, events.APIGatewayProxyResponse{
		StatusCode:      http.StatusOK,
		IsBase64Encoded: false,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(marshalledUser),
	})
}

func addUserHandler(c *gin.Context) {

	// TODO: Error handling
	//	- IDOR
	//  - gracefully send error response

	var createUserParams UserAccount
	if err := c.BindJSON(&createUserParams); err != nil {
		log.Printf("error binding user params JSON to UserAccount struct", err)
		return
	}

	log.Println(createUserParams.UserID)

	marshalledUser, err := json.Marshal(createUserParams)
	log.Printf("\nmarshalled: %s\n", string(marshalledUser))
	if err != nil {
		log.Println("error marshalling user to JSON object", err)
		return
	}

	err = putUser(&createUserParams)
	if err != nil {
		log.Println("error putting using in DDB", err)
		return
	}

	c.JSON(http.StatusOK, events.APIGatewayProxyResponse{
		StatusCode:      http.StatusOK,
		IsBase64Encoded: false,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(marshalledUser),
	})
}

func addFeedHandler(c *gin.Context) {
	// TODO:
	// recv payload, parse userId and
	// marshall feed json into feed struct
	// if feedId already exists, error
	// fetch user profile
	// append marshalled field to feedList
	// return 200
	// return appropriate errors as required

	addFeedParams := struct {
		UserId  string `json:"userId"`
		NewFeed []Feed `json:"newFeed"`
	}{}

	if err := c.BindJSON(&addFeedParams); err != nil {
		log.Println("error binding addFeed params JSON to addFeedParams struct", err)
		return
	}

	marshalledFeed, err := dynamodbattribute.Marshal(addFeedParams.NewFeed)
	if err != nil {
		log.Println("error marshalling feed struct into dynamodbattribute map", err)
		return
	}

	addFeedInput := &dynamodb.UpdateItemInput{
		ExpressionAttributeNames: map[string]*string{
			"#fL": aws.String("feedList"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":f": marshalledFeed,
		},
		Key: map[string]*dynamodb.AttributeValue{
			"userId": {
				N: aws.String(addFeedParams.UserId),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("SET #fL = list_append(#fL, :f)"),
		TableName:        aws.String(USER_TABLE_NAME),
	}

	_, err = ddbSvc.UpdateItem(addFeedInput)
	if err != nil {
		log.Printf("error updating user: %s's feed list with feed: %v, %s\n", addFeedParams.UserId, addFeedParams.NewFeed, err.Error())
		return
	}

	jsonMarshalledFeed, err := json.Marshal(addFeedParams.NewFeed)
	if err != nil {
		log.Printf("error converting addFeedParams.NewFeed to json string", err)
		return
	}

	c.JSON(http.StatusOK, events.APIGatewayProxyResponse{
		StatusCode:      http.StatusOK,
		IsBase64Encoded: false,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(jsonMarshalledFeed),
	})
}

func helloWorldHandler(c *gin.Context) {
	c.JSON(http.StatusOK, events.APIGatewayProxyResponse{
		StatusCode:      http.StatusOK,
		IsBase64Encoded: false,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: "Hello, World!",
	})
}

func notFoundHandler(c *gin.Context) {
	c.JSON(http.StatusNotFound, events.APIGatewayProxyResponse{
		StatusCode:      http.StatusNotFound,
		IsBase64Encoded: false,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: fmt.Sprintf("Not Found: %s", c.Request.URL.Path),
	})
}

func main() {
	isLocal = os.Getenv("LAMBDA_TASK_ROOT") == ""

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
		os.Exit(1)
	}

	g := gin.Default()
	var jwtMiddleware gin.HandlerFunc
	if isLocal {
		jwtMiddleware = func(c *gin.Context) {}
	} else {
		jwtMiddleware = adapter.Wrap(EnsureValidToken())
	}

	log.Println("Configuring API methods")
	g.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"POST", "PATCH", "PUT", "GET", "OPTIONS"},
		AllowHeaders:     []string{"*", "Authorization"},
		AllowCredentials: true,
	}))
	g.GET("/hello", helloWorldHandler)

	userRoute := g.Group("/user")
	{
		userRoute.GET("", jwtMiddleware, getUserHandler)
		userRoute.POST("", jwtMiddleware, addUserHandler)

		userRoute.GET("/feeds", jwtMiddleware /*, removeFeedHandler */)
		userRoute.POST("/feeds", jwtMiddleware, addFeedHandler)
	}

	awsSession, err := getAWSSession()
	if err != nil {
		log.Println("error opening AWS session: ", err)
		return
	}
	ddbSvc = dynamodb.New(awsSession)
	secretsmanagerSvc = secretsmanager.New(awsSession)

	if isLocal {
		log.Println("Inside LOCAL environment, using default router")
		g.Run("0.0.0.0:9001")
	} else {
		log.Println("Inside REMOTE lambda environment, using ginLambda router")
		ginLambda = ginLambdaAdapter.New(g)
		lambda.Start(lambdaHandler)
	}
}

//func localHandler(ctx context.Context, request events.Lambda)

func lambdaHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, request)
}
