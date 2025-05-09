package api

import (
	"TelegaFeed/internal/app/api/middlewares"
	abstractrepositories "TelegaFeed/internal/pkg/core/abstractions/infrastructure/repositories"
	abstractservices "TelegaFeed/internal/pkg/core/abstractions/services"
	"github.com/labstack/echo/v4"
	"net/http"
)

type GetFeedDigestEndpoint struct {
	llmService      abstractservices.LlmService
	usersRepository abstractrepositories.UsersRepository
}

func NewGetFeedDigestEndpoint(
	llmService abstractservices.LlmService,
	usersRepository abstractrepositories.UsersRepository,
) *GetFeedDigestEndpoint {
	return &GetFeedDigestEndpoint{
		llmService:      llmService,
		usersRepository: usersRepository,
	}
}

func (g GetFeedDigestEndpoint) Setup(e *echo.Echo) {
	e.GET(
		"/api/feed/digest",
		g.Execute,
		middlewares.ParseUserIDMiddleware,
		middlewares.UserExistsMiddleware(g.usersRepository),
	)
}

func (g GetFeedDigestEndpoint) Execute(e echo.Context) error {
	ctx := e.Request().Context()
	userId := e.Get(middlewares.UserIdContextKey).(string)

	digest, err := g.llmService.GetDailyDigest(ctx, userId)
	if err != nil {
		e.Logger().Errorf("failed to get digest for user %s: %v", userId, err)

		return e.NoContent(http.StatusInternalServerError)
	}

	return e.JSON(http.StatusOK, digest)
}
