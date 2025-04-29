package api

import (
	"TelegaFeed/internal/app/api/middlewares"
	abstractservices "TelegaFeed/internal/pkg/core/abstractions/services"
	"TelegaFeed/internal/pkg/core/entities"
	"github.com/labstack/echo/v4"
	"net/http"
)

type addFeedSourceEndpointRequest struct {
	Name    string `json:"name"`
	FeedUrl string `json:"feed_url"`
}

type AddFeedSourceEndpoint struct {
	feedSourcesService abstractservices.FeedSourcesService
}

func NewAddFeedSourceEndpoint(feedSourcesService abstractservices.FeedSourcesService) *AddFeedSourceEndpoint {
	return &AddFeedSourceEndpoint{feedSourcesService: feedSourcesService}
}

func (a AddFeedSourceEndpoint) Setup(e *echo.Echo) {
	e.POST("/api/feed-sources", a.Execute, middlewares.ParseUserIDMiddleware)
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
		Id:       "",
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
