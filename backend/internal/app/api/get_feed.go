package api

import (
	"TelegaFeed/internal/app/api/middlewares"
	abstractservices "TelegaFeed/internal/pkg/core/abstractions/services"
	"github.com/labstack/echo/v4"
	"net/http"
)

type GetFeedEndpoint struct {
	feedService abstractservices.FeedService
}

func NewGetFeedEndpoint(feedService abstractservices.FeedService) *GetFeedEndpoint {
	return &GetFeedEndpoint{feedService: feedService}
}

func (g *GetFeedEndpoint) Setup(e *echo.Echo) {
	e.GET("/api/feed", g.Execute, middlewares.ParseUserIDMiddleware)
}

func (g *GetFeedEndpoint) Execute(c echo.Context) error {
	ctx := c.Request().Context()
	userId := c.Get(middlewares.UserIdContextKey).(string)

	feed, err := g.feedService.GetFeed(ctx, userId)
	if err != nil {
		c.Logger().Errorf("failed to fetch feed: %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, feed)
}
