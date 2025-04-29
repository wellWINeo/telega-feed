package api

import (
	"TelegaFeed/internal/app/api/middlewares"
	abstractservices "TelegaFeed/internal/pkg/core/abstractions/services"
	"TelegaFeed/internal/pkg/core/entities"
	"github.com/labstack/echo/v4"
	"net/http"
)

type patchFeedSourceEndpointRequest struct {
	Name     *string `json:"name,omitempty"`
	Disabled *bool   `json:"disabled,omitempty"`
}

type PatchFeedSourceEndpoint struct {
	feedSourcesService abstractservices.FeedSourcesService
}

func NewPatchFeedSourceEndpoint(feedSourcesService abstractservices.FeedSourcesService) *PatchFeedSourceEndpoint {
	return &PatchFeedSourceEndpoint{feedSourcesService: feedSourcesService}
}

func (p PatchFeedSourceEndpoint) Setup(e *echo.Echo) {
	e.PATCH("/api/feed-source/:id", p.Execute, middlewares.ParseUserIDMiddleware)
}

func (p PatchFeedSourceEndpoint) Execute(c echo.Context) error {
	ctx := c.Request().Context()
	userId := c.Get(middlewares.UserIdContextKey).(string)
	sourceId := c.Param("id")

	var req patchFeedSourceEndpointRequest
	if err := c.Bind(&req); err != nil {
		c.Logger().Warnf("failed to bind json to request: %v", err)

		return c.NoContent(http.StatusBadRequest)
	}

	patch := &entities.FeedSourcePatch{
		Name:     entities.OptionFromNilable(req.Name),
		Disabled: entities.OptionFromNilable(req.Disabled),
	}

	err := p.feedSourcesService.UpdateSource(ctx, userId, sourceId, patch)
	if err != nil {
		c.Logger().Warnf("failed to update feed source: %v", err)

		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}
