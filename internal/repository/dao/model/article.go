package model

// Article 制作库的文章表
// 设计索引可以从业务角度出发，高频的 SQL 查询 Where 条件是什么
// 制作库的文章表，从业务角度高频的操作就是打开草稿箱或者打开已经发布的文章，里面都是我的文章
// 因此高频的查询就是 where author_id = ? 所以可以考虑建立 author_id 索引
// 其次文章列表往往是按照时间排序
// 所以总结下来可以设计一个 author_id + create_time 联合索引
type Article struct {
	Id int64 `gorm:"primaryKey;autoIncrement"`

	// 标题、content 一般只是用于存储的，并不会 like 等进行查询，如果要查询的话，得用 ES 等工具
	Title   string `gorm:"type:varchar(100);not null"`
	Content string `gorm:"type:text"`
	// 时间戳，单位毫秒
	AuthorId   int64 `gorm:"index:idx_author_id_create_time;"`
	CreateTime int64 `gorm:"index:idx_author_id_create_time;"`
	UpdateTime int64
}
