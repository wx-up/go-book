package dao

import (
	"context"
	"github.com/wx-up/go-book/interactive/repository/dao/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

var ErrRecordNotFound = gorm.ErrRecordNotFound

type InteractiveDAO interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	BatchIncrReadCnt(ctx context.Context, bizs []string, bizIds []int64) error
	InsertLikeInfo(ctx context.Context, biz string, id int64, uid int64) error
	DeleteLikeInfo(ctx context.Context, biz string, id int64, uid int64) error
	InsertCollectionBiz(ctx context.Context, cb UserCollectionBiz) error
	GetLikeInfo(ctx context.Context,
		biz string, id int64, uid int64) (UserLikeBiz, error)
	GetCollectInfo(ctx context.Context,
		biz string, id int64, uid int64) (UserCollectionBiz, error)
	Get(ctx context.Context, biz string, id int64) (Interactive, error)
	GetByIds(ctx context.Context, biz string, ids []int64) ([]Interactive, error)
}

type GORMInteractiveDAO struct {
	db *gorm.DB
}

func (dao *GORMInteractiveDAO) GetByIds(ctx context.Context, biz string, ids []int64) ([]Interactive, error) {
	var res []Interactive
	err := dao.db.WithContext(ctx).
		Where("biz = ? AND biz_id IN ?", biz, ids).
		First(&res).Error
	return res, err
}

func (dao *GORMInteractiveDAO) Get(ctx context.Context, biz string, id int64) (Interactive, error) {
	var res Interactive
	err := dao.db.WithContext(ctx).
		Where("biz = ? AND biz_id = ?", biz, id).
		First(&res).Error
	return res, err
}

func (dao *GORMInteractiveDAO) GetLikeInfo(ctx context.Context,
	biz string, id int64, uid int64) (UserLikeBiz, error) {
	var res UserLikeBiz
	err := dao.db.WithContext(ctx).
		Where("biz = ? AND biz_id = ? AND uid = ? AND status = ?",
			biz, id, uid, 1).
		First(&res).Error
	return res, err
}

func (dao *GORMInteractiveDAO) GetCollectInfo(ctx context.Context,
	biz string, id int64, uid int64) (UserCollectionBiz, error) {
	var res UserCollectionBiz
	err := dao.db.WithContext(ctx).
		Where("biz = ? AND biz_id = ? AND uid = ?", biz, id, uid).
		First(&res).Error
	return res, err
}

func (dao *GORMInteractiveDAO) InsertCollectionBiz(ctx context.Context,
	cb UserCollectionBiz) error {
	now := time.Now().UnixMilli()
	cb.CreateTime = now
	cb.UpdateTime = now
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&cb).Error
		if err != nil {
			return err
		}
		return tx.WithContext(ctx).Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"collect_cnt": gorm.Expr("`collect_cnt` + 1"),
				"utime":       now,
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

func (dao *GORMInteractiveDAO) InsertLikeInfo(ctx context.Context,
	biz string, id int64, uid int64) error {
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"update_time": now,
				"status":      1,
			}),
		}).Create(&model.UserLikeBiz{
			Uid:        uid,
			Biz:        biz,
			BizId:      id,
			Status:     1,
			UpdateTime: now,
			CreateTime: now,
		}).Error
		if err != nil {
			return err
		}
		return tx.WithContext(ctx).Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"like_cnt":    gorm.Expr("`like_cnt` + 1"),
				"update_time": now,
			}),
		}).Create(&model.Interactive{
			Biz:        biz,
			BizId:      id,
			LikeCnt:    1,
			CreateTime: now,
			UpdateTime: now,
		}).Error
	})
}

func (dao *GORMInteractiveDAO) DeleteLikeInfo(ctx context.Context,
	biz string, id int64, uid int64) error {
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&UserLikeBiz{}).
			Where("uid=? AND biz_id = ? AND biz=?", uid, id, biz).
			Updates(map[string]interface{}{
				"update_time": now,
				"status":      0,
			}).Error
		if err != nil {
			return err
		}
		return tx.Model(&model.Interactive{}).
			Where("biz =? AND biz_id=?", biz, id).
			Updates(map[string]interface{}{
				"like_cnt":    gorm.Expr("`like_cnt` - 1"),
				"update_time": now,
			}).Error
	})
}

func NewGORMInteractiveDAO(db *gorm.DB) InteractiveDAO {
	return &GORMInteractiveDAO{db: db}
}

func (dao *GORMInteractiveDAO) BatchIncrReadCnt(ctx context.Context, bizs []string, bizIds []int64) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txDAO := NewGORMInteractiveDAO(tx)
		for i := 0; i < len(bizs); i++ {
			err := txDAO.IncrReadCnt(ctx, bizs[i], bizIds[i])
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (dao *GORMInteractiveDAO) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{
			"read_cnt":    gorm.Expr("`read_cnt` + 1"),
			"update_time": now,
		}),
	}).Create(&model.Interactive{
		Biz:        biz,
		BizId:      bizId,
		ReadCnt:    1,
		CreateTime: now,
		UpdateTime: now,
	}).Error
}

type UserLikeBiz struct {
	model.UserLikeBiz
}

type UserCollectionBiz struct {
	model.UserCollectionBiz
}

type Interactive struct {
	model.Interactive
}
