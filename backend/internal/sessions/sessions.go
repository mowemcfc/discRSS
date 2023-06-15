package sessions

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/bwmarrin/discordgo"
)

func GetAWSSession(isLocal bool) (*session.Session, error) {
  sess := session.Must(session.NewSession())
  if err := validateAWSSession(sess); err != nil {
    panic(err)
  }

	log.Println("Opened AWS session")

	return sess, nil
}

// Check if the given AWS Session's credentials are invalid or expired
func validateAWSSession(session *session.Session) error {
  stsSvc := sts.New(session)
  input := &sts.GetCallerIdentityInput{}
  _, err := stsSvc.GetCallerIdentity(input)
  if err != nil {
    return err
  }
  return nil
}

func GetDiscordSession(token string) (*discordgo.Session, error) {
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

	log.Println("Opened Discord session")

	return discord, nil
}
