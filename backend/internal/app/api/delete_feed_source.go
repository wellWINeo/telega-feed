package api

import (
	"TelegaFeed/internal/app/api/middlewares"
	abstractrepositories "TelegaFeed/internal/pkg/core/abstractions/infrastructure/repositories"
	abstractservices "TelegaFeed/internal/pkg/core/abstractions/services"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

type DeleteFeedSourceEndpoint struct {
	feedSourcesService abstractservices.FeedSourcesService
	usersRepository    abstractrepositories.UsersRepository
}

func NewDeleteFeedSourceEndpoint(
	feedSourcesService abstractservices.FeedSourcesService,
	usersRepository abstractrepositories.UsersRepository,
) *DeleteFeedSourceEndpoint {
	return &DeleteFeedSourceEndpoint{
		feedSourcesService: feedSourcesService,
		usersRepository:    usersRepository,
	}
}

func (d DeleteFeedSourceEndpoint) Setup(e *echo.Echo) {
	e.DELETE(
		"/api/feed-sources/:id",
		d.Execute,
		middlewares.ParseUserIDMiddleware,
		middlewares.UserExistsMiddleware(d.usersRepository),
	)
}

func (d DeleteFeedSourceEndpoint) Execute(c echo.Context) error {
	ctx := c.Request().Context()
	userId := c.Get(middlewares.UserIdContextKey).(string)

	sourceId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if err := d.feedSourcesService.DeleteSource(ctx, userId, sourceId); err != nil {
		c.Logger().Errorf("failed to delete feed source: %v", err)

		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}
