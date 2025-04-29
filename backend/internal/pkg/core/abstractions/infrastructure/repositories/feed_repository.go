package abstractrepositories

import (
	"TelegaFeed/internal/pkg/core/entities"
	"context"
)

type FeedRepository interface {
	AddArticlesToFeed(ctx context.Context, articles []*entities.Article) error
	GetFeedByUser(ctx context.Context, userId entities.UserId) ([]*entities.Article, error)
	GetTodayArticles(ctx context.Context, userId entities.UserId) ([]*entities.Article, error)
	GetArticleById(ctx context.Context, userId entities.UserId, articleId entities.ArticleId) (*entities.Article, error)
	UpdateArticle(ctx context.Context, userId entities.UserId, articleId entities.ArticleId, patch *entities.ArticlePatch) (*entities.Article, error)
}
