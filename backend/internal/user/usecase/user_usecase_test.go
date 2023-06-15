package usecase

import (
  "net/http/httptest"
  "reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mowemcfc/discRSS/models"
)

type mockUserRepo struct {
  getUser func(*gin.Context) (*models.UserAccount, error)
  createUser func(*gin.Context) (*models.UserAccount, error)
  addFeed func(*gin.Context) (*models.Feed, error)
  updateFeed func(*gin.Context) (*models.Feed, error)
  removeFeed func(*gin.Context) (error)
}

func (m *mockUserRepo) GetUser(c *gin.Context) (*models.UserAccount, error) { return m.getUser(c) }
func (m *mockUserRepo) CreateUser(c *gin.Context) (*models.UserAccount, error) { return m.createUser(c) }
func (m *mockUserRepo) AddFeed(c *gin.Context) (*models.Feed, error) { return m.addFeed(c) }
func (m *mockUserRepo) UpdateFeed(c *gin.Context) (*models.Feed, error) { return m.updateFeed(c) }
func (m *mockUserRepo) RemoveFeed(c *gin.Context) (error) { return m.removeFeed(c) }

func TestGetUser(t *testing.T) {
  gin.SetMode(gin.TestMode)

  expectedBody := &models.UserAccount{
    UserID: "1",
    Username: "username",
    FeedList: map[string]*models.Feed{},
    ChannelList: map[string]*models.DiscordChannel{},
  }

  m := &mockUserRepo{getUser: func (c *gin.Context) (*models.UserAccount, error) { return expectedBody, nil }}
  usecase := userUsecase{ userRepo: m }

  c, _ := gin.CreateTestContext(httptest.NewRecorder())
  c.Params = []gin.Param{gin.Param{Key: "id", Value: "1"}}

  user, _ := usecase.GetUser(c)
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

  m := &mockUserRepo{createUser: func (c *gin.Context) (*models.UserAccount, error) { return expectedBody, nil }}
  usecase := userUsecase{ userRepo: m }

  c, _ := gin.CreateTestContext(httptest.NewRecorder())

  user, _ := usecase.CreateUser(c)
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

  m := &mockUserRepo{addFeed: func (c *gin.Context) (*models.Feed, error) { return expectedBody, nil }}
  usecase := userUsecase{ userRepo: m }

  c, _ := gin.CreateTestContext(httptest.NewRecorder())

  user, _ := usecase.AddFeed(c)
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

  m := &mockUserRepo{updateFeed: func (c *gin.Context) (*models.Feed, error) { return expectedBody, nil }}
  usecase := userUsecase{ userRepo: m }

  c, _ := gin.CreateTestContext(httptest.NewRecorder())

  user, _ := usecase.UpdateFeed(c)
  if !reflect.DeepEqual(user, expectedBody) {
    t.Errorf("want %v, got %v", expectedBody, user)
  }
}

func TestRemoveFeed(t *testing.T) {
  gin.SetMode(gin.TestMode)

  m := &mockUserRepo{removeFeed: func (c *gin.Context) (error) { return nil }}
  usecase := userUsecase{ userRepo: m }

  c, _ := gin.CreateTestContext(httptest.NewRecorder())

  err := usecase.RemoveFeed(c)
  if err != nil {
    t.Errorf("want nil err, got %s", err)
  }
}
