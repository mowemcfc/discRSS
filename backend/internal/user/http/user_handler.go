package user

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"go.opentelemetry.io/otel"

	"github.com/mowemcfc/discRSS/internal/auth0"
	"github.com/mowemcfc/discRSS/internal/config"
	"github.com/mowemcfc/discRSS/internal/response"
	"github.com/mowemcfc/discRSS/internal/user/usecase"
	"github.com/mowemcfc/discRSS/models"
)

type UserHandler struct {
  userUsecase usecase.UserUsecase
}


func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	logrus.Error(err)
	switch err {
	case models.ErrInternalServerError:
		return http.StatusInternalServerError
	case models.ErrNotFound:
		return http.StatusNotFound
	case models.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

func NewUserHandler(g *gin.Engine, usecase usecase.UserUsecase) UserHandler {
  handler := UserHandler{ userUsecase: usecase }

	g.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"POST", "PATCH", "PUT", "DELETE", "GET", "OPTIONS"},
		AllowHeaders:     []string{"*", "Authorization"},
		AllowCredentials: true,
	}))
  pprof.Register(g, "dev/pprof")

  g.Use(auth0.EnsureValidToken())
  g.Use(auth0.EnsureValidClaims())
  g.GET("/user/:userId", handler.GetUser)
  g.POST("/user", handler.CreateUser)
  g.DELETE("/user/:userId", handler.DeleteUser)

  g.GET("/user/:userId/feeds", handler.ListFeedsAll)
  g.GET("/user/:userId/feed/:feedId", handler.GetFeed)
  g.POST("/user/:userId/feed", handler.AddFeed)
  g.DELETE("/user/:userId/feed/:feedId", handler.RemoveFeed)

  return handler
}

func (handler *UserHandler) GetUser(c *gin.Context) {
	appG := response.Gin{C: c}
  tracer := otel.GetTracerProvider().Tracer(config.AppName)
  ctx, span := tracer.Start(c.Request.Context(), "user_handler")
  defer span.End()
  userId := appG.C.Param("userId")
  res, err := handler.userUsecase.GetUser(ctx, userId)
  if err != nil {
    code := getStatusCode(err)
    appG.Response(code, err.Error())
    return
  }

	appG.Response(http.StatusOK, res)
}

type CreateUserParams struct {
	UserID      string                     `json:"userId" dynamodbav:"userId"`
	Username    string                     `json:"username" dynamodbav:"username"`
}
func (handler *UserHandler) CreateUser(c *gin.Context) {
	appG := response.Gin{C: c}

	var createUserParams CreateUserParams
	if err := c.BindJSON(&createUserParams); err != nil {
    logrus.Error("error binding user params JSON to models.UserAccount struct: ", err)
    appG.Response(http.StatusBadRequest, "bad request format")
    return
	}

  user := &models.UserAccount{
    UserID: createUserParams.UserID,
    Username: createUserParams.Username,
    FeedList: make(map[string]*models.Feed),
    ChannelList: make(map[string]*models.DiscordChannel),
  }
  res, err := handler.userUsecase.CreateUser(c.Request.Context(), user)
  if err != nil {
    code := getStatusCode(err)
    appG.Response(code, err.Error())
    return
  }

	appG.Response(http.StatusOK, res)
}

func (handler *UserHandler) DeleteUser(c *gin.Context) {
	appG := response.Gin{C: c}
  appG.Response(http.StatusNotImplemented, "Method DELETE for resource /user not implemented")
}

type AddFeedParams struct {
	Title      string `json:"title" dynamodbav:"title"`
	URL        string `json:"url" dynamodbav:"url"`
}
func (handler *UserHandler) AddFeed(c *gin.Context) {
	appG := response.Gin{C: c}

  var addFeedParams AddFeedParams

	if err := appG.C.BindJSON(&addFeedParams); err != nil {
    logrus.Error("error binding addFeed params JSON to addFeedParams struct: ", err)
		appG.Response(http.StatusBadRequest, "request data is not structured correctly")
		return
	}

	_, err := url.ParseRequestURI(addFeedParams.URL)
	if err != nil {
		logrus.Errorf("error parsing AddFeedHandler request URL %s: %s ", addFeedParams.URL, err)
		appG.Response(http.StatusBadRequest, "given feed URL is not a valid URI")
		return
	}

	userIdS := appG.C.Param("userId")
  userId, err := strconv.Atoi(userIdS)
	if err != nil {
    logrus.Error("error parsing request userId: ", err)
    appG.Response(http.StatusBadRequest, "error parsing request userId")
		return
	}

	if userId < 0 {
		logrus.Errorf("error: request userId was less than 0: %d", userId)
		appG.Response(http.StatusBadRequest, "request userId is less than 0")
		return
	}

	newFeedId := strconv.FormatInt(time.Now().UnixNano()/(1<<22), 10)
	newFeed := &models.Feed{
		FeedID:     newFeedId,
		Title:      addFeedParams.Title,
		Url:        addFeedParams.URL,
		TimeFormat: "z",
	}
  res, err := handler.userUsecase.AddFeed(c.Request.Context(), newFeed, userIdS)
  if err != nil {
    code := getStatusCode(err)
    appG.Response(code, err.Error())
    return
  }

	appG.Response(http.StatusOK, res)
}

func (handler *UserHandler) GetFeed(c *gin.Context) {
	appG := response.Gin{C: c}

	userIdS := appG.C.Param("userId")
  userId, err := strconv.Atoi(userIdS)
	if err != nil {
    logrus.Error("error parsing request userId: ", err)
    appG.Response(http.StatusBadRequest, "error parsing request userId")
		return
	}

	if userId < 0 {
		logrus.Error("error: userId value is less than 0")
		appG.Response(http.StatusBadRequest, "userId value is less than 0")
		return
	}

  feedIdS := appG.C.Param("feedId")
  feedId, err := strconv.Atoi(feedIdS)
	if err != nil {
    logrus.Println("error parsing feedId: ", err)
		appG.Response(http.StatusBadRequest, "error parsing feedId parameter")
		return
	}

	if feedId < 0 {
		logrus.Error("error: feedId value is less than 0")
		appG.Response(http.StatusBadRequest, "userId value is less than 0")
		return
	}

  res, err := handler.userUsecase.GetFeed(c.Request.Context(), feedIdS, userIdS)
  if err != nil {
    code := getStatusCode(err)
    appG.Response(code, err.Error())
    return
  }

  appG.Response(http.StatusOK, res)
}

func (handler *UserHandler) UpdateFeed(c *gin.Context) {
	appG := response.Gin{C: c}
  appG.Response(http.StatusNotImplemented, "Method PATCH for resource /user/:userId/feed not implemented")
}

func (handler *UserHandler) RemoveFeed(c *gin.Context) {
	appG := response.Gin{C: c}

	userIdS := appG.C.Param("userId")
	feedIdS := appG.C.Param("feedId")
  feedId, err := strconv.Atoi(feedIdS)
  if err != nil {
    logrus.Error("error converting request feed ID to int: ", err)
    appG.Response(http.StatusBadRequest, "error parsing feedId param")
    return
  }
  if feedId < 0 {
    logrus.Error("error: feedId value was less than 0: ", err)
    appG.Response(http.StatusBadRequest, "feedId value is less than 0")
    return
  }

  err = handler.userUsecase.RemoveFeed(c.Request.Context(), feedIdS, userIdS)
  if err != nil {
    code := getStatusCode(err)
    appG.Response(code, err.Error())
    return
  }

  appG.Response(http.StatusNoContent, interface{}(nil))
}

func (handler *UserHandler) ListFeedsAll(c *gin.Context) {
	appG := response.Gin{C: c}

	userIdS := appG.C.Param("userId")

  res, err := handler.userUsecase.ListFeedsAll(c.Request.Context(), userIdS)
  if err != nil {
    code := getStatusCode(err)
    appG.Response(code, err.Error())
    return
  }

  appG.Response(http.StatusOK, res)
}
