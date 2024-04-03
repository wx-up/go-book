package repository

import (
	"context"

	"github.com/wx-up/go-book/internal/repository/dao"

	"github.com/wx-up/go-book/internal/repository/dao/model"

	"github.com/wx-up/go-book/internal/domain"
)

var ErrArticleNotFound = dao.ErrArticleNotFound

type ArticleRepository interface {
	Create(ctx context.Context, d domain.Article) (int64, error)
	Update(ctx context.Context, d domain.Article) error
}

type CacheArticleRepository struct {
	dao dao.ArticleDAO
}

func (c *CacheArticleRepository) Update(ctx context.Context, d domain.Article) error {
	return c.dao.UpdateById(ctx, c.toModel(d))
}

func NewCacheArticleRepository(dao dao.ArticleDAO) *CacheArticleRepository {
	return &CacheArticleRepository{dao: dao}
}

func (c *CacheArticleRepository) Create(ctx context.Context, d domain.Article) (int64, error) {
	return c.dao.Insert(ctx, c.toModel(d))
}

func (c *CacheArticleRepository) toModel(d domain.Article) model.Article {
	return model.Article{
		Id:       d.Id,
		Title:    d.Title,
		Content:  d.Content,
		AuthorId: d.Author.Id,
	}
}
