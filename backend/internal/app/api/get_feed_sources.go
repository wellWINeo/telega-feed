package api

import (
	"TelegaFeed/internal/app/api/middlewares"
	abstractservices "TelegaFeed/internal/pkg/core/abstractions/services"
	"github.com/labstack/echo/v4"
	"net/http"
)

type GetFeedSourcesEndpoint struct {
	feedSourcesService abstractservices.FeedSourcesService
}

func NewGetFeedSourcesEndpoint(feedSourcesService abstractservices.FeedSourcesService) *GetFeedSourcesEndpoint {
	return &GetFeedSourcesEndpoint{feedSourcesService: feedSourcesService}
}

func (g GetFeedSourcesEndpoint) Setup(e *echo.Echo) {
	e.GET("/api/feed-sources", g.Execute, middlewares.ParseUserIDMiddleware)
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
