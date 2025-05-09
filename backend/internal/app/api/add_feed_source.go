package api

import (
	"TelegaFeed/internal/app/api/middlewares"
	abstractrepositories "TelegaFeed/internal/pkg/core/abstractions/infrastructure/repositories"
	abstractservices "TelegaFeed/internal/pkg/core/abstractions/services"
	"TelegaFeed/internal/pkg/core/entities"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

type addFeedSourceEndpointRequest struct {
	Name    string `json:"name"`
	FeedUrl string `json:"feed_url"`
}

type AddFeedSourceEndpoint struct {
	feedSourcesService abstractservices.FeedSourcesService
	usersRepository    abstractrepositories.UsersRepository
}

func NewAddFeedSourceEndpoint(
	feedSourcesService abstractservices.FeedSourcesService,
	usersRepository abstractrepositories.UsersRepository,
) *AddFeedSourceEndpoint {
	return &AddFeedSourceEndpoint{
		feedSourcesService: feedSourcesService,
		usersRepository:    usersRepository,
	}
}

func (a AddFeedSourceEndpoint) Setup(e *echo.Echo) {
	e.POST(
		"/api/feed-sources",
		a.Execute,
		middlewares.ParseUserIDMiddleware,
		middlewares.UserExistsMiddleware(a.usersRepository),
	)
}

func (a AddFeedSourceEndpoint) Execute(c echo.Context) error {
	ctx := c.Request().Context()
	userId := c.Get(middlewares.UserIdContextKey).(string)

	var req addFeedSourceEndpointRequest
	if err := c.Bind(&req); err != nil {
		c.Logger().Warnf("failed to bind json to request: %v", err)

		return c.NoContent(http.StatusBadRequest)
	}

	feedSource := entities.FeedSource{
		Id:       uuid.Nil,
		Name:     req.Name,
		FeedUrl:  req.FeedUrl,
		Disabled: false,
	}

	err := a.feedSourcesService.AddSource(ctx, userId, &feedSource)
	if err != nil {
		c.Logger().Errorf("failed to add feed source: %v", err)

		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}
