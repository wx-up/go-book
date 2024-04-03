package dao

import (
	"context"
	"errors"
	"time"

	"github.com/wx-up/go-book/internal/repository/dao/model"
)

var ErrArticleNotFound = errors.New("article not found")

type ArticleDAO interface {
	Insert(ctx context.Context, article model.Article) (int64, error)
	UpdateById(ctx context.Context, article model.Article) error
}

type GORMArticleDAO struct {
	p DBProvider
}

func (a *GORMArticleDAO) UpdateById(ctx context.Context, article model.Article) error {
	// 不推荐下面的方式更新，它依赖gorm忽略零值的特性，可读性很差，如果不熟悉gorm你都不知道它会更新哪些字段
	// 当然它也有好处，如果需要更新新的字段，这种写法不需要改什么
	// a.p().WithContext(ctx).Where("id = ?", article.Id).Updates(article)
	// 指定字段更新的话，当需要更新新的字段时，map中要添加新的字段
	res := a.p().WithContext(ctx).Model(&model.Article{}).
		Where("id =? and author_id =?", article.Id, article.AuthorId).
		Updates(map[string]any{
			"title":       article.Title,
			"content":     article.Content,
			"update_time": time.Now().UnixMilli(),
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected <= 0 {
		return ErrArticleNotFound
	}
	return nil
}

func (a *GORMArticleDAO) Insert(ctx context.Context, article model.Article) (int64, error) {
	t := time.Now().UnixMilli()
	article.CreateTime = t
	article.UpdateTime = t
	err := a.p().WithContext(ctx).Create(&article).Error
	return article.Id, err
}

func NewGORMArticleDAO(p DBProvider) *GORMArticleDAO {
	return &GORMArticleDAO{
		p: p,
	}
}
