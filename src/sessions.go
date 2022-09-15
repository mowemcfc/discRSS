package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/bwmarrin/discordgo"
)

func getDiscordSession(token string) (*discordgo.Session, error) {
	if token == "" {
		return nil, fmt.Errorf("discord token is empty. has it been properly retrieved from AWS?")
	}

	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("error creating Discord session:\n  %s", err)
	}

	if err = discord.Open(); err != nil {
		return nil, fmt.Errorf("error opening Discord session:\n  %s", err)
	}

	fmt.Println("successfully opened discord session")

	return discord, nil
}

func getAWSSession() (*session.Session, error) {
	sess, err := session.NewSession()

	if err != nil {
		return nil, fmt.Errorf("error creating AWS session: \n%s", err)
	}

	fmt.Println("successfully opened AWS session")

	return sess, nil
}
