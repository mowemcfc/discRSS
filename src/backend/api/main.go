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
	"github.com/gin-gonic/gin"
	adapter "github.com/gwatts/gin-adapter"
	"github.com/joho/godotenv"
)

type Feed struct {
	FeedID     int    `json:"feedId"`
	Title      string `json:"title"`
	Url        string `json:"url"`
	TimeFormat string `json:"timeFormat"`
}

type UserAccount struct {
	UserID      int              `json:"userId"`
	Username    string           `json:"username"`
	FeedList    []Feed           `json:"feedList"`
	ChannelList []DiscordChannel `json:"channelList"`
}

type DiscordChannel struct {
	ChannelName string `json:"channelName"`
	ServerName  string `json:"serverName"`
	ChannelID   int    `json:"channelID"`
}

type AppConfig struct {
	AppName               string `json:"appName"`
	LastCheckedTime       string `json:"lastCheckedTime"`
	LastCheckedTimeFormat string `json:"lastCheckedTimeFormat"`
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

func fetchAppConfig(appName string) (*AppConfig, error) {

	getAppConfigInput := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"appName": {
				S: aws.String(appName),
			},
		},
		TableName: aws.String(APP_CONFIG_TABLE_NAME),
	}

	config, err := ddbSvc.GetItem(getAppConfigInput)
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

	unmarshalled := AppConfig{}
	err = dynamodbattribute.UnmarshalMap(config.Item, &unmarshalled)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling returned appconfig item: %s", err)
	}

	return &unmarshalled, nil
}

func fetchUser(userID int) (*UserAccount, error) {

	getUserInput := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"userID": {
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
		return fmt.Errorf(err.Error())
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

func userGetHandler(c *gin.Context) {
	requestUserID, err := strconv.Atoi(c.Request.URL.Query().Get("userId"))
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("userID: %d\n", requestUserID)

	user, err := fetchUser(requestUserID)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("user %d channels: %+v", user.UserID, user.ChannelList)
	log.Printf("user %d feeds: %+v", user.UserID, user.FeedList)

	marshalledUser, err := json.Marshal(user)
	log.Printf("\nmarshalled: %s\n", string(marshalledUser))
	if err != nil {
		log.Println(err)
		return
	}

	c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
	c.JSON(http.StatusOK, events.APIGatewayProxyResponse{
		StatusCode:      http.StatusOK,
		IsBase64Encoded: false,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(marshalledUser),
	})
}

func userPostHandler(c *gin.Context) {

	// TODO: Error handling
	//	- IDOR
	//  - gracefully send error response

	var createUserParams UserAccount
	err := c.BindJSON(&createUserParams)
	if err != nil {
		log.Println(err)
		return
	}

	user, err := fetchUser(createUserParams.UserID) // TODO: replace with request body userId
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(createUserParams.UserID)
	log.Printf("%v", user)

	marshalledUser, err := json.Marshal(createUserParams)
	log.Printf("\nmarshalled: %s\n", string(marshalledUser))
	if err != nil {
		log.Println(err)
		return
	}

	err = putUser(&createUserParams)
	if err != nil {
		log.Println(err)
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

func corsPreflightHandler(c *gin.Context) {
	c.JSON(http.StatusOK, events.APIGatewayProxyResponse{
		StatusCode:      http.StatusOK,
		IsBase64Encoded: false,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: "Hello from discRSS",
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

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Header("Access-Control-Allow-Methods", "POST, PATCH, PUT, GET, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "*, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")
	}
}

func main() {
	isLocal = os.Getenv("LAMBDA_TASK_ROOT") == ""

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
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
	g.GET("/hello", corsMiddleware(), helloWorldHandler)
	g.GET("/user", corsMiddleware(), jwtMiddleware, userGetHandler)
	g.POST("/user", corsMiddleware(), userPostHandler)
	g.OPTIONS("/user", corsMiddleware(), corsPreflightHandler)

	awsSession, err := getAWSSession()
	if err != nil {
		log.Println(err)
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
