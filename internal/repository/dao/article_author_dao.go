package dao

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/wx-up/go-book/internal/repository/dao/model"
)

var ErrArticleNotFound = errors.New("article not found")

// ArticleDAO 制作库DAO，针对创作者
type ArticleDAO interface {
	Insert(ctx context.Context, article model.Article) (int64, error)
	UpdateById(ctx context.Context, article model.Article) error
	Sync(ctx context.Context, article model.Article) (int64, error)
	Transaction(ctx context.Context, f func(dao ArticleDAO, readerDao ReaderArticleDAO) error) error
	SyncStatus(ctx context.Context, uid int64, id int64, status uint8) error
}

type GORMArticleDAO struct {
	p DBProvider
}

func (a *GORMArticleDAO) SyncStatus(ctx context.Context, uid int64, id int64, status uint8) error {
	now := time.Now().UnixMilli()
	return a.p().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&model.Article{}).Where("id = ? and author_id = ?", uid, id).Updates(map[string]any{
			"status":      status,
			"update_time": now,
		})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			// 可能有人在搞，打点
			return fmt.Errorf("同步文章状态异常，作者id：%d，文章id：%d", uid, id)
		}
		return tx.Model(&model.PublishArticle{}).Where("id = ?", id).Updates(map[string]any{
			"status":      status,
			"update_time": now,
		}).Error
	})
}

func (a *GORMArticleDAO) Transaction(ctx context.Context, f func(dao ArticleDAO, readerDao ReaderArticleDAO) error) error {
	return a.p().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txDAO := NewGORMArticleDAO(func() *gorm.DB {
			return tx
		})
		readerTxDAO := NewGORMReaderArticleDAO(func() *gorm.DB {
			return tx
		})
		return f(txDAO, readerTxDAO)
	})
}

func (a *GORMArticleDAO) Sync(ctx context.Context, article model.Article) (int64, error) {
	var (
		id  = article.Id
		err error
	)
	err = a.p().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txDAO := NewGORMArticleDAO(func() *gorm.DB {
			return tx
		})
		if id > 0 {
			err = txDAO.UpdateById(ctx, article)
		} else {
			id, err = txDAO.Insert(ctx, article)
		}
		if err != nil {
			return err
		}

		// 同步给线上库
		id, err = NewGORMReaderArticleDAO(func() *gorm.DB {
			return tx
		}).Upsert(ctx, model.PublishArticle{Article: article})
		return err
	})
	return id, err
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
