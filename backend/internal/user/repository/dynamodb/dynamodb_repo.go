package dynamodb

import (
	"github.com/mowemcfc/discRSS/internal/config"
	"github.com/mowemcfc/discRSS/models"
	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/gin-gonic/gin"

	"go.opentelemetry.io/otel"
)

type DynamoDBUserRepository struct {
  client dynamodbiface.DynamoDBAPI
}

type UserRepository interface {
  GetUser(ctx *gin.Context, userId string) (*models.UserAccount, error)
  CreateUser(ctx *gin.Context, user *models.UserAccount) (*models.UserAccount, error)
  AddFeed(ctx *gin.Context, feed *models.Feed, userId string) (*models.Feed, error)
  GetFeed(ctx *gin.Context, feedId string, userId string) (*models.Feed, error)
  UpdateFeed(ctx *gin.Context, feed *models.Feed) (*models.Feed, error)
  RemoveFeed(ctx *gin.Context, feedId string, userId string) (error)
  ListFeedsAll(ctx *gin.Context, userId string) ([]*models.Feed, error)
}

func NewDynamoDBUserRepository (client dynamodbiface.DynamoDBAPI) UserRepository {
  return &DynamoDBUserRepository{client}
}

func (d *DynamoDBUserRepository) GetUser(ctx *gin.Context, userId string) (*models.UserAccount, error) { 
  tr := otel.Tracer("api_user")
  _, span := tr.Start(ctx, "db.getUser")
  defer span.End()

	getUserInput := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"userId": {
				S: aws.String(userId),
			},
		},
		TableName: aws.String(config.UserTableName),
	}

	user, err := d.client.GetItem(getUserInput)
	if err != nil {
    logrus.Errorf("error getting user from ddb", err)
    return nil, models.ErrInternalServerError
	}

	unmarshalled := models.UserAccount{}
	err = dynamodbattribute.UnmarshalMap(user.Item, &unmarshalled)
	if err != nil {
    logrus.Errorf("error unmarshalling response payload into user struct", err)
		return nil, models.ErrInternalServerError
	}

	return &unmarshalled, nil
}

func (d *DynamoDBUserRepository) CreateUser(ctx *gin.Context, user *models.UserAccount) (*models.UserAccount, error) {  
  dynamoEncoder := dynamodbattribute.NewEncoder(func(e *dynamodbattribute.Encoder) {
		e.EnableEmptyCollections = true
	})
	marshalledUser, err := dynamoEncoder.Encode(user)
	if err != nil {
    logrus.Errorf("error marshalling new user to ddb attribute map: ", err)
		return nil, models.ErrInternalServerError
	}

	input := &dynamodb.PutItemInput{
		Item:                   marshalledUser.M,
		ReturnConsumedCapacity: aws.String("TOTAL"),
		TableName:              aws.String(config.UserTableName),
    ConditionExpression:    aws.String("attribute_not_exists(userId)"),
	}

	_, err = d.client.PutItem(input)
	if err != nil {
    if aerr, ok := err.(awserr.Error); ok {
      if aerr.Code() == dynamodb.ErrCodeConditionalCheckFailedException {
          logrus.Errorf("item already exists in ddb: %v", err)
          return nil, models.ErrConflict
      }
    }
    logrus.Errorf("error putting new user into ddb: ", err)
    return nil, models.ErrInternalServerError
	}

	return user, nil
}

func (d *DynamoDBUserRepository) AddFeed(ctx *gin.Context, feed *models.Feed, userId string) (*models.Feed, error) { 
	marshalledFeed, err := dynamodbattribute.Marshal(feed)
	if err != nil {
    logrus.Errorf("error marshalling feed struct into dynamodbattribute map: ", err)
    return nil, models.ErrInternalServerError
	}

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeNames: map[string]*string{
			"#fID": aws.String(feed.FeedID),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":f": marshalledFeed,
		},
		Key: map[string]*dynamodb.AttributeValue{
			"userId": {
				S: aws.String(userId),
			},
		},
		ConditionExpression: aws.String("attribute_not_exists(feedList.#fID)"),
		UpdateExpression:    aws.String("SET feedList.#fID = :f"),
		TableName:           aws.String(config.UserTableName),
	}

  _, err = d.client.UpdateItem(input)
  if err != nil {
    if aerr, ok := err.(awserr.Error); ok {
      if aerr.Code() == dynamodb.ErrCodeResourceNotFoundException {
        logrus.Errorf("error adding feed for user %d: %s %s", userId, dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
        return nil, models.ErrNotFound
      } else {
        logrus.Errorf("error adding feed for user %d: %s", userId, err)
        return nil, models.ErrInternalServerError
      }
    }
  }

  return feed, nil
}

func (d *DynamoDBUserRepository) GetFeed(ctx *gin.Context, feedId string, userId string) (*models.Feed, error) { 
  input := &dynamodb.GetItemInput{
		ExpressionAttributeNames: map[string]*string{
			"#fID": aws.String(feedId),
      "#fL": aws.String("feedList"),
		},
		Key: map[string]*dynamodb.AttributeValue{
			"userId": {
				S: aws.String(userId),
			},
		},
    ProjectionExpression: aws.String("#fL.#fID"),
    TableName: aws.String(config.UserTableName),
  }

  feed, err := d.client.GetItem(input)
  if err != nil {
    logrus.Error("error getting feed %d for user %d: %s", feedId, userId, err)
    return nil, models.ErrInternalServerError
  }

	unmarshalled := models.Feed{}
	err = dynamodbattribute.UnmarshalMap(feed.Item, &unmarshalled)
	if err != nil {
    logrus.Errorf("error unmarshalling response payload into feed struct", err)
		return nil, models.ErrInternalServerError
	}

  return &models.Feed{}, nil
}

func (d *DynamoDBUserRepository) RemoveFeed(ctx *gin.Context, feedId string, userId string) (error) {
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeNames: map[string]*string{
			"#fID": aws.String(feedId),
      "#fL": aws.String("feedList"),
		},
		Key: map[string]*dynamodb.AttributeValue{
			"userId": {
				S: aws.String(userId),
			},
		},
		ConditionExpression: aws.String("attribute_exists(#fL.#fID)"),
		UpdateExpression: aws.String("REMOVE #fL.#fID"),
		TableName:        aws.String(config.UserTableName),
	}

  _, err := d.client.UpdateItem(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == dynamodb.ErrCodeResourceNotFoundException {
				logrus.Errorf("error deleting feed %d for user %d: %s %s", feedId, userId, dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
        return models.ErrNotFound
      } else if aerr.Code() == dynamodb.ErrCodeConditionalCheckFailedException {
				logrus.Errorf("error deleting feed %d for user %d: %s %s", feedId, userId, dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
        return models.ErrNotFound
      } else {
        logrus.Errorf("error deleting feed %s for user %d: %s", feedId, userId, err)
        return models.ErrInternalServerError
      }
    }
	}

  return nil
}

func (d *DynamoDBUserRepository) UpdateFeed(ctx *gin.Context, feed *models.Feed) (*models.Feed, error) { 
  return &models.Feed{}, nil 
}

func (d *DynamoDBUserRepository) ListFeedsAll(ctx *gin.Context, userId string) ([]*models.Feed, error) {
  input := &dynamodb.GetItemInput{
    TableName: aws.String(config.UserTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"userId": {
				S: aws.String(userId),
			},
		},
    ProjectionExpression: aws.String("feedList"),
  }

  feedList, err := d.client.GetItem(input)
  if err != nil {
    logrus.Errorf("error getting feed list for user %d: %s", userId, err)
    return nil, models.ErrInternalServerError
  }

	unmarshalled := struct {
    FeedList map[string]*models.Feed     `json:"feedList" dynamodbav:"feedList"`
  }{}
	err = dynamodbattribute.UnmarshalMap(feedList.Item, &unmarshalled)
	if err != nil {
    logrus.Errorf("error unmarshalling response payload into feed list struct", err)
		return nil, models.ErrInternalServerError
	}

  list := valuesFromMap(unmarshalled.FeedList)
  return list, nil
}

func valuesFromMap[M ~map[K]V, K comparable, V any](m M) []V {
    r := make([]V, 0, len(m))
    for _, v := range m {
        r = append(r, v)
    }
    return r
}
