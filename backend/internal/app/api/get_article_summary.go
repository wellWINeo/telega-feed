package api

import (
	"TelegaFeed/internal/app/api/middlewares"
	abstractrepositories "TelegaFeed/internal/pkg/core/abstractions/infrastructure/repositories"
	abstractservices "TelegaFeed/internal/pkg/core/abstractions/services"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

type GetArticleSummaryEndpoint struct {
	llmService      abstractservices.LlmService
	usersRepository abstractrepositories.UsersRepository
}

func NewGetArticleSummaryEndpoint(
	llmService abstractservices.LlmService,
	usersRepository abstractrepositories.UsersRepository,
) *GetArticleSummaryEndpoint {
	return &GetArticleSummaryEndpoint{
		llmService:      llmService,
		usersRepository: usersRepository,
	}
}

func (g GetArticleSummaryEndpoint) Setup(e *echo.Echo) {
	e.GET(
		"/api/articles/:id/summary",
		g.Execute,
		middlewares.ParseUserIDMiddleware,
		middlewares.UserExistsMiddleware(g.usersRepository),
	)
}

func (g GetArticleSummaryEndpoint) Execute(c echo.Context) error {
	ctx := c.Request().Context()
	userId := c.Get(middlewares.UserIdContextKey).(string)

	articleId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	summary, err := g.llmService.GetArticleSummary(ctx, userId, articleId)
	if err != nil {
		c.Logger().Errorf("failed to get summary by id: %v", err)

		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, summary)
}
