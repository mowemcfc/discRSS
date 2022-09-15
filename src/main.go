package main

import (
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

var UserAccounts map[string]UserAccount

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

const LAST_CHECKED_TIME = "2022-08-30T00:00:00+10:00"
const LAST_CHECKED_TIME_FORMAT = time.RFC3339

func fetchDiscordToken(sess *session.Session) (string, error) {
	scm := secretsmanager.New(sess)

	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String("discRSS/discord-bot-secret"),
	}
	result, err := scm.GetSecretValue(input)
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

	ddb := dynamodb.New(sess)
	var tableName string = "discRSS-UserRecords"

	getUserInput := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"userID": {
				N: aws.String(strconv.Itoa(userID)),
			},
		},
		TableName: aws.String(tableName),
	}

	user, err := ddb.GetItem(getUserInput)
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

func commentNewPosts(sess *discordgo.Session, wg *sync.WaitGroup, feed Feed, channelList []DiscordChannel) {
	defer wg.Done()
	fp := gofeed.NewParser()

	parsedFeed, err := fp.ParseURL(feed.Url)

	if err != nil {
		fmt.Printf("unable to parse URL %s for feed %s: %s", feed.Url, feed.Title, err)
		return
	}

	lastChecked, err := time.Parse(LAST_CHECKED_TIME_FORMAT, LAST_CHECKED_TIME)

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
func start() {

	aws, err := getAWSSession()
	if err != nil {
		fmt.Println(err)
		return
	}

	discordToken, err := fetchDiscordToken(aws)
	if err != nil {
		fmt.Println(err)
		return
	}

	discord, err := getDiscordSession(discordToken)
	if err != nil {
		fmt.Println(err)
		return
	}

	user, err := fetchUser(aws, 1)
	if err != nil {
		fmt.Println(err)
		return
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
}

func main() {
	lambda.Start(start)
}
