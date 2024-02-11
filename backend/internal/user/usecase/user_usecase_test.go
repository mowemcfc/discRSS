package usecase

import (
  "reflect"
  "context"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mowemcfc/discRSS/models"
)

type mockUserRepo struct {
  getUser func(context.Context, string) (*models.UserAccount, error)
  createUser func(context.Context, *models.UserAccount) (*models.UserAccount, error)
  addFeed func(context.Context, *models.Feed, string) (*models.Feed, error)
  getFeed func(context.Context, string, string) (*models.Feed, error)
  updateFeed func(context.Context, *models.Feed) (*models.Feed, error)
  removeFeed func(context.Context, string, string) (error)
  listFeedsAll func(context.Context, string) ([]*models.Feed, error)
}

func (m *mockUserRepo) GetUser(c context.Context, userId string) (*models.UserAccount, error) { return m.getUser(c, userId) }
func (m *mockUserRepo) CreateUser(c context.Context, user *models.UserAccount) (*models.UserAccount, error) { return m.createUser(c, user) }
func (m *mockUserRepo) AddFeed(c context.Context, f *models.Feed, userId string) (*models.Feed, error) { return m.addFeed(c, f, userId) }
func (m *mockUserRepo) GetFeed(c context.Context, feedId string, userId string) (*models.Feed, error) { return m.getFeed(c, feedId, userId) }
func (m *mockUserRepo) UpdateFeed(c context.Context, f *models.Feed) (*models.Feed, error) { return m.updateFeed(c, f) }
func (m *mockUserRepo) RemoveFeed(c context.Context, feedId string, userId string) (error) { return m.removeFeed(c, feedId, userId) }
func (m *mockUserRepo) ListFeedsAll(c context.Context, userId string) ([]*models.Feed, error) { return m.listFeedsAll(c, userId) }

func TestGetUser(t *testing.T) {
  gin.SetMode(gin.TestMode)

  expectedBody := &models.UserAccount{
    UserID: "1",
    Username: "username",
    FeedList: map[string]*models.Feed{},
    ChannelList: map[string]*models.DiscordChannel{},
  }

  m := &mockUserRepo{getUser: func (c context.Context, userId string) (*models.UserAccount, error) { return expectedBody, nil }}
  usecase := userUsecase{ userRepo: m }
  c := context.Background()
  user, _ := usecase.GetUser(c, "1")
  if !reflect.DeepEqual(user, expectedBody) {
    t.Errorf("want %v, got %v", expectedBody, user)
  }
}

func TestCreateUser(t *testing.T) {
  gin.SetMode(gin.TestMode)



  expectedBody := &models.UserAccount{
    UserID: "1",
    Username: "username",
    FeedList: map[string]*models.Feed{},
    ChannelList: map[string]*models.DiscordChannel{},
  }

  m := &mockUserRepo{createUser: func (c context.Context, user *models.UserAccount) (*models.UserAccount, error) { return expectedBody, nil }}
  usecase := userUsecase{ userRepo: m }

  c := context.Background()
  u := &models.UserAccount{
    UserID: "1",
    Username: "username",
    FeedList: map[string]*models.Feed{},
    ChannelList: map[string]*models.DiscordChannel{},
  }

  user, _ := usecase.CreateUser(c, u)
  if !reflect.DeepEqual(user, expectedBody) {
    t.Errorf("want %v, got %v", expectedBody, user)
  }
}

func TestAddFeed(t *testing.T) {
  gin.SetMode(gin.TestMode)
  
  expectedBody := &models.Feed{
    FeedID: "1",
    Title: "feed",
    Url: "https://feed.com",
    TimeFormat: "none",
  }

  m := &mockUserRepo{addFeed: func (c context.Context, feed *models.Feed, userId string) (*models.Feed, error) { return expectedBody, nil }}
  usecase := userUsecase{ userRepo: m }
  

  c := context.Background()
  f := &models.Feed{
    FeedID: "1",
    Title: "feed",
    Url: "https://feed.com",
    TimeFormat: "none",
  }
  uid := "1"

  user, _ := usecase.AddFeed(c, f, uid)
  if !reflect.DeepEqual(user, expectedBody) {
    t.Errorf("want %v, got %v", expectedBody, user)
  }
}

func TestUpdateFeed(t *testing.T) {
  gin.SetMode(gin.TestMode)
  
  expectedBody := &models.Feed{
    FeedID: "1",
    Title: "feed",
    Url: "https://feed.com",
    TimeFormat: "none",
  }

  m := &mockUserRepo{updateFeed: func (c context.Context, feed *models.Feed) (*models.Feed, error) { return expectedBody, nil }}
  usecase := userUsecase{ userRepo: m }

  c := context.Background()
  f := &models.Feed{
    FeedID: "1",
    Title: "feed",
    Url: "https://feed.com",
    TimeFormat: "none",
  }

  user, _ := usecase.UpdateFeed(c, f)
  if !reflect.DeepEqual(user, expectedBody) {
    t.Errorf("want %v, got %v", expectedBody, user)
  }
}

func TestRemoveFeed(t *testing.T) {
  gin.SetMode(gin.TestMode)

  m := &mockUserRepo{removeFeed: func (c context.Context, feedId string, userId string) (error) { return nil }}
  usecase := userUsecase{ userRepo: m }

  c := context.Background()

  err := usecase.RemoveFeed(c, "1", "1")
  if err != nil {
    t.Errorf("want nil err, got %s", err)
  }
}
