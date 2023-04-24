package handlers

import (
  "time"
	"strconv"
	"log"
	"net/http"
  "net/url"

	"github.com/mowemcfc/discRSS/models"
	"github.com/mowemcfc/discRSS/internal/response"
	"github.com/mowemcfc/discRSS/internal/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gin-gonic/gin"
)


type AddFeedParams struct {
	Title string
	URL   string
}


func (app *App) AddFeedHandler(c *gin.Context) {
	appG := response.Gin{C: c}

	addFeedParams := AddFeedParams{}

	if err := appG.C.BindJSON(&addFeedParams); err != nil {
		log.Println("error binding addFeed params JSON to addFeedParams struct", err)
    appG.Response(http.StatusBadRequest, interface{}(nil))
		return
	}

  _, err := url.ParseRequestURI(addFeedParams.URL)
  if err != nil {
    log.Printf("error parsing AddFeedHandler request URL %s: %s ", addFeedParams.URL, err)
    appG.Response(http.StatusBadRequest, interface{}(nil))
  }

	requestUserID, err := strconv.Atoi(appG.C.Param("userId"))
	if err != nil {
		log.Println(err)
		return
	}

  if requestUserID < 0 {
    log.Printf("error: request userId was less than 0: %d", requestUserID)
    appG.Response(http.StatusBadRequest, interface{}(nil))
  }

  newFeedId := strconv.FormatInt(time.Now().UnixNano()/(1<<22), 10)
	newFeed := models.Feed{
		FeedID:     newFeedId,
		Title:      addFeedParams.Title,
		Url:        addFeedParams.URL,
		TimeFormat: "z",
	}

	marshalledFeed, err := dynamodbattribute.Marshal(newFeed)
	if err != nil {
		log.Println("error marshalling feed struct into dynamodbattribute map", err)
    appG.Response(http.StatusInternalServerError, interface{}(nil))
		return
	}

	addFeedInput := &dynamodb.UpdateItemInput{
		ExpressionAttributeNames: map[string]*string{
			"#fID": aws.String(newFeedId),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":f": marshalledFeed,
		},
		Key: map[string]*dynamodb.AttributeValue{
			"userId": {
				N: aws.String(strconv.Itoa(requestUserID)),
			},
		},
		ConditionExpression: aws.String("attribute_not_exists(feedList.#fID)"),
		ReturnValues:        aws.String("UPDATED_NEW"),
		UpdateExpression:    aws.String("SET feedList.#fID = :f"),
		TableName:           aws.String(config.UserTableName),
	}

	updatedValues, err := app.DdbSvc.UpdateItem(addFeedInput)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
        log.Printf("error adding feed for user %d: %s %s", requestUserID, dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
        appG.Response(http.StatusInternalServerError, interface{}(nil))
			case dynamodb.ErrCodeResourceNotFoundException:
        log.Printf("error adding feed for user %d: %s %s", requestUserID, dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
        appG.Response(http.StatusNotFound, interface{}(nil))
			case dynamodb.ErrCodeRequestLimitExceeded:
				log.Printf("error adding feed for user %d: %s %s", requestUserID, dynamodb.ErrCodeRequestLimitExceeded, aerr.Error())
        appG.Response(http.StatusInternalServerError, interface{}(nil))
			case dynamodb.ErrCodeInternalServerError:
				log.Printf("error adding feed for user %d: %s %s", requestUserID, dynamodb.ErrCodeInternalServerError, aerr.Error())
        appG.Response(http.StatusInternalServerError, interface{}(nil))
			default:
				log.Printf("error adding feed for user %d: %s", requestUserID, aerr.Error())
			}
		} else {
			log.Printf("error adding feed for user %d: %s", requestUserID, err.Error())
		}
		return
	}
	log.Printf("%v", updatedValues.Attributes)

	appG.Response(http.StatusOK, newFeed)
}

func (app *App) GetFeedHandler(c *gin.Context) {
	appG := response.Gin{C: c}

	requestUserID, err := strconv.Atoi(appG.C.Param("userId"))
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("userId: %d\n", requestUserID)

	feedId := appG.C.Param("feedId")
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("feedId: %s\n", feedId)

	user, err := app.FetchUser(requestUserID)
	if err != nil {
		log.Println("error fetching user from DDB", err)
		return
	}


  feed, found := user.FeedList[feedId]
  if (!found) {
    appG.Response(http.StatusNotFound, "Unable to find feed")
    return
  }

	appG.Response(http.StatusOK, feed)
}

func (app *App) DeleteFeedHandler(c *gin.Context) {
	appG := response.Gin{C: c}

	requestUserId := appG.C.Param("userId")
	requestFeedId := appG.C.Param("feedId")
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeNames: map[string]*string{
			"#fID": aws.String(requestFeedId),
			"#fL":  aws.String("feedList"),
		},
		Key: map[string]*dynamodb.AttributeValue{
			"userId": {
				N: aws.String(requestUserId),
			},
		},
		// todo(mowemcfc): actually figure out how to add a working ConditionExpression
		//ConditionExpression: aws.String("attribute_exists (#fL)"),
		UpdateExpression: aws.String("REMOVE #fL.#fID"),
		TableName:        aws.String(config.UserTableName),
	}

	_, err := app.DdbSvc.UpdateItem(input)
	if err != nil {
		log.Println("error deleting feed row:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	appG.Response(http.StatusNoContent, interface{}(nil))
}

