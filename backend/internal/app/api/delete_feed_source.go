package api

import (
	"TelegaFeed/internal/app/api/middlewares"
	abstractservices "TelegaFeed/internal/pkg/core/abstractions/services"
	"github.com/labstack/echo/v4"
	"net/http"
)

type DeleteFeedSourceEndpoint struct {
	feedSourcesService abstractservices.FeedSourcesService
}

func NewDeleteFeedSourceEndpoint(feedSourcesService abstractservices.FeedSourcesService) *DeleteFeedSourceEndpoint {
	return &DeleteFeedSourceEndpoint{feedSourcesService: feedSourcesService}
}

func (d DeleteFeedSourceEndpoint) Setup(e *echo.Echo) {
	e.DELETE("/api/feed-source/:id", d.Execute, middlewares.ParseUserIDMiddleware)
}

func (d DeleteFeedSourceEndpoint) Execute(c echo.Context) error {
	ctx := c.Request().Context()
	userId := c.Get(middlewares.UserIdContextKey).(string)
	sourceId := c.Param("id")

	err := d.feedSourcesService.DeleteSource(ctx, userId, sourceId)
	if err != nil {
		c.Logger().Errorf("failed to delete feed source: %v", err)

		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}
