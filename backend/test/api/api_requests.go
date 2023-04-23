package test

import (
  "errors"
  "testing"
  "github.com/stretchr/testify/assert"

  "github.com/mowemcfc/discRSS/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type mockDynamoDBSvc struct {
	dynamodbiface.DynamoDBAPI
	mockGetItem func(*dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
}

func (m *mockDynamoDBSvc) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	return m.mockGetItem(input)
}

func TestFetchUser(t *testing.T) {
	// Test data
	testUserID := "1"
	testUser := &models.UserAccount{
		UserID:   testUserID,
		Username: "testUser",
    FeedList: make(map[string]*models.Feed),
    ChannelList: make(map[string]*models.DiscordChannel),
	}

	// Test successful fetch
	{
		ddbSvc := &mockDynamoDBSvc{
			mockGetItem: func(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
				item, _ := dynamodbattribute.MarshalMap(testUser)
				return &dynamodb.GetItemOutput{Item: item}, nil
			},
		}

		user, err := FetchUser(testUserID, ddbSvc)
		assert.NoError(t, err)
		assert.Equal(t, testUser, user)
	}

	// Test DynamoDB error
	{
		ddbSvc := &mockDynamoDBSvc{
			mockGetItem: func(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
				return nil, errors.New("DynamoDB error")
			},
		}

		user, err := FetchUser(testUserID, ddbSvc)
		assert.Error(t, err)
		assert.Nil(t, user)
	}

	// Test unmarshalling error
	{
		ddbSvc := &mockDynamoDBSvc{
			mockGetItem: func(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
				return &dynamodb.GetItemOutput{
					Item: map[string]*dynamodb.AttributeValue{
						"userId": {
							N: aws.String("notAnInt"),
						},
					},
				}, nil
			},
		}

		user, err := FetchUser(testUserID, ddbSvc)
		assert.Error(t, err)
		assert.Nil(t, user)
	}
}
