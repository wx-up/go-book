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

	// Sync 存储并同步到线上库
	Sync(ctx context.Context, d domain.Article) (int64, error)
	// SyncStatus 同步状态
	SyncStatus(ctx context.Context, authorId int64, id int64, status domain.ArticleStatus) error
}

type CacheArticleRepository struct {
	// 制作库的 dao
	dao dao.ArticleDAO
}

func (c *CacheArticleRepository) SyncStatus(ctx context.Context, authorId int64, id int64, status domain.ArticleStatus) error {
	return c.dao.SyncStatus(ctx, authorId, id, status.ToUint8())
}

func (c *CacheArticleRepository) Sync(ctx context.Context, d domain.Article) (int64, error) {
	var (
		id  = d.Id
		err error
		obj = c.toModel(d)
	)
	err = c.dao.Transaction(ctx, func(dao dao.ArticleDAO, readerDao dao.ReaderArticleDAO) error {
		if id > 0 {
			err = dao.UpdateById(ctx, obj)
		} else {
			id, err = dao.Insert(ctx, obj)
		}
		if err != nil {
			return err
		}

		// 同步给线上库
		id, err = readerDao.Upsert(ctx, model.PublishArticle{Article: obj})
		return err
	})
	return id, err
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
