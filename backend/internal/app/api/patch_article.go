package api

import (
	"TelegaFeed/internal/app/api/middlewares"
	abstractservices "TelegaFeed/internal/pkg/core/abstractions/services"
	"TelegaFeed/internal/pkg/core/entities"
	"github.com/labstack/echo/v4"
	"net/http"
)

type patchArticleEndpointRequest struct {
	Starred *bool `json:"starred,omitempty"`
	Read    *bool `json:"read,omitempty"`
}

type PatchArticleEndpoint struct {
	feedService abstractservices.FeedService
}

func NewPatchArticleEndpoint(feedService abstractservices.FeedService) *PatchArticleEndpoint {
	return &PatchArticleEndpoint{feedService: feedService}
}

func (p PatchArticleEndpoint) Setup(e *echo.Echo) {
	e.PATCH("/api/articles/:id", p.Execute, middlewares.ParseUserIDMiddleware)
}

func (p PatchArticleEndpoint) Execute(c echo.Context) error {
	ctx := c.Request().Context()
	userId := c.Get(middlewares.UserIdContextKey).(string)
	articleId := c.Param("id")

	var req patchArticleEndpointRequest
	if err := c.Bind(&req); err != nil {
		c.Logger().Warn("failed to bind json to request's model")

		return c.NoContent(http.StatusBadRequest)
	}

	articlePatch := entities.ArticlePatch{
		Starred: entities.OptionFromNilable(req.Starred),
		Read:    entities.OptionFromNilable(req.Read),
	}

	err := p.feedService.UpdateArticle(ctx, userId, articleId, &articlePatch)

	if err != nil {
		c.Logger().Errorf("failed to update article: %v", err)

		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}
