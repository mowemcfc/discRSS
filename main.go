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

func fetchLocalFeeds() any {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Profile: "carter-dev",
		Config: aws.Config{
			Region: aws.String("ap-southeast-2"),
		},
	}))

	ddb := dynamodb.New(sess)
	var tableName string = "discRSS-UserRecords"

	getUserInput := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"userID": {
				N: aws.String("1"),
			},
		},
		TableName: aws.String(tableName),
	}

	user, err := ddb.GetItem(getUserInput)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeRequestLimitExceeded:
				fmt.Println(dynamodb.ErrCodeRequestLimitExceeded, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
	}
	return user
}

func commentNewPosts(sess *discordgo.Session, wg *sync.WaitGroup, feed Feed, channelList []DiscordChannel) {
	defer wg.Done()
	fp := gofeed.NewParser()

	parsedFeed, err := fp.ParseURL(feed.Url)

	if err != nil {
		fmt.Printf("Unable to parse URL %s for feed %s: %s", feed.Url, feed.Title, err)
		return
	}

	lastChecked, err := time.Parse(LAST_CHECKED_TIME_FORMAT, LAST_CHECKED_TIME)

	if err != nil {
		fmt.Printf("Unable to parse last_checked datetime string: %s", err)
		return
	}

	for _, item := range parsedFeed.Items {
		publishedTime, err := time.Parse(feed.TimeFormat, item.Published)
		if err != nil {
			fmt.Printf("Unable to parse published_time datetime string for post %s in blog %s: %s", item.Title, feed.Title, err)
			return
		}

		if publishedTime.After(lastChecked) {
			var message string = fmt.Sprintf("**%s**\n%s\n", item.Title, item.Link)
			for _, channel := range channelList {
				if _, err := sess.ChannelMessageSend(strconv.Itoa(channel.ChannelID), message); err != nil {
					fmt.Printf("Error sending message: %s", err)
					return
				}
			}
		}
	}

}

func main() {

	discord, err := discordgo.New("Bot " + os.Getenv("DISCORD_BOT_TOKEN"))

	if err != nil {
		fmt.Printf("Error creating Discord session:\n  %s", err)
		return
	}

	if err = discord.Open(); err != nil {
		fmt.Printf("Error opening discord session\n  %s", err)
		return
	}

	UserAccounts = make(map[string]UserAccount)
	UserAccounts[currentUser] = UserAccount{
		1,
		currentUser,
		localFeeds,
		localChannels,
	}

	fetchLocalFeeds()

	// Initialise a WaitGroup that will spawn a goroutine per subscribed RSS feed to post all new content
	var wg sync.WaitGroup
	for _, feed := range UserAccounts[currentUser].FeedList {
		wg.Add(1)
		go commentNewPosts(discord, &wg, feed, UserAccounts[currentUser].ChannelList)
	}

	wg.Wait()
}
