package usecase

import (
	"context"

	"github.com/mowemcfc/discRSS/internal/user/repository/dynamodb"
	"github.com/mowemcfc/discRSS/models"
)

type userUsecase struct {
  userRepo dynamodb.DynamoDBUserRepository
}

func (usecase *userUsecase) GetUser(ctx context.Context) (*models.UserAccount, error) { 
  return &models.UserAccount{}, nil
}
