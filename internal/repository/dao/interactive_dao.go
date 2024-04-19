package dao

import (
	"context"
	"time"

	"github.com/wx-up/go-book/internal/repository/dao/model"

	"gorm.io/gorm/clause"

	"gorm.io/gorm"
)

var (
	ErrInteractiveLikedNotFound     = gorm.ErrRecordNotFound
	ErrInteractiveCollectedNotFound = gorm.ErrRecordNotFound
)

type InteractiveDao interface {
	IncrReadCnt(ctx context.Context, biz string, bid int64) error
	BatchIncrReadCnt(ctx context.Context, biz string, bid []int64) error

	InsertLikeInfo(ctx context.Context, biz string, bid int64, uid int64) error
	DelLikeInfo(ctx context.Context, biz string, bid int64, uid int64) error

	InsertCollectionInfo(ctx context.Context, cbo model.UserCollectionBiz) error
	DelCollectionInfo(ctx context.Context, cbo model.UserCollectionBiz) error

	GetLikeInfo(ctx context.Context, info model.UserLikeBiz) (model.UserLikeBiz, error)
	GetCollectionInfo(ctx context.Context, info model.UserCollectionBiz) (model.UserCollectionBiz, error)

	Get(ctx context.Context, biz string, bid int64) (model.Interactive, error)
}

type GORMInteractiveDao struct {
	db *gorm.DB
}

// BatchIncrReadCnt
// 批处理增加阅读计数
// 虽然事务里面还是for循环更新记录，但是事务的次数只有一次
// 如果没有批操作的话，需要for循环更新1000条记录，事务的次数就会是1000
// 事务操作在mysql中也是比较重的操作
func (g *GORMInteractiveDao) BatchIncrReadCnt(ctx context.Context, biz string, bid []int64) error {
	// 可以进一步使用 map 优化
	// 以 biz+bid 作为 key，出现的次数作为 value
	// 计数直接增加 value 就可以，不用一次一次的增加
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		dao := NewGORMInteractiveDao(tx)
		for _, b := range bid {
			err := dao.IncrReadCnt(ctx, biz, b)
			if err != nil {
				// 阅读计数失败就失败，无关紧要，不用回滚事务
				// 记录日志
			}
		}
		return nil
	})
}

func (g *GORMInteractiveDao) Get(ctx context.Context, biz string, bid int64) (model.Interactive, error) {
	// TODO implement me
	panic("implement me")
}

func (g *GORMInteractiveDao) GetLikeInfo(ctx context.Context, info model.UserLikeBiz) (model.UserLikeBiz, error) {
	var res model.UserLikeBiz
	err := g.db.WithContext(ctx).Where("biz = ? AND biz_id = ? AND uid = ? AND status = 1", info.Biz, info.BizId, info.Uid).First(&res).Error
	return res, err
}

func (g *GORMInteractiveDao) GetCollectionInfo(ctx context.Context, info model.UserCollectionBiz) (model.UserCollectionBiz, error) {
	var res model.UserCollectionBiz
	err := g.db.WithContext(ctx).Where("biz = ? AND biz_id = ? AND uid = ? AND status = 1", info.Biz, info.BizId, info.Uid).First(&res).Error
	return res, err
}

func (g *GORMInteractiveDao) DelCollectionInfo(ctx context.Context, cbo model.UserCollectionBiz) error {
	now := time.Now().UnixMilli()
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&model.UserLikeBiz{}).
			Where("biz = ? AND biz_id = ? AND uid = ? AND cid = ?", cbo.Biz, cbo.BizId, cbo.Uid, cbo.Cid).Updates(map[string]any{
			"status":      0, // 0 表示删除
			"update_time": now,
		}).Error
		if err != nil {
			return err
		}
		return tx.Model(&model.Interactive{}).Where("biz = ? AND biz_id = ?", cbo.Biz, cbo.BizId).Updates(map[string]any{
			"collect_cnt": gorm.Expr("collect_cnt - ?", 1),
			"update_time": now,
		}).Error
	})
}

func (g *GORMInteractiveDao) InsertCollectionInfo(ctx context.Context, cb model.UserCollectionBiz) error {
	now := time.Now().UnixMilli()
	cb.CreateTime = now
	cb.UpdateTime = now
	cb.Status = 1
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"update_time": now,
				"status":      1,
			}),
		}).Create(&cb).Error
		if err != nil {
			return err
		}
		return tx.WithContext(ctx).Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"collect_cnt": gorm.Expr("`collect_cnt` + 1"),
				"update_time": now,
			}),
		}).Create(&model.Interactive{
			Biz:        cb.Biz,
			BizId:      cb.BizId,
			CollectCnt: 1,
			CreateTime: now,
			UpdateTime: now,
		}).Error
	})
}

func (g *GORMInteractiveDao) DelLikeInfo(ctx context.Context, biz string, bid int64, uid int64) error {
	now := time.Now().UnixMilli()
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 用户点赞记录标记成软删除
		err := tx.Model(&model.UserLikeBiz{}).
			Where("biz = ? AND biz_id = ? AND uid = ?", biz, bid, uid).Updates(map[string]any{
			"status":      0, // 0 表示删除
			"update_time": now,
		}).Error
		if err != nil {
			return err
		}
		// 减少点赞数
		return tx.Model(&model.Interactive{}).Where("biz = ? AND biz_id = ?", biz, bid).Updates(map[string]any{
			"like_cnt":    gorm.Expr("like_cnt - ?", 1),
			"update_time": now,
		}).Error
	})
}

// InsertLikeInfo 增加点赞数
func (g *GORMInteractiveDao) InsertLikeInfo(ctx context.Context, biz string, bid int64, uid int64) error {
	now := time.Now().UnixMilli()
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 用户点赞记录表
		// 用户没有点赞的场景下，表记录为空或者存在一条记录但是status=0
		// 所以这里使用 upsert 语义
		err := tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"update_time": now,
				"status":      1,
			}),
		}).Create(&model.UserLikeBiz{
			Uid:        uid,
			BizId:      bid,
			Biz:        biz,
			CreateTime: now,
			UpdateTime: now,
			Status:     1,
		}).Error
		if err != nil {
			return err
		}

		// 增加点赞数
		return tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"like_cnt":    gorm.Expr("like_cnt + ?", 1),
				"update_time": now,
			}),
		}).Create(&model.Interactive{
			LikeCnt:    1,
			CreateTime: now,
			UpdateTime: now,
			Biz:        biz,
			BizId:      bid,
		}).Error
	})
}

func (g *GORMInteractiveDao) IncrReadCnt(ctx context.Context, biz string, bid int64) error {
	now := time.Now().UnixMilli()
	// withContext 控制事务超时
	return g.db.WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{
			"read_cnt":    gorm.Expr("read_cnt + ?", 1),
			"update_time": now,
		}),
	}).Create(&model.Interactive{
		ReadCnt:    1,
		CreateTime: now,
		UpdateTime: now,
		Biz:        biz,
		BizId:      bid,
	}).Error
}

func NewGORMInteractiveDao(db *gorm.DB) *GORMInteractiveDao {
	return &GORMInteractiveDao{db: db}
}
