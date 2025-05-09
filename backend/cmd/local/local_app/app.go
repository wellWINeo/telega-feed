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
	"os"
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

	dsn := os.Getenv("YDB_CONNECTIONSTRING")
	if dsn == "" {
		dsn = "grpc://127.0.0.1:2136/local"
	}

	ydbDB, err := ydb.Open(ctx, dsn)
	if err != nil {
		return err
	}

	// infra
	app.llmProvider = providers.NewEchoLlmProvider()
	app.fetchProvidersMap = map[string]abstractproviders.FetchProvider{
		"rss":  providers.NewRssProvider(),
		"atom": providers.NewAtomProvider(),
		"rdf":  providers.NewRDFProvider(),
	}
	app.digestsRepository = repositories.NewYdbDigestsRepository(ydbDB)
	app.feedRepository = repositories.NewYdbFeedRepository(ydbDB)
	app.feedSourceRepository = repositories.NewYdbFeedSourceRepository(ydbDB)
	app.summariesRepository = repositories.NewYdbSummariesRepository(ydbDB)
	app.usersRepository = repositories.NewYdbUsersRepository(ydbDB)

	// services
	app.fetchService = services.NewFetchService(app.fetchProvidersMap)
	app.feedSourcesService = services.NewFeedSourcesService(app.fetchService, app.feedSourceRepository)
	app.llmService = services.NewLlmService(app.digestsRepository, app.summariesRepository, app.feedRepository, app.llmProvider)
	app.feedService = services.NewFeedService(app.llmService, app.fetchService, app.feedRepository, app.feedSourceRepository, app.usersRepository)

	// endpoints
	app.endpoints = []api.Endpoint{
		api.NewAddFeedSourceEndpoint(app.feedSourcesService, app.usersRepository),
		api.NewDeleteFeedSourceEndpoint(app.feedSourcesService, app.usersRepository),
		api.NewGetArticleSummaryEndpoint(app.llmService, app.usersRepository),
		api.NewGetFeedEndpoint(app.feedService, app.usersRepository),
		api.NewGetFeedDigestEndpoint(app.llmService, app.usersRepository),
		api.NewGetFeedSourceEndpoint(app.feedSourcesService, app.usersRepository),
		api.NewGetFeedSourcesEndpoint(app.feedSourcesService, app.usersRepository),
		api.NewPatchArticleEndpoint(app.feedService, app.usersRepository),
		api.NewPatchFeedSourceEndpoint(app.feedSourcesService, app.usersRepository),
		api.NewExecuteUpdateFeedEndpoint(app.feedService),
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
