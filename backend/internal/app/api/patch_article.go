package api

import (
	"TelegaFeed/internal/app/api/middlewares"
	abstractrepositories "TelegaFeed/internal/pkg/core/abstractions/infrastructure/repositories"
	abstractservices "TelegaFeed/internal/pkg/core/abstractions/services"
	"TelegaFeed/internal/pkg/core/entities"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

type patchArticleEndpointRequest struct {
	Starred *bool `json:"starred,omitempty"`
	Read    *bool `json:"read,omitempty"`
}

type PatchArticleEndpoint struct {
	feedService     abstractservices.FeedService
	usersRepository abstractrepositories.UsersRepository
}

func NewPatchArticleEndpoint(
	feedService abstractservices.FeedService,
	usersRepository abstractrepositories.UsersRepository,
) *PatchArticleEndpoint {
	return &PatchArticleEndpoint{
		feedService:     feedService,
		usersRepository: usersRepository,
	}
}

func (p PatchArticleEndpoint) Setup(e *echo.Echo) {
	e.PATCH(
		"/api/articles/:id",
		p.Execute,
		middlewares.ParseUserIDMiddleware,
		middlewares.UserExistsMiddleware(p.usersRepository),
	)
}

func (p PatchArticleEndpoint) Execute(c echo.Context) error {
	ctx := c.Request().Context()
	userId := c.Get(middlewares.UserIdContextKey).(string)

	articleId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	var req patchArticleEndpointRequest
	if err := c.Bind(&req); err != nil {
		c.Logger().Warn("failed to bind json to request's model")

		return c.NoContent(http.StatusBadRequest)
	}

	articlePatch := entities.ArticlePatch{
		Starred: entities.OptionFromNilable(req.Starred),
		Read:    entities.OptionFromNilable(req.Read),
	}

	err = p.feedService.UpdateArticle(ctx, userId, articleId, &articlePatch)

	if err != nil {
		c.Logger().Errorf("failed to update article: %v", err)

		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}
