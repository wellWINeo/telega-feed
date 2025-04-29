package api

import (
	"TelegaFeed/internal/app/api/middlewares"
	abstractservices "TelegaFeed/internal/pkg/core/abstractions/services"
	"github.com/labstack/echo/v4"
	"net/http"
)

type GetArticleSummaryEndpoint struct {
	llmService abstractservices.LlmService
}

func NewGetArticleSummaryEndpoint(llmService abstractservices.LlmService) *GetArticleSummaryEndpoint {
	return &GetArticleSummaryEndpoint{llmService: llmService}
}

func (g GetArticleSummaryEndpoint) Setup(e *echo.Echo) {
	e.GET("/api/articles/:id/summary", g.Execute, middlewares.ParseUserIDMiddleware)
}

func (g GetArticleSummaryEndpoint) Execute(c echo.Context) error {
	ctx := c.Request().Context()
	userId := c.Get(middlewares.UserIdContextKey).(string)
	articleId := c.Param("id")

	summary, err := g.llmService.GetArticleSummary(ctx, userId, articleId)
	if err != nil {
		c.Logger().Errorf("failed to get summary by id: %v", err)

		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, summary)
}
