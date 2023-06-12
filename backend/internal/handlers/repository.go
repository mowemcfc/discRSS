package handlers

import (
  "github.com/mowemcfc/discRSS/models"
)

type UserRepository interface {
  GetUser() (*models.UserAccount, error)
  CreateUser() (error)
  AddFeed() (*models.Feed, error)
  UpdateFeed(*models.Feed, error)
  RemoveFeed(error)
}
