package dynamodb

import (
  //"log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type MockDynamoDB struct {
  dynamodbiface.DynamoDBAPI
  store map[string]map[string]*dynamodb.AttributeValue
}

func NewMockDynamoDB() *MockDynamoDB {
  user1 := map[string]*dynamodb.AttributeValue{
		"userId":   {N: aws.String("0")},
		"username": {S: aws.String("John Doe")},
		"feedList": {
      M: map[string]*dynamodb.AttributeValue{
        "0": { 
          M: map[string]*dynamodb.AttributeValue{
            "feedId": {N: aws.String("123")},
            "title": {S: aws.String("feed1")},
            "url": {S: aws.String("https://feed1.com/rss")},
            "timeFormat": {S: aws.String("Mon, 02 Jan 2006 15:04:05 MST")},
          },
        },
      },
    },
		"channelList": {
      M: map[string]*dynamodb.AttributeValue{
        "0": {
          M: map[string]*dynamodb.AttributeValue{
            "channelId": {N: aws.String("123")},
            "channelName": {S: aws.String("channel1")},
            "serverName": {S: aws.String("server1")},
          },
        },
      },
    },
	}
  return &MockDynamoDB{
    store: map[string]map[string]*dynamodb.AttributeValue{
      "0": user1,
    },
  }
}

func (m *MockDynamoDB) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
  return &dynamodb.PutItemOutput{}, nil
}

func (m *MockDynamoDB) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
  _, found := m.store[*input.Key["userId"].N]
  if !found {
    return nil, awserr.New(dynamodb.ErrCodeResourceNotFoundException, "unable to find userId", nil)
  }

  return &dynamodb.GetItemOutput{
    Item: m.store[*input.Key["userId"].N],
  }, nil
}

func (m *MockDynamoDB) UpdateItem(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
  _, found := m.store[*input.Key["userId"].N]
  if !found {
    return nil, awserr.New(dynamodb.ErrCodeResourceNotFoundException, "unable to find userId", nil)
  }

  return &dynamodb.UpdateItemOutput{}, nil
}

