package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/mmcdole/gofeed"
)

type feed struct {
	title      string
	url        string
	timeFormat string
}

type discordChannel struct {
	channel_name string
	server_name  string
	channel_id   int
}

const LAST_CHECKED_TIME = "2022-07-10T00:00:00+10:00"
const LAST_CHECKED_TIME_FORMAT = time.RFC3339

var feedURLS = [...]feed{
	{title: "The Future Does Not Fit In The Containers Of The Past", url: "https://rishad.substack.com/feed", timeFormat: time.RFC1123},
	{title: "Dan Luu", url: "https://danluu.com/atom.xml", timeFormat: time.RFC1123Z},
	{title: "Scattered Thoughts", url: "https://www.scattered-thoughts.net/feed", timeFormat: time.RFC3339},
	{title: "Ben Kuhn", url: "https://www.benkuhn.net/rss", timeFormat: time.RFC3339},
	{title: "Carefree Wandering", url: "https://www.youtube.com/feeds/videos.xml?channel_id=UCnEuIogVV2Mv6Q1a3nHIRsQ", timeFormat: time.RFC3339},
}

var subscribedChannels = [...]discordChannel{
	{channel_name: "mowes mate", server_name: "mines", channel_id: 985831956203851786},
	{channel_name: "pisser", server_name: "klnkn (pers)", channel_id: 1000661720215343114},
}

func commentNewPosts(sess *discordgo.Session, wg *sync.WaitGroup, feed feed) {
	defer wg.Done()
	fp := gofeed.NewParser()

	parsedFeed, _ := fp.ParseURL(feed.url)

	lastChecked, err := time.Parse(LAST_CHECKED_TIME_FORMAT, LAST_CHECKED_TIME)

	if err != nil {
		fmt.Printf("Unable to parse last_checked datetime string: %s", err)
		return
	}

	for _, item := range parsedFeed.Items {
		publishedTime, _ := time.Parse(feed.timeFormat, item.Published)
		if err != nil {
			fmt.Printf("Unable to parse published_time datetime string for post %s in blog %s: %s", item.Title, feed.title, err)
			return
		}

		if publishedTime.After(lastChecked) {
			var message string = fmt.Sprintf("**%s**\n%s\n", item.Title, item.Link)
			for _, channel := range subscribedChannels {
				if _, err := sess.ChannelMessageSend(strconv.Itoa(channel.channel_id), message); err != nil {
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

	var wg sync.WaitGroup
	for _, feedURL := range feedURLS {
		wg.Add(1)
		go commentNewPosts(discord, &wg, feedURL)
	}

	wg.Wait()
}
