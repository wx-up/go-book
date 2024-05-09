package dao

import (
	"context"
	"time"

	"gorm.io/gorm/clause"

	"github.com/wx-up/go-book/internal/repository/dao/model"
)

// ReaderArticleDAO 线上库DAO，针对读者
type ReaderArticleDAO interface {
	// Upsert 有就更新，没有就插入
	Upsert(ctx context.Context, obj model.PublishArticle) (int64, error)

	// GetById 获取线上库的文章详情
	GetById(ctx context.Context, id int64) (model.PublishArticle, error)

	ListPub(ctx context.Context, startTime time.Time, offset, limit int) ([]model.PublishArticle, error)
}

type GORMReaderArticleDAO struct {
	p DBProvider
}

func (g *GORMReaderArticleDAO) ListPub(ctx context.Context, startTime time.Time, offset, limit int) ([]model.PublishArticle, error) {
	res := make([]model.PublishArticle, 0, limit)
	err := g.p().WithContext(ctx).
		Where("update_time < ?", startTime.UnixMilli()).
		Limit(limit).
		Offset(offset).
		Order("update_time DESC").
		Find(&res).Error
	return res, err
}

func (g *GORMReaderArticleDAO) GetById(ctx context.Context, id int64) (model.PublishArticle, error) {
	// TODO implement me
	panic("implement me")
}

func NewGORMReaderArticleDAO(p DBProvider) *GORMReaderArticleDAO {
	return &GORMReaderArticleDAO{
		p: p,
	}
}

func (g *GORMReaderArticleDAO) Upsert(ctx context.Context, obj model.PublishArticle) (int64, error) {
	now := time.Now()
	obj.UpdateTime = now.UnixMilli()
	obj.CreateTime = now.UnixMilli()
	// 在 `id` 冲突时，更新 `update_time` 和 `title` 和 `content` 列，反之插入新行
	err := g.p().WithContext(ctx).Clauses(clause.OnConflict{
		// SQL2003的规范， sqlite 这种是符合该规范的，支持指定列的数据库，可以指定冲突的列
		// INSERT INTO xxx ON CONFLICT (xxx) DO UPDATE xxx WHERE xxx
		// INSERT INTO xxx ON CONFLICT (xxx) DO NOTHING
		// mysql 不遵循该规范，只支持如下格式的语句，不支持指定列，它只能根据主键、唯一索引冲突做 upsert 也不支持 where 条件的
		// INSERT INTO xxx ON DUPLICATE KEY UPDATE xxx

		// 指定冲突列
		// Columns: []clause.Column{{Name: "id"}},

		// 冲突了啥也不干
		// DoNothing: true,

		// 冲突了，并且符合 where 条件的才会更新
		// mysql 也不支持 where 条件的
		// Where:

		// AssignmentColumns 指定更新的列
		// Assignments 指定更新的列和值
		DoUpdates: clause.Assignments(map[string]interface{}{
			"update_time": obj.UpdateTime,
			"title":       obj.Title,
			"content":     obj.Content,
		}),
	}).Create(&obj).Error
	return obj.Id, err
}
