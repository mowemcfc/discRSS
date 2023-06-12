package dynamodb

import (
  "log"
  "fmt"

  "github.com/mowemcfc/discRSS/models"
  "github.com/mowemcfc/discRSS/internal/config"

  "github.com/gin-gonic/gin"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DynamoDBUserRepository struct {
  client dynamodbiface.DynamoDBAPI
}

type UserRepository interface {
  GetUser(gin.Context) (*models.UserAccount, error)
  CreateUser(gin.Context) (error)
  AddFeed(gin.Context) (*models.Feed, error)
  UpdateFeed(gin.Context) (*models.Feed, error)
  RemoveFeed(gin.Context) (error)
}

func NewDynamoDBUserRepository (client dynamodbiface.DynamoDBAPI) UserRepository {
  return &DynamoDBUserRepository{client}
}

func (d *DynamoDBUserRepository) GetUser(ctx *gin.Context) (*models.UserAccount, error) { 
  userId := ctx.Query("userId")
	getUserInput := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"userId": {
				N: aws.String(userId),
			},
		},
		TableName: aws.String(config.UserTableName),
	}

	user, err := d.client.GetItem(getUserInput)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
				log.Printf("error fetching user %s: %s", userId, aerr.Error())
		} else {
			log.Printf("error fetching user %s: %s", userId, err.Error())
		}
    return nil, fmt.Errorf("error fetching user %s: %s", userId, err.Error())
	}

	unmarshalled := models.UserAccount{}
	err = dynamodbattribute.UnmarshalMap(user.Item, &unmarshalled)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling returned user item: %s", err)
	}

	return &unmarshalled, nil
}
func (d *DynamoDBUserRepository) CreateUser(ctx *gin.Context) (error) { return nil }
func (d *DynamoDBUserRepository) AddFeed(ctx *gin.Context) (*models.Feed, error) { return &models.Feed{}, nil }
func (d *DynamoDBUserRepository) UpdateFeed(ctx *gin.Context) (*models.Feed, error) { return &models.Feed{}, nil }
func (d *DynamoDBUserRepository) RemoveFeed(ctx *gin.Context) (error) { return nil }
