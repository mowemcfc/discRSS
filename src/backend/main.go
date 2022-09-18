package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/secretsmanager"

	"github.com/bwmarrin/discordgo"
	"github.com/mmcdole/gofeed"
)

type Feed struct {
	FeedID     int    `json:"feedID"`
	Title      string `json:"title"`
	Url        string `json:"url"`
	TimeFormat string `json:"TimeFormat"`
}

type UserAccount struct {
	UserID      int              `json:"userID"`
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

var secretsmanagerSvc *secretsmanager.SecretsManager
var ddbSvc *dynamodb.DynamoDB

const APP_NAME string = "discRSS"
const APP_CONFIG_TABLE_NAME string = "discRSS-AppConfigs"
const USER_TABLE_NAME string = "discRSS-UserRecords"
const BOT_TOKEN_SECRET_NAME string = "discRSS/discord-bot-secret"

func fetchAppConfig(sess *session.Session, appName string) (*AppConfig, error) {

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

func fetchDiscordToken(sess *session.Session) (string, error) {

	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(BOT_TOKEN_SECRET_NAME),
	}

	result, err := secretsmanagerSvc.GetSecretValue(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeResourceNotFoundException:
				return "", fmt.Errorf("%s %s", secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			case secretsmanager.ErrCodeInvalidParameterException:
				return "", fmt.Errorf("%s %s", secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
			case secretsmanager.ErrCodeInvalidRequestException:
				return "", fmt.Errorf("%s %s", secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
			case secretsmanager.ErrCodeDecryptionFailure:
				return "", fmt.Errorf("%s %s", secretsmanager.ErrCodeDecryptionFailure, aerr.Error())
			case secretsmanager.ErrCodeInternalServiceError:
				return "", fmt.Errorf("%s %s", secretsmanager.ErrCodeInternalServiceError, aerr.Error())
			default:
				return "", fmt.Errorf("%s", aerr.Error())
			}
		} else {
			return "", fmt.Errorf(fmt.Sprintln(err.Error()))
		}
	}

	fmt.Println("successfully fetched discord token")

	return *result.SecretString, nil
}

func fetchUser(sess *session.Session, userID int) (*UserAccount, error) {

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

func updateLastCheckedTime(sess *session.Session, t time.Time) error {
	formatted := t.Format(discRssConfig.LastCheckedTimeFormat)

	updateTimeInput := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":t": {
				S: aws.String(formatted),
			},
		},
		Key: map[string]*dynamodb.AttributeValue{
			"appName": {
				S: aws.String(discRssConfig.AppName),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set lastCheckedTime = :t"),
		TableName:        aws.String(APP_CONFIG_TABLE_NAME),
	}

	_, err := ddbSvc.UpdateItem(updateTimeInput)
	if err != nil {
		return fmt.Errorf("error updating last checked time: %s", err)
	}

	fmt.Printf("successfully updated last checked time: %s\n", formatted)

	return nil
}

func commentNewPosts(sess *discordgo.Session, wg *sync.WaitGroup, feed Feed, channelList []DiscordChannel) {
	defer wg.Done()
	fp := gofeed.NewParser()

	parsedFeed, err := fp.ParseURL(feed.Url)

	if err != nil {
		fmt.Printf("unable to parse URL %s for feed %s: %s", feed.Url, feed.Title, err)
		return
	}

	lastChecked, err := time.Parse(discRssConfig.LastCheckedTimeFormat, discRssConfig.LastCheckedTime)

	if err != nil {
		fmt.Printf("unable to parse last_checked datetime string: %s", err)
		return
	}

	for _, item := range parsedFeed.Items {
		publishedTime, err := time.Parse(feed.TimeFormat, item.Published)
		if err != nil {
			fmt.Printf("unable to parse published_time datetime string for post %s in blog %s: %s", item.Title, feed.Title, err)
			return
		}

		if publishedTime.After(lastChecked) {
			var message string = fmt.Sprintf("**%s**\n%s\n", item.Title, item.Link)
			for _, channel := range channelList {
				if _, err := sess.ChannelMessageSend(strconv.Itoa(channel.ChannelID), message); err != nil {
					fmt.Printf("error sending message: %s", err)
					return
				}
				fmt.Printf("successfully sent message: %s to channel: %s %d\n", message, channel.ChannelName, channel.ChannelID)
			}
		}
	}
}

//func start(ctx context.Context) {
func start() (map[string]interface{}, error) {

	aws, err := getAWSSession()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	secretsmanagerSvc = secretsmanager.New(aws)
	ddbSvc = dynamodb.New(aws)

	discRssConfig, err = fetchAppConfig(aws, "discRSS")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	discordToken, err := fetchDiscordToken(aws)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	discord, err := getDiscordSession(discordToken)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	user, err := fetchUser(aws, 1)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Printf("user %d channels: %+v", user.UserID, user.ChannelList)
	fmt.Printf("user %d feeds: %+v", user.UserID, user.FeedList)

	// Initialise a WaitGroup that will spawn a goroutine per subscribed RSS feed to post all new content
	var wg sync.WaitGroup
	for _, feed := range user.FeedList {
		wg.Add(1)
		go commentNewPosts(discord, &wg, feed, user.ChannelList)
	}
	wg.Wait()

	if err := updateLastCheckedTime(aws, time.Now()); err != nil {
		fmt.Println(err)
		return nil, err
	}

	marshalled, err := json.Marshal(user)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Printf("returning: %+v", user)
	return map[string]interface{}{
			"statusCode":      200,
			"isBase64Encoded": false,
			"headers":         map[string]string{"Content-Type": "application/json"},
			"body":            string(marshalled),
		},
		nil
}

func main() {
	lambda.Start(start)
}
