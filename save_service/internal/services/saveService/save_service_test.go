package saveService

import (
	"context"
	"errors"
	"testing"
	"time"

	"gonews/save_service/internal/models"
	"gonews/save_service/internal/services/saveService/mocks"

	"github.com/stretchr/testify/suite"
	"gotest.tools/v3/assert"
)

type SaveServiceSuite struct {
	suite.Suite
	ctx         context.Context
	newsStorage *mocks.MockNewsStorage
	saveService *SaveService
}

func (s *SaveServiceSuite) SetupTest() {
	s.newsStorage = mocks.NewMockNewsStorage(s.T())
	s.ctx = context.Background()
	s.saveService = NewSaveService(s.ctx, s.newsStorage)
}

func (s *SaveServiceSuite) TestCreateUserSuccess() {
	userName := "Test User"
	expectedID := uint64(123)

	s.newsStorage.EXPECT().CreateUser(s.ctx, userName).Return(expectedID, nil)

	actualID, err := s.saveService.CreateUser(s.ctx, userName)

	assert.NilError(s.T(), err)
	assert.Equal(s.T(), expectedID, actualID)
}

func (s *SaveServiceSuite) TestCreateUserStorageError() {
	userName := "Test User"
	wantErr := errors.New("storage error")

	s.newsStorage.EXPECT().CreateUser(s.ctx, userName).Return(uint64(0), wantErr)

	actualID, err := s.saveService.CreateUser(s.ctx, userName)

	assert.ErrorIs(s.T(), err, wantErr)
	assert.Equal(s.T(), uint64(0), actualID)
}

func (s *SaveServiceSuite) TestAddFavouriteSuccess() {
	userID := uint64(1)
	newsID := uint64(100)

	s.newsStorage.EXPECT().AddFavourite(s.ctx, userID, newsID).Return(nil)

	err := s.saveService.AddFavourite(s.ctx, userID, newsID)

	assert.NilError(s.T(), err)
}

func (s *SaveServiceSuite) TestAddFavouriteStorageError() {
	userID := uint64(1)
	newsID := uint64(100)
	wantErr := errors.New("storage error")

	s.newsStorage.EXPECT().AddFavourite(s.ctx, userID, newsID).Return(wantErr)

	err := s.saveService.AddFavourite(s.ctx, userID, newsID)

	assert.ErrorIs(s.T(), err, wantErr)
}

func (s *SaveServiceSuite) TestGetFavouritesSuccess() {
	userID := uint64(1)
	expectedNews := []*models.News{
		{
			ID:          100,
			Source:      "CNN",
			Author:      "John Doe",
			Title:       "Breaking News",
			Description: "Test description",
			URL:         "https://example.com",
			ImageURL:    "https://example.com/image.jpg",
			PublishedAt: time.Now(),
		},
	}

	s.newsStorage.EXPECT().GetFavourites(s.ctx, userID).Return(expectedNews, nil)

	actualNews, err := s.saveService.GetFavourites(s.ctx, userID)

	assert.NilError(s.T(), err)
	assert.DeepEqual(s.T(), expectedNews, actualNews)
}

func (s *SaveServiceSuite) TestGetFavouritesStorageError() {
	userID := uint64(1)
	wantErr := errors.New("storage error")

	s.newsStorage.EXPECT().GetFavourites(s.ctx, userID).Return(nil, wantErr)

	actualNews, err := s.saveService.GetFavourites(s.ctx, userID)

	assert.ErrorIs(s.T(), err, wantErr)
	assert.Assert(s.T(), actualNews == nil)
}

func (s *SaveServiceSuite) TestAddToSearchHistorySuccess() {
	userID := uint64(1)
	query := "bitcoin news"
	results := []uint64{100, 101}

	s.newsStorage.EXPECT().AddToSearchHistory(s.ctx, userID, query, results).Return(nil)

	err := s.saveService.AddToSearchHistory(s.ctx, userID, query, results)

	assert.NilError(s.T(), err)
}

func (s *SaveServiceSuite) TestAddToSearchHistoryStorageError() {
	userID := uint64(1)
	query := "bitcoin news"
	results := []uint64{100, 101}
	wantErr := errors.New("storage error")

	s.newsStorage.EXPECT().AddToSearchHistory(s.ctx, userID, query, results).Return(wantErr)

	err := s.saveService.AddToSearchHistory(s.ctx, userID, query, results)

	assert.ErrorIs(s.T(), err, wantErr)
}

func (s *SaveServiceSuite) TestGetSearchHistorySuccess() {
	userID := uint64(1)
	expectedQueries := []string{"bitcoin", "ethereum", "crypto"}

	s.newsStorage.EXPECT().GetSearchHistory(s.ctx, userID).Return(expectedQueries, nil)

	actualQueries, err := s.saveService.GetSearchHistory(s.ctx, userID)

	assert.NilError(s.T(), err)
	assert.DeepEqual(s.T(), expectedQueries, actualQueries)
}

func (s *SaveServiceSuite) TestGetSearchHistoryStorageError() {
	userID := uint64(1)
	wantErr := errors.New("storage error")

	s.newsStorage.EXPECT().GetSearchHistory(s.ctx, userID).Return(nil, wantErr)

	actualQueries, err := s.saveService.GetSearchHistory(s.ctx, userID)

	assert.ErrorIs(s.T(), err, wantErr)
	assert.Assert(s.T(), actualQueries == nil)
}

func (s *SaveServiceSuite) TestSubscribeSuccess() {
	userID := uint64(1)
	keyword := "technology"

	s.newsStorage.EXPECT().Subscribe(s.ctx, userID, keyword).Return(nil)

	err := s.saveService.Subscribe(s.ctx, userID, keyword)

	assert.NilError(s.T(), err)
}

func (s *SaveServiceSuite) TestSubscribeStorageError() {
	userID := uint64(1)
	keyword := "technology"
	wantErr := errors.New("storage error")

	s.newsStorage.EXPECT().Subscribe(s.ctx, userID, keyword).Return(wantErr)

	err := s.saveService.Subscribe(s.ctx, userID, keyword)

	assert.ErrorIs(s.T(), err, wantErr)
}

func (s *SaveServiceSuite) TestGetSubscriptionsSuccess() {
	expectedSubscriptions := []*models.Subscription{
		{
			ID:      1,
			UserID:  100,
			Keyword: "bitcoin",
		},
		{
			ID:      2,
			UserID:  101,
			Keyword: "ethereum",
		},
	}

	s.newsStorage.EXPECT().GetSubscriptions(s.ctx).Return(expectedSubscriptions, nil)

	actualSubscriptions, err := s.saveService.GetSubscriptions(s.ctx)

	assert.NilError(s.T(), err)
	assert.DeepEqual(s.T(), expectedSubscriptions, actualSubscriptions)
}

func (s *SaveServiceSuite) TestGetSubscriptionsStorageError() {
	wantErr := errors.New("storage error")

	s.newsStorage.EXPECT().GetSubscriptions(s.ctx).Return(nil, wantErr)

	actualSubscriptions, err := s.saveService.GetSubscriptions(s.ctx)

	assert.ErrorIs(s.T(), err, wantErr)
	assert.Assert(s.T(), actualSubscriptions == nil)
}

func (s *SaveServiceSuite) TestMarkNewsAsSeenSuccess() {
	userID := uint64(1)
	newsID := uint64(100)

	s.newsStorage.EXPECT().MarkNewsAsSeen(s.ctx, userID, newsID).Return(nil)

	err := s.saveService.MarkNewsAsSeen(s.ctx, userID, newsID)

	assert.NilError(s.T(), err)
}

func (s *SaveServiceSuite) TestMarkNewsAsSeenStorageError() {
	userID := uint64(1)
	newsID := uint64(100)
	wantErr := errors.New("storage error")

	s.newsStorage.EXPECT().MarkNewsAsSeen(s.ctx, userID, newsID).Return(wantErr)

	err := s.saveService.MarkNewsAsSeen(s.ctx, userID, newsID)

	assert.ErrorIs(s.T(), err, wantErr)
}

func (s *SaveServiceSuite) TestSaveNewsSuccess() {
	news := []*models.News{
		{
			Source:      "CNN",
			Author:      "John Doe",
			Title:       "Breaking News",
			Description: "Something important happened",
			URL:         "https://example.com/news1",
			ImageURL:    "https://example.com/image1.jpg",
			PublishedAt: time.Now(),
		},
	}

	s.newsStorage.EXPECT().UpsertNews(s.ctx, news).Return(nil)

	err := s.saveService.SaveNews(s.ctx, news)

	assert.NilError(s.T(), err)
}

func (s *SaveServiceSuite) TestSaveNewsStorageError() {
	news := []*models.News{
		{
			Source:      "CNN",
			Author:      "John Doe",
			Title:       "Breaking News",
			Description: "Something important happened",
			URL:         "https://example.com/news1",
			ImageURL:    "https://example.com/image1.jpg",
			PublishedAt: time.Now(),
		},
	}
	wantErr := errors.New("storage error")

	s.newsStorage.EXPECT().UpsertNews(s.ctx, news).Return(wantErr)

	err := s.saveService.SaveNews(s.ctx, news)

	assert.ErrorIs(s.T(), err, wantErr)
}

func (s *SaveServiceSuite) TestGetNewsByIDsSuccess() {
	ids := []uint64{100, 101}
	expectedNews := []*models.News{
		{
			ID:          100,
			Source:      "CNN",
			Author:      "John Doe",
			Title:       "News 1",
			Description: "Description 1",
			URL:         "https://example.com/news1",
			ImageURL:    "https://example.com/image1.jpg",
			PublishedAt: time.Now(),
		},
	}

	s.newsStorage.EXPECT().GetNewsByIDs(s.ctx, ids).Return(expectedNews, nil)

	actualNews, err := s.saveService.GetNewsByIDs(s.ctx, ids)

	assert.NilError(s.T(), err)
	assert.DeepEqual(s.T(), expectedNews, actualNews)
}

func (s *SaveServiceSuite) TestGetNewsByIDsStorageError() {
	ids := []uint64{100, 101}
	wantErr := errors.New("storage error")

	s.newsStorage.EXPECT().GetNewsByIDs(s.ctx, ids).Return(nil, wantErr)

	actualNews, err := s.saveService.GetNewsByIDs(s.ctx, ids)

	assert.ErrorIs(s.T(), err, wantErr)
	assert.Assert(s.T(), actualNews == nil)
}

func (s *SaveServiceSuite) TestNewSaveService() {
	service := NewSaveService(s.ctx, s.newsStorage)
	assert.Assert(s.T(), service != nil)
	assert.Equal(s.T(), s.newsStorage, service.newsStorage)
}

func TestSaveServiceSuite(t *testing.T) {
	suite.Run(t, new(SaveServiceSuite))
}
