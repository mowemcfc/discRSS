package usecase

import (
  "context"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"

	"github.com/mowemcfc/discRSS/internal/config"
	"github.com/mowemcfc/discRSS/internal/user/repository/dynamodb"
	"github.com/mowemcfc/discRSS/models"
)

type userUsecase struct {
  userRepo dynamodb.UserRepository
}

type UserUsecase interface {
  GetUser(ctx context.Context, userId string) (*models.UserAccount, error)
  CreateUser(ctx context.Context, user *models.UserAccount) (*models.UserAccount, error)
  AddFeed(ctx context.Context, feed *models.Feed, userId string) (*models.Feed, error)
  GetFeed(ctx context.Context, feedId string, userId string) (*models.Feed, error)
  UpdateFeed(ctx context.Context, feed *models.Feed) (*models.Feed, error)
  RemoveFeed(ctx context.Context, feedId string, userId string) (error)
  ListFeedsAll(ctx context.Context, userId string) ([]*models.Feed, error)
}

func NewUserUsecase(userRepo dynamodb.UserRepository) UserUsecase {
  return &userUsecase{ userRepo: userRepo }
}

func (usecase *userUsecase) GetUser(ctx context.Context, userId string) (*models.UserAccount, error) { 
  tracer := otel.GetTracerProvider().Tracer(config.AppName)
  ctx, span := tracer.Start(ctx, "user_usecase")
  defer span.End()
  user, err := usecase.userRepo.GetUser(ctx, userId)
  if err != nil {
    logrus.Errorf("usecase GetUser error: %s", err)
    return nil, err
  }

  return user, nil
}

func (usecase *userUsecase) CreateUser(ctx context.Context, u *models.UserAccount) (*models.UserAccount, error) {
  user, err := usecase.userRepo.CreateUser(ctx, u)
  if err != nil {
    logrus.Errorf("usecase CreateUser error: %s", err)
    return nil, err
  }

  return user, nil
}

func (usecase *userUsecase) AddFeed(ctx context.Context, f *models.Feed, userId string) (*models.Feed, error) {
  feed, err := usecase.userRepo.AddFeed(ctx, f, userId)
  if err != nil {
    logrus.Errorf("usecase AddFeed error: %s", err)
    return nil, err
  }

  return feed, nil
}

func (usecase *userUsecase) GetFeed(ctx context.Context, feedId string, userId string) (*models.Feed, error) {
  feed, err := usecase.userRepo.GetFeed(ctx, feedId, userId)
  if err != nil {
    logrus.Errorf("usecase GetFeed error: %s", err)
    return nil, err
  }

  return feed, nil
}

func (usecase *userUsecase) UpdateFeed(ctx context.Context, f *models.Feed) (*models.Feed, error) {
  feed, err := usecase.userRepo.UpdateFeed(ctx, f)
  if err != nil {
    logrus.Errorf("usecase UpdateFeed error: %s", err)
    return nil, err
  }

  return feed, nil
}

func (usecase *userUsecase) RemoveFeed(ctx context.Context, feedId string, userId string) (error) {
  err := usecase.userRepo.RemoveFeed(ctx, feedId, userId)
  if err != nil {
    logrus.Errorf("usecase RemoveFeed error: %s", err)
    return err
  }

  return nil
}

func (usecase *userUsecase) ListFeedsAll(ctx context.Context, userId string) ([]*models.Feed, error) {
  feeds, err := usecase.userRepo.ListFeedsAll(ctx, userId)
  if err != nil {
    logrus.Errorf("usecase UpdateFeed error: %s", err)
    return nil, err
  }

  return feeds, nil
}
