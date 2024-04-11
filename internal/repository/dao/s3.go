package dao

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/wx-up/go-book/internal/domain"

	"github.com/spf13/cast"
	"github.com/wx-up/go-book/pkg"

	"gorm.io/gorm/clause"

	"github.com/aws/aws-sdk-go-v2/service/s3"

	"gorm.io/gorm"

	"github.com/wx-up/go-book/internal/repository/dao/model"
)

type S3ArticleDAO struct {
	dao ArticleDAO

	p      DBProvider
	Oss    *s3.Client
	bucket string
}

func NewS3Dao(p DBProvider, oss *s3.Client) *S3ArticleDAO {
	return &S3ArticleDAO{
		dao:    NewGORMArticleDAO(p),
		p:      p,
		Oss:    oss,
		bucket: "public_articles",
	}
}

func (s *S3ArticleDAO) Insert(ctx context.Context, article model.Article) (int64, error) {
	return s.dao.Insert(ctx, article)
}

func (s *S3ArticleDAO) UpdateById(ctx context.Context, article model.Article) error {
	return s.dao.UpdateById(ctx, article)
}

func (s *S3ArticleDAO) Sync(ctx context.Context, article model.Article) (int64, error) {
	var (
		id      = article.Id
		err     error
		content = article.Content
	)
	err = s.p().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
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

		article.Id = id
		now := time.Now()
		article.UpdateTime = now.UnixMilli()
		article.CreateTime = now.UnixMilli()

		// 内容置为空，不存在数据库中，保存到 OSS 中
		article.Content = ""
		publishArticle := model.PublishArticle{
			Article: article,
		}
		return tx.Model(&model.PublishArticle{}).Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"update_time": publishArticle.UpdateTime,
				"title":       publishArticle.Title,
			}),
		}).Create(&publishArticle).Error
	})
	if err != nil {
		return 0, err
	}

	// 保存到 OSS 中
	// 这里需要有监控、有重试，因为同步到OSS的调用不在事务中，有可能部分失败导致数据不一致
	_, err = s.Oss.PutObject(ctx, &s3.PutObjectInput{
		Bucket: pkg.ToPtr[string](s.bucket),
		Key:    pkg.ToPtr[string](cast.ToString(id)),
		Body:   bytes.NewReader([]byte(content)),
		// 需要设置 Content-Type 否则 body 有中文的话会乱码
		ContentType: pkg.ToPtr[string]("text/plain;charset=utf-8"),
	})
	return id, err
}

func (s *S3ArticleDAO) Transaction(ctx context.Context, f func(dao ArticleDAO, readerDao ReaderArticleDAO) error) error {
	// TODO implement me
	panic("implement me")
}

func (s *S3ArticleDAO) SyncStatus(ctx context.Context, uid int64, id int64, status uint8) error {
	now := time.Now().UnixMilli()
	err := s.p().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&model.Article{}).Where("id = ? and author_id = ?", id, uid).Updates(map[string]any{
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
	if err != nil {
		return err
	}

	// 删除 OSS 数据
	if status == domain.ArticleStatusPrivate.ToUint8() {
		_, err = s.Oss.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: pkg.ToPtr[string](s.bucket),
			Key:    pkg.ToPtr[string](cast.ToString(id)),
		})
	}
	return err
}
