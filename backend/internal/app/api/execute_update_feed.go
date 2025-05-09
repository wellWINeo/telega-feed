package api

import (
	abstractservices "TelegaFeed/internal/pkg/core/abstractions/services"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ExecuteUpdateFeedEndpoint struct {
	feedService abstractservices.FeedService
}

func NewExecuteUpdateFeedEndpoint(feedService abstractservices.FeedService) *ExecuteUpdateFeedEndpoint {
	return &ExecuteUpdateFeedEndpoint{feedService: feedService}
}

func (e ExecuteUpdateFeedEndpoint) Setup(ee *echo.Echo) {
	ee.POST("/api/execute/update-feed", e.Execute)
}

func (e ExecuteUpdateFeedEndpoint) Execute(c echo.Context) error {
	err := e.feedService.UpdateFeed(c.Request().Context())

	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}
