package dynamodb

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type DynamoDBClient struct {
  client dynamodbiface.DynamoDBAPI
}

func (d *DynamoDBClient) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
    return d.client.PutItem(input)
}

func (d *DynamoDBClient) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
    return d.client.GetItem(input)
}

func (d *DynamoDBClient) UpdateItem(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
    return d.client.UpdateItem(input)
}
