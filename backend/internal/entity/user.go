package entity

type Feed struct {
	FeedID     string `json:"feedId" dynamodbav:"feedId"`
	Title      string `json:"title" dynamodbav:"title"`
	Url        string `json:"url" dynamodbav:"url"`
	TimeFormat string `json:"timeFormat" dynamodbav:"timeFormat"`
}

type UserAccount struct {
	UserID      string                     `json:"userId" dynamodbav:"userId"`
	Username    string                     `json:"username" dynamodbav:"username"`
	FeedList    map[string]*Feed           `json:"feedList" dynamodbav:"feedList"`
	ChannelList map[string]*DiscordChannel `json:"channelList" dynamodbav:"channelList"`
}

type DiscordChannel struct {
	ChannelName string `json:"channelName" dynamodbav:"channelName"`
	ServerName  string `json:"serverName" dynamodbav:"serverName"`
	ChannelID   int    `json:"channelID" dynamodbav:"channelID"`
}

type AppConfig struct {
	AppName               string `json:"appName" dynamodbav:"appName"`
	LastCheckedTime       string `json:"lastCheckedTime" dynamodbav:"lastCheckedTime"`
	LastCheckedTimeFormat string `json:"lastCheckedTimeFormat" dynamodbav:"lastCheckedTimeFormat"`
}

