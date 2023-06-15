package user

import (
  "net/http"
	"github.com/sirupsen/logrus"
	"github.com/gin-gonic/gin"

  "github.com/mowemcfc/discRSS/models"
  "github.com/mowemcfc/discRSS/internal/user/usecase"
  "github.com/mowemcfc/discRSS/internal/response"
)

type UserHandler struct {
  userUsecase usecase.UserUsecase
}


func NewUserHandler(g *gin.Engine, usecase usecase.UserUsecase) UserHandler {
  handler := UserHandler{ userUsecase: usecase }

	userRoute := g.Group("/user")
	{
		userRoute.GET("/:userId", handler.GetUser)
		userRoute.POST("/:userId", handler.CreateUser)
		userRoute.DELETE("/:userId", handler.DeleteUser)

		userRoute.GET("/:userId/feed/:feedId", handler.GetFeed)
		userRoute.POST("/:userId/feed", handler.AddFeed)
		userRoute.DELETE("/:userId/feed/:feedId", handler.RemoveFeed)
	}

  return handler
}

func (handler *UserHandler) GetUser(c *gin.Context) {
	appG := response.Gin{C: c}

  user, err := handler.userUsecase.GetUser(c)
  if err != nil {
    code := getStatusCode(err)
    appG.Response(code, err)
  }

	appG.Response(http.StatusOK, user)
}
func (handler *UserHandler) CreateUser(c *gin.Context) {
	appG := response.Gin{C: c}
	appG.Response(http.StatusNotImplemented, "Method PUT for resource /user not implemented")
}
func (handler *UserHandler) DeleteUser(c *gin.Context) {
	appG := response.Gin{C: c}
  appG.Response(http.StatusNotImplemented, "Method DELETE for resource /user not implemented")
}
func (handler *UserHandler) AddFeed(c *gin.Context) {
	appG := response.Gin{C: c}
	appG.Response(http.StatusNotImplemented, "Method PUT for resource /user/:userId/feed not implemented")
}
func (handler *UserHandler) GetFeed(c *gin.Context) {
	appG := response.Gin{C: c}
	appG.Response(http.StatusNotImplemented, "Method GET for resource /user/:userId/feed not implemented")
}
func (handler *UserHandler) UpdateFeed(c *gin.Context) {
	appG := response.Gin{C: c}
  appG.Response(http.StatusNotImplemented, "Method PATCH for resource /user/:userId/feed not implemented")
}
func (handler *UserHandler) RemoveFeed(c *gin.Context) {
	appG := response.Gin{C: c}
	appG.Response(http.StatusNotImplemented, "Method DELETE for resource /user/:userId/feed not implemented")
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
