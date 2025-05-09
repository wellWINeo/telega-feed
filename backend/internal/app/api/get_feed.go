package api

import (
	"TelegaFeed/internal/app/api/middlewares"
	abstractrepositories "TelegaFeed/internal/pkg/core/abstractions/infrastructure/repositories"
	abstractservices "TelegaFeed/internal/pkg/core/abstractions/services"
	"github.com/labstack/echo/v4"
	"net/http"
)

type GetFeedEndpoint struct {
	feedService     abstractservices.FeedService
	usersRepository abstractrepositories.UsersRepository
}

func NewGetFeedEndpoint(
	feedService abstractservices.FeedService,
	usersRepository abstractrepositories.UsersRepository,
) *GetFeedEndpoint {
	return &GetFeedEndpoint{
		feedService:     feedService,
		usersRepository: usersRepository,
	}
}

func (g *GetFeedEndpoint) Setup(e *echo.Echo) {
	e.GET(
		"/api/feed",
		g.Execute,
		middlewares.ParseUserIDMiddleware,
		middlewares.UserExistsMiddleware(g.usersRepository),
	)
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
