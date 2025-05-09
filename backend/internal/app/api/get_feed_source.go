package api

import (
	"TelegaFeed/internal/app/api/middlewares"
	abstractrepositories "TelegaFeed/internal/pkg/core/abstractions/infrastructure/repositories"
	abstractservices "TelegaFeed/internal/pkg/core/abstractions/services"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

type GetFeedSourceEndpoint struct {
	feedSourcesService abstractservices.FeedSourcesService
	usersRepository    abstractrepositories.UsersRepository
}

func NewGetFeedSourceEndpoint(
	feedSourcesService abstractservices.FeedSourcesService,
	usersRepository abstractrepositories.UsersRepository,
) *GetFeedSourceEndpoint {
	return &GetFeedSourceEndpoint{
		feedSourcesService: feedSourcesService,
		usersRepository:    usersRepository,
	}
}

func (g GetFeedSourceEndpoint) Setup(e *echo.Echo) {
	e.GET(
		"/api/feed-sources/:id",
		g.Execute,
		middlewares.ParseUserIDMiddleware,
		middlewares.UserExistsMiddleware(g.usersRepository),
	)
}

func (g GetFeedSourceEndpoint) Execute(c echo.Context) error {
	ctx := c.Request().Context()
	userId := c.Get(middlewares.UserIdContextKey).(string)
	sourceId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	source, err := g.feedSourcesService.GetSource(ctx, userId, sourceId)

	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, source)
}
