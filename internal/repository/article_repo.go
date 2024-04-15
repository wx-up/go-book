package repository

import (
	"context"
	"time"

	"github.com/wx-up/go-book/pkg/logger"

	"github.com/wx-up/go-book/internal/repository/cache"

	"github.com/wx-up/go-book/pkg/slice"

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

	GetByAuthorId(ctx context.Context, authorId int64, page, pageSize int64) ([]domain.Article, error)

	// GetById 获取制作库详情
	GetById(ctx context.Context, id int64) (domain.Article, error)
	// GetPublishedById 获取线上库详情
	GetPublishedById(ctx context.Context, id int64) (domain.Article, error)
}

type CacheArticleRepository struct {
	// 制作库的 dao
	dao dao.ArticleDAO
	// 线上库的 dao
	readerDAO dao.ReaderArticleDAO

	userRepo UserRepository

	cache cache.ArticleCache

	logger logger.Logger
}

func (c *CacheArticleRepository) GetPublishedById(ctx context.Context, id int64) (domain.Article, error) {
	res, err := c.readerDAO.GetById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}

	// 获取用户信息
	user, err := c.userRepo.FindById(ctx, res.AuthorId)
	if err != nil {
		return domain.Article{}, err
	}
	dm := c.toDomain(res.Article)
	dm.Author.Id = user.Id
	dm.Author.Name = user.Nickname
	return dm, nil
}

func (c *CacheArticleRepository) GetById(ctx context.Context, id int64) (domain.Article, error) {
	// TODO implement me
	panic("implement me")
}

func (c *CacheArticleRepository) GetByAuthorId(ctx context.Context, authorId int64, page, pageSize int64) ([]domain.Article, error) {
	// 缓存第一页数据，即前端页面打开时默认请求的查询参数
	// 假设默认请求的参数 page == 1，pageSize == 100
	if page == 1 && pageSize == 100 {
		data, err := c.cache.GetFirstPage(ctx, authorId)
		if err == nil {
			return data, nil
		}
	}
	res, err := c.dao.GetByAuthorId(ctx, authorId, page, pageSize)
	if err != nil {
		return nil, err
	}
	data := slice.Map[model.Article, domain.Article](res, func(idx int, val model.Article) domain.Article {
		return c.toDomain(val)
	})

	// 回写缓存
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if page == 1 && pageSize == 100 {
			err = c.cache.SetFirstPage(ctx, authorId, data)
			if err != nil {
				// TODO 记录日志
			}
		}
	}()

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		// 预加载
		if err = c.preCache(ctx, data); err != nil {
			// TODO 记录日志
			c.logger.Error("")
		}
	}()

	return data, nil
}

func (c *CacheArticleRepository) preCache(ctx context.Context, objs []domain.Article) error {
	if len(objs) <= 0 {
		return nil
	}
	return c.cache.Set(ctx, objs[0], time.Second*30)
}

func (c *CacheArticleRepository) toDomain(m model.Article) domain.Article {
	return domain.Article{
		Id:      m.Id,
		Title:   m.Title,
		Content: m.Content,
		Author: domain.Author{
			Id: m.AuthorId,
		},
		Status:     domain.ArticleStatus(m.Status),
		CreateTime: time.UnixMilli(m.CreateTime),
		UpdateTime: time.UnixMilli(m.UpdateTime),
	}
}

func (c *CacheArticleRepository) SyncStatus(ctx context.Context, authorId int64, id int64, status domain.ArticleStatus) error {
	defer func() {
		// 清空缓存
		c.cache.DelFirstPage(ctx, authorId)
	}()
	return c.dao.SyncStatus(ctx, authorId, id, status.ToUint8())
}

func (c *CacheArticleRepository) Sync(ctx context.Context, d domain.Article) (int64, error) {
	defer func() {
		// 清空缓存
		c.cache.DelFirstPage(ctx, d.Author.Id)
	}()
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
	defer func() {
		// 清空缓存
		c.cache.DelFirstPage(ctx, d.Author.Id)
	}()
	return c.dao.UpdateById(ctx, c.toModel(d))
}

func NewCacheArticleRepository(dao dao.ArticleDAO) *CacheArticleRepository {
	return &CacheArticleRepository{dao: dao}
}

func (c *CacheArticleRepository) Create(ctx context.Context, d domain.Article) (int64, error) {
	defer func() {
		c.cache.DelFirstPage(ctx, d.Author.Id)
	}()
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
