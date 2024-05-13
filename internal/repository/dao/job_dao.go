package dao

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/wx-up/go-book/internal/repository/dao/model"
)

var ErrPreemptFailed = errors.New("preempt failed")

type JobDAO interface {
	Preempt(ctx context.Context) (model.Job, error)
	Release(ctx context.Context, obj model.Job) error
	UpdateNextTime(ctx context.Context, id int64, nextTime time.Time) error
	UpdateTime(ctx context.Context, id int64) error
	Stop(ctx context.Context, id int64) error
}

type GORMJobDAO struct {
	db *gorm.DB
}

func (dao *GORMJobDAO) UpdateNextTime(ctx context.Context, id int64, nextTime time.Time) error {
	return dao.db.WithContext(ctx).Model(&model.Job{}).Where("id = ?", id).Updates(map[string]any{
		"next_time": nextTime.UnixMilli(),
	}).Error
}

func (dao *GORMJobDAO) UpdateTime(ctx context.Context, id int64) error {
	return dao.db.WithContext(ctx).Model(&model.Job{}).Where("id = ?", id).Updates(map[string]any{
		"update_time": time.Now().UnixMilli(),
	}).Error
}

func (dao *GORMJobDAO) Stop(ctx context.Context, id int64) error {
	return dao.db.WithContext(ctx).Model(&model.Job{}).Where("id = ?", id).Updates(map[string]any{
		"status":      model.JobStatusPaused,
		"update_time": time.Now().UnixMilli(),
	}).Error
}

func (dao *GORMJobDAO) Release(ctx context.Context, obj model.Job) error {
	// where 条件需要增加 version 的判断
	// 实例A抢占了任务A，在执行，但是它和数据库失去的联系，导致过段时间实例B也抢占了任务A
	// 这时候实例A重新取得了与数据库的联系，如果不加版本号的筛选，就会释放实例B抢占的任务A
	return dao.db.WithContext(ctx).
		Model(&model.Job{}).Where("id = ? AND version = ?", obj.Id, obj.Version).Updates(map[string]any{
		"status":      model.JobStatusWaiting,
		"update_time": time.Now().UnixMilli(),
	}).Error
}

func (dao *GORMJobDAO) Preempt(ctx context.Context) (model.Job, error) {
	// 高并发的情况下，下面的代码会有问题的
	// 比如 100 个高并发的请求同时抢占，但是只有一个能抢到，其他的请求都抢不到，继续下一次抢占
	// 但是实际场景，抢占的 goroutine 数量会比较少，所以影响不大
	// 执行 job 的 goroutine 可以设置的比较多

	// 如果一定要优化的话可以使用下面的方案：
	//   1. 一次取出100条，然后从中随机获取一条，再执行下面的 updates 语句
	//   2. 搞一个随机偏移量（ 比如 0-100 ），在查询的时候添加 offset 语句，如果查询不到数据，则兜底，将 offset = 0 再次查询
	//   3. 搞一个随机取余数的条件 status = ? AND next_time <= ? AND id%10 = 0 其中 10 是随机的（ 比如1-10随机 ）
	//  上面的三个方案，在执行一段时间之后，可以把条件去掉，直接走查询 status = ? AND next_time <= ? 然后 order by next_time asc limit 1
	//  取最早的那条抢占，这样可以解决随机不均匀，导致一些记录一直没有被抢占的问题
	for {
		now := time.Now()
		var obj model.Job
		err := dao.db.WithContext(ctx).Where("status = ? AND next_time <= ?",
			model.JobStatusWaiting,
			now.UnixMilli(),
		).First(&obj).Error
		// 记录不存在或者数据库错误
		if err != nil {
			return model.Job{}, err
		}

		// CAS 操作 Compare AND Swap
		// 乐观锁
		// 一个很常见的优化手段：用乐观锁取代 for update
		// for update 容易导致死锁问题，参考文章：https://blog.csdn.net/zhouwenjun0820/article/details/108790922
		res := dao.db.WithContext(ctx).Model(&model.Job{}).
			Where("id = ? AND version = ?", obj.Id, obj.Version).
			Updates(map[string]any{
				"status":      model.JobStatusPreempted,
				"update_time": now.UnixMilli(),
				"version":     obj.Version + 1,
			})
		if res.Error != nil {
			return model.Job{}, res.Error
		}

		// 抢占失败，则继续尝试抢占
		if res.RowsAffected <= 0 {
			continue
		}

		return obj, nil
	}
}
