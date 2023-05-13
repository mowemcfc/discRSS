package handlers


import (
	"fmt"
	"strconv"
	"log"
	"net/http"

	"github.com/mowemcfc/discRSS/models"
	"github.com/mowemcfc/discRSS/internal/response"
	"github.com/mowemcfc/discRSS/internal/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gin-gonic/gin"
)


func (app *App) FetchUser(userID int) (*models.UserAccount, error) {
	getUserInput := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"userId": {
				N: aws.String(strconv.Itoa(userID)),
			},
		},
		TableName: aws.String(config.UserTableName),
	}

	user, err := app.DdbSvc.GetItem(getUserInput)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
        log.Printf("error fetching user %d: %s %s", userID, dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
        log.Printf("error fetching user %d: %s %s", userID, dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeRequestLimitExceeded:
				log.Printf("error fetching user %d: %s %s", userID, dynamodb.ErrCodeRequestLimitExceeded, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				log.Printf("error fetching user %d: %s %s", userID, dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				log.Printf("error fetching user %d: %s", userID, aerr.Error())
			}
		} else {
			log.Printf("error fetching user %d: %s", userID, err.Error())
		}
    return nil, fmt.Errorf("error fetching user %d: %s", userID, err.Error())
	}

	unmarshalled := models.UserAccount{}
	err = dynamodbattribute.UnmarshalMap(user.Item, &unmarshalled)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling returned user item: %s", err)
	}

	return &unmarshalled, nil
}

func (app *App) PutUser(user *models.UserAccount) error {
	marshalledUser, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return fmt.Errorf("error marshalling user into ddb Item: %s", err.Error())
	}

	input := &dynamodb.PutItemInput{
		Item:                   marshalledUser,
		ReturnConsumedCapacity: aws.String("TOTAL"),
		TableName:              aws.String(config.UserTableName),
	}

	_, err = app.DdbSvc.PutItem(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				return fmt.Errorf("%s %s", dynamodb.ErrCodeConditionalCheckFailedException, aerr.Error())
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				return fmt.Errorf("%s %s", dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				return fmt.Errorf("%s %s", dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeItemCollectionSizeLimitExceededException:
				return fmt.Errorf("%s %s", dynamodb.ErrCodeItemCollectionSizeLimitExceededException, aerr.Error())
			case dynamodb.ErrCodeTransactionConflictException:
				return fmt.Errorf("%s %s", dynamodb.ErrCodeTransactionConflictException, aerr.Error())
			case dynamodb.ErrCodeRequestLimitExceeded:
				return fmt.Errorf("%s %s", dynamodb.ErrCodeRequestLimitExceeded, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				return fmt.Errorf("%s %s", dynamodb.ErrCodeInternalServerError, aerr.Error())
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

func (app *App) GetUserHandler(c *gin.Context) {
	appG := response.Gin{C: c}

	requestUserID, err := strconv.Atoi(appG.C.Param("userId"))
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("userId: %d\n", requestUserID)

	user, err := app.FetchUser(requestUserID)
	if err != nil {
		log.Println("error fetching user from DDB", err)
		return
	}

	log.Printf("user %s channels: %+v", user.UserID, user.ChannelList)
	log.Printf("user %s feeds: %+v", user.UserID, user.FeedList)

	appG.Response(http.StatusOK, user)
}


func (app *App) AddUserHandler(c *gin.Context) {
	appG := response.Gin{C: c}

	var createUserParams models.UserAccount
	if err := appG.C.BindJSON(&createUserParams); err != nil {
		log.Println("error binding user params JSON to models.UserAccount struct", err)
		return
	}

	log.Println(createUserParams.UserID)

	err := app.PutUser(&createUserParams)
	if err != nil {
		log.Println("error putting using in DDB", err)
		return
	}

	appG.Response(http.StatusOK, createUserParams)
}

func (app *App) DeleteUserHandler(c *gin.Context) {
	appG := response.Gin{C: c}
	appG.Response(http.StatusNotImplemented, "Method DELETE for resource /user not implemented")
}

