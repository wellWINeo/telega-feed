package api

import (
	"TelegaFeed/internal/app/api/middlewares"
	abstractrepositories "TelegaFeed/internal/pkg/core/abstractions/infrastructure/repositories"
	abstractservices "TelegaFeed/internal/pkg/core/abstractions/services"
	"github.com/labstack/echo/v4"
	"net/http"
)

type GetFeedSourcesEndpoint struct {
	feedSourcesService abstractservices.FeedSourcesService
	usersRepository    abstractrepositories.UsersRepository
}

func NewGetFeedSourcesEndpoint(
	feedSourcesService abstractservices.FeedSourcesService,
	usersRepository abstractrepositories.UsersRepository,
) *GetFeedSourcesEndpoint {
	return &GetFeedSourcesEndpoint{
		feedSourcesService: feedSourcesService,
		usersRepository:    usersRepository,
	}
}

func (g GetFeedSourcesEndpoint) Setup(e *echo.Echo) {
	e.GET(
		"/api/feed-sources",
		g.Execute,
		middlewares.ParseUserIDMiddleware,
		middlewares.UserExistsMiddleware(g.usersRepository),
	)
}

func (g GetFeedSourcesEndpoint) Execute(c echo.Context) error {
	ctx := c.Request().Context()
	userId := c.Get(middlewares.UserIdContextKey).(string)

	sources, err := g.feedSourcesService.GetSources(ctx, userId)
	if err != nil {
		c.Logger().Errorf("failed to fetch feed-sources: %v", err)

		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, sources)
}
