package localapp

import (
	"TelegaFeed/internal/app/api"
	abstractproviders "TelegaFeed/internal/pkg/core/abstractions/infrastructure/providers"
	abstractrepositories "TelegaFeed/internal/pkg/core/abstractions/infrastructure/repositories"
	abstractservices "TelegaFeed/internal/pkg/core/abstractions/services"
	"TelegaFeed/internal/pkg/infrastructure/providers"
	"TelegaFeed/internal/pkg/infrastructure/repositories"
	"TelegaFeed/internal/pkg/services"
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"time"
)

type LocalApp struct {
	// framework
	echo *echo.Echo

	// endpoints
	endpoints []api.Endpoint

	// services
	feedService        abstractservices.FeedService
	feedSourcesService abstractservices.FeedSourcesService
	fetchService       abstractservices.FetchService
	llmService         abstractservices.LlmService

	// infrastructure
	llmProvider          abstractproviders.LlmProvider
	fetchProvidersMap    map[string]abstractproviders.FetchProvider
	digestsRepository    abstractrepositories.DigestsRepository
	feedRepository       abstractrepositories.FeedRepository
	feedSourceRepository abstractrepositories.FeedSourceRepository
	summariesRepository  abstractrepositories.SummariesRepository
	usersRepository      abstractrepositories.UsersRepository
}

func NewLocalApp() *LocalApp {
	return &LocalApp{}
}

func (app *LocalApp) Build() error {
	ctx, cancelFunc := context.WithTimeout(context.TODO(), 15*time.Second)

	defer cancelFunc()

	db, err := ydb.Open(ctx, "grpc://127.0.0.1:2136/local")
	if err != nil {
		return err
	}

	// infra
	app.llmProvider = providers.NewEchoLlmProvider()
	app.fetchProvidersMap = map[string]abstractproviders.FetchProvider{
		"rss": providers.NewRssProvider(),
	}
	app.digestsRepository = repositories.NewYdbDigestsRepository(db)
	app.feedRepository = repositories.NewYdbFeedRepository(db)
	app.feedSourceRepository = repositories.NewYdbFeedSourceRepository(db)
	app.summariesRepository = repositories.NewYdbSummariesRepository(db)
	app.usersRepository = repositories.NewYdbUsersRepository(db)

	// services
	app.feedSourcesService = app.feedSourceRepository
	app.fetchService = services.NewFetchService(app.fetchProvidersMap)
	app.llmService = services.NewLlmService(app.digestsRepository, app.summariesRepository, app.feedRepository, app.llmProvider)
	app.feedService = services.NewFeedService(app.llmService, app.fetchService, app.feedRepository, app.feedSourceRepository, app.usersRepository)

	// endpoints
	app.endpoints = []api.Endpoint{
		api.NewAddFeedSourceEndpoint(app.feedSourcesService),
		api.NewDeleteFeedSourceEndpoint(app.feedSourcesService),
		api.NewGetArticleSummaryEndpoint(app.llmService),
		api.NewGetFeedEndpoint(app.feedService),
		api.NewGetFeedDigestEndpoint(app.llmService),
		api.NewGetFeedSourceEndpoint(app.feedSourcesService),
		api.NewGetFeedSourcesEndpoint(app.feedSourcesService),
		api.NewPatchArticleEndpoint(app.feedService),
		api.NewPatchFeedSourceEndpoint(app.feedSourcesService),
	}

	// framework
	app.echo = echo.New()

	return nil
}

func (app *LocalApp) Setup() error {
	for _, endpoint := range app.endpoints {
		endpoint.Setup(app.echo)
	}

	return nil
}

func (app *LocalApp) Start(port int32) error {
	app.echo.Logger.Fatal(app.echo.Start(fmt.Sprintf(":%d", port)))

	return nil
}
