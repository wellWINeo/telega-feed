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

type patchFeedSourceEndpointRequest struct {
	Name     *string `json:"name,omitempty"`
	Disabled *bool   `json:"disabled,omitempty"`
}

type PatchFeedSourceEndpoint struct {
	feedSourcesService abstractservices.FeedSourcesService
	usersRepository    abstractrepositories.UsersRepository
}

func NewPatchFeedSourceEndpoint(
	feedSourcesService abstractservices.FeedSourcesService,
	usersRepository abstractrepositories.UsersRepository,
) *PatchFeedSourceEndpoint {
	return &PatchFeedSourceEndpoint{
		feedSourcesService: feedSourcesService,
		usersRepository:    usersRepository,
	}
}

func (p PatchFeedSourceEndpoint) Setup(e *echo.Echo) {
	e.PATCH(
		"/api/feed-sources/:id",
		p.Execute,
		middlewares.ParseUserIDMiddleware,
		middlewares.UserExistsMiddleware(p.usersRepository),
	)
}

func (p PatchFeedSourceEndpoint) Execute(c echo.Context) error {
	ctx := c.Request().Context()
	userId := c.Get(middlewares.UserIdContextKey).(string)
	sourceId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	var req patchFeedSourceEndpointRequest
	if err := c.Bind(&req); err != nil {
		c.Logger().Warnf("failed to bind json to request: %v", err)

		return c.NoContent(http.StatusBadRequest)
	}

	patch := &entities.FeedSourcePatch{
		Name:     entities.OptionFromNilable(req.Name),
		Disabled: entities.OptionFromNilable(req.Disabled),
	}

	if err := p.feedSourcesService.UpdateSource(ctx, userId, sourceId, patch); err != nil {
		c.Logger().Warnf("failed to update feed source: %v", err)

		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}
