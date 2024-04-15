package model

type Interactive struct {
	Id int64 `gorm:"primaryKey;autoIncrement"`

	ReadCnt int64

	LikeCnt    int64 // 点赞数
	CollectCnt int64 // 收藏数

	CreateTime int64 `gorm:"index:idx_author_id_create_time;"`
	UpdateTime int64

	// biz + bizId 应该建立唯一索引
	// 这时候有一个问题，应该建立 biz+bizId 索引还是 bizId+biz 索引呢？
	// 推荐后者，因为 bizId 的区分度比较高，确定了 bizId 之后，很快就能确定 biz 进而确定记录
	// 如果后续有 where biz = xxx 的查询条件的话，那么推荐前者
	// 所以联合索引列的顺序如何确定：首先根据查询场景（ ESR ），其次再根据列的区分度。

	// 对于查询比较低频的场景，不用特意建立索引，因为索引也是有成本的
	BizId int64  `gorm:"uniqueIndex:biz_type_id"`
	Biz   string `gorm:"uniqueIndex:biz_type_id"`
}

type InteractiveV1 struct {
	Cnt     int64
	CntType int8
}

type InteractiveV2 struct {
	ReadCnt    int64
	LikeCnt    int64 // 点赞数
	CollectCnt int64 // 收藏数
}

// UserLikeBiz 点赞记录表
type UserLikeBiz struct {
	Id    int64  `gorm:"primaryKey;autoIncrement"`
	Uid   int64  `gorm:"uniqueIndex:uid_biz_id_biz"`
	BizId int64  `gorm:"uniqueIndex:uid_biz_id_biz"`
	Biz   string `gorm:"uniqueIndex:uid_biz_id_biz"`

	CreateTime int64
	UpdateTime int64

	// Status 软删除标记，它是一个存储状态，业务层面就是点赞和取消点赞
	// 这也说明了，domain 对象和 dao 对象状态不一定一一对应的
	Status int8 // 0 无效 1 有效
}

// UserCollectionBiz 收藏记录表
type UserCollectionBiz struct {
	Id int64 `gorm:"primaryKey;autoIncrement"`

	Cid   int64  `gorm:"uniqueIndex:uid_biz_id_collect"`
	Uid   int64  `gorm:"uniqueIndex:uid_biz_id_collect"`
	BizId int64  `gorm:"uniqueIndex:uid_biz_id_collect"`
	Biz   string `gorm:"uniqueIndex:uid_biz_id_collect"`

	Status int8 // 0 无效 1 有效

	// 收藏夹的ID
	// 收藏夹ID本身有索引
	UpdateTime int64
	CreateTime int64
}

// Collection 收藏夹
type Collection struct {
	Id   int64  `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"type:varchar(50);not null;default:'';unique"`
}
