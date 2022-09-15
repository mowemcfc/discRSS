package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
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
