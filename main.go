package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
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
const currentUser = "mowemcfc"

var localFeeds = []Feed{
	{FeedID: 1, Title: "The Future Does Not Fit In The Containers Of The Past", Url: "https://rishad.substack.com/feed", TimeFormat: time.RFC1123},
	{FeedID: 2, Title: "Dan Luu", Url: "https://danluu.com/atom.xml", TimeFormat: time.RFC1123Z},
	{FeedID: 3, Title: "Scattered Thoughts", Url: "https://www.scattered-thoughts.net/feed", TimeFormat: time.RFC3339},
	{FeedID: 4, Title: "Ben Kuhn", Url: "https://www.benkuhn.net/rss", TimeFormat: time.RFC3339},
	{FeedID: 5, Title: "Carefree Wandering", Url: "https://www.youtube.com/feeds/videos.xml?channel_id=UCnEuIogVV2Mv6Q1a3nHIRsQ", TimeFormat: time.RFC3339},
}

var localChannels = []DiscordChannel{
	{ChannelID: 985831956203851786, ChannelName: "mowes mate", ServerName: "mines"},
	{ChannelID: 958948046606053406, ChannelName: "rss", ServerName: "klnkn (pers)"},
}

func getDiscordSession() (*discordgo.Session, error) {
	DISCORD_BOT_TOKEN := ""
	if DISCORD_BOT_TOKEN = os.Getenv("DISCORD_BOT_TOKEN"); DISCORD_BOT_TOKEN == "" {
		return nil, fmt.Errorf("error retrieving DISCORD_BOT_TOKEN environment variable. is it set?")
	}

	discord, err := discordgo.New("Bot " + DISCORD_BOT_TOKEN)
	if err != nil {
		return nil, fmt.Errorf("error creating Discord session:\n  %s", err)
	}

	if err = discord.Open(); err != nil {
		return nil, fmt.Errorf("error opening Discord session:\n  %s", err)
	}

	return discord, nil
}

func getDDBSession() (*session.Session, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: "carter-dev",
		Config: aws.Config{
			Region: aws.String("ap-southeast-2"),
		},
	})

	if err != nil {
		return nil, fmt.Errorf("error creating AWS session: \n%s", err)
	}

	return sess, nil
}

func fetchUser(sess *session.Session, userID int) (*dynamodb.GetItemOutput, error) {

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
	return user, nil
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
			}
		}
	}

}

func main() {

	discord, err := getDiscordSession()
	if err != nil {
		fmt.Println(err)
		return
	}

	ddbSession, err := getDDBSession()
	if err != nil {
		fmt.Println(err)
		return
	}

	UserAccounts = make(map[string]UserAccount)
	UserAccounts[currentUser] = UserAccount{
		1,
		currentUser,
		localFeeds,
		localChannels,
	}

	user, err := fetchUser(ddbSession, 1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(user)

	// Initialise a WaitGroup that will spawn a goroutine per subscribed RSS feed to post all new content
	var wg sync.WaitGroup
	for _, feed := range UserAccounts[currentUser].FeedList {
		wg.Add(1)
		go commentNewPosts(discord, &wg, feed, UserAccounts[currentUser].ChannelList)
	}

	wg.Wait()
}
