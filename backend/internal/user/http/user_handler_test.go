package user

import (
	"github.com/gin-gonic/gin"
	"github.com/mowemcfc/discRSS/models"
)

type mockUserUsecase struct {}

func (*mockUserUsecase) GetUser(ctx *gin.Context, userId string) (*models.UserAccount, error) {return nil, nil}
func (*mockUserUsecase) CreateUser(ctx *gin.Context, user *models.UserAccount) (*models.UserAccount, error) {return nil, nil}
func (*mockUserUsecase) AddFeed(ctx *gin.Context, feed *models.Feed, userId string) (*models.Feed, error) {return nil, nil}
func (*mockUserUsecase) GetFeed(ctx *gin.Context, feedId string, userId string) (*models.Feed, error) {return nil, nil}
func (*mockUserUsecase) UpdateFeed(ctx *gin.Context, feed *models.Feed) (*models.Feed, error) {return nil, nil}
func (*mockUserUsecase) RemoveFeed(ctx *gin.Context, feedId string, userId string) (error) {return nil}
func (*mockUserUsecase) ListFeedsAll(ctx *gin.Context, userId string) ([]*models.Feed, error) {return nil, nil}
