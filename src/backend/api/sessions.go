package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func getAWSSession() (*session.Session, error) {
	var sess *session.Session
	var err error
	if isLocal {
		sess, err = session.NewSessionWithOptions(session.Options{
			Config:  aws.Config{Region: aws.String(os.Getenv("AWS_LOCAL_REGION"))},
			Profile: os.Getenv("AWS_LOCAL_NAMED_PROFILE"),
		})
	} else {
		sess, err = session.NewSession()
	}

	if err != nil {
		return nil, fmt.Errorf("error creating AWS session: \n%s", err)
	}

	fmt.Println("successfully opened AWS session")

	return sess, nil
}
