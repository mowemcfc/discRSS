package usecase

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/mowemcfc/discRSS/internal/user/repository/dynamodb"
	"github.com/mowemcfc/discRSS/models"
)

type userUsecase struct {
  userRepo dynamodb.UserRepository
}

type UserUsecase interface {
  GetUser(ctx *gin.Context) (*models.UserAccount, error)
  CreateUser(ctx *gin.Context) (*models.UserAccount, error)
  AddFeed(ctx *gin.Context) (*models.Feed, error)
  UpdateFeed(ctx *gin.Context) (*models.Feed, error)
  RemoveFeed(ctx *gin.Context) (error)
}

func NewUserUsecase(userRepo dynamodb.UserRepository) UserUsecase {
  return &userUsecase{ userRepo: userRepo }
}

func (usecase *userUsecase) GetUser(ctx *gin.Context) (*models.UserAccount, error) { 
  user, err := usecase.userRepo.GetUser(ctx)
  if err != nil {
    logrus.Errorf("usecase GetUser error: %s", err)
    return nil, err
  }

  return user, nil
}

func (usecase *userUsecase) CreateUser(ctx *gin.Context) (*models.UserAccount, error) {
  user, err := usecase.userRepo.CreateUser(ctx)
  if err != nil {
    logrus.Errorf("usecase CreateUser error: %s", err)
    return nil, err
  }

  return user, nil
}

func (usecase *userUsecase) AddFeed(ctx *gin.Context) (*models.Feed, error) {
  user, err := usecase.userRepo.AddFeed(ctx)
  if err != nil {
    logrus.Errorf("usecase AddFeed error: %s", err)
    return nil, err
  }

  return user, nil
}

func (usecase *userUsecase) UpdateFeed(ctx *gin.Context) (*models.Feed, error) {
  user, err := usecase.userRepo.UpdateFeed(ctx)
  if err != nil {
    logrus.Errorf("usecase UpdateFeed error: %s", err)
    return nil, err
  }

  return user, nil
}

func (usecase *userUsecase) RemoveFeed(ctx *gin.Context) (error) {
  err := usecase.userRepo.RemoveFeed(ctx)
  if err != nil {
    logrus.Errorf("usecase UpdateFeed error: %s", err)
    return err
  }

  return nil
}
