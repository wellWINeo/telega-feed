package api

import (
	"TelegaFeed/internal/app/api/middlewares"
	abstractservices "TelegaFeed/internal/pkg/core/abstractions/services"
	"github.com/labstack/echo/v4"
	"net/http"
)

type GetFeedSourceEndpoint struct {
	feedSourcesService abstractservices.FeedSourcesService
}

func NewGetFeedSourceEndpoint(feedSourcesService abstractservices.FeedSourcesService) *GetFeedSourceEndpoint {
	return &GetFeedSourceEndpoint{feedSourcesService: feedSourcesService}
}

func (g GetFeedSourceEndpoint) Setup(e *echo.Echo) {
	e.GET("/api/feed-source/:id", g.Execute, middlewares.ParseUserIDMiddleware)
}

func (g GetFeedSourceEndpoint) Execute(c echo.Context) error {
	ctx := c.Request().Context()
	userId := c.Get(middlewares.UserIdContextKey).(string)
	sourceId := c.Param("id")

	source, err := g.feedSourcesService.GetSource(ctx, userId, sourceId)

	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, source)
}
