package main

import (
	"fmt"
	"os"
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

const CHANNEL_ID = "985831956203851786"
const LAST_CHECKED_TIME = "2022-06-01T00:00:00+10:00"
const LAST_CHECKED_TIME_FORMAT = time.RFC3339

var feedURLS = [...]feed{
	{title: "The Future Does Not Fit In The Containers Of The Past", url: "https://rishad.substack.com/feed", timeFormat: time.RFC1123},
	{title: "Dan Luu", url: "https://danluu.com/atom.xml", timeFormat: time.RFC1123Z},
	{title: "Scattered Thoughts", url: "https://www.scattered-thoughts.net/feed", timeFormat: time.RFC3339},
	{title: "Ben Kuhn", url: "https://www.benkuhn.net/rss", timeFormat: time.RFC3339},
	{title: "Carefree Wandering", url: "https://www.youtube.com/feeds/videos.xml?channel_id=UCnEuIogVV2Mv6Q1a3nHIRsQ", timeFormat: time.RFC3339},
}

func commentNewPosts(sess *discordgo.Session, wg *sync.WaitGroup, feed feed) {
	defer wg.Done()
	fp := gofeed.NewParser()

	resp, _ := fp.ParseURL(feed.url)

	lastChecked, _ := time.Parse(LAST_CHECKED_TIME_FORMAT, LAST_CHECKED_TIME)
	fmt.Println(lastChecked)

	fmt.Printf("feed: %s\n", feed.title)
	for _, item := range resp.Items {
		publishedTime, _ := time.Parse(feed.timeFormat, item.Published)
		if publishedTime.After(lastChecked) {
			var message = &discordgo.MessageSend{
				Content: fmt.Sprintf("**%s**\n\n%s", item.Title, item.Link),
			}

			if _, err := sess.ChannelMessageSendComplex(CHANNEL_ID, message); err != nil {
				fmt.Printf("Error sending message: %s", err)
				return
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
