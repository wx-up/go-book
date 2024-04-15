package web

// VO view object 对标前端，包括请求和响应参数

type ArticleListReq struct {
	Page int64 `form:"page"`
	Size int64 `form:"size"`
}
type ArticleVO struct {
	Id    int64  `json:"id"`
	Title string `json:"title"`
	// 摘要
	Abstract   string `json:"abstract"`
	Content    string `json:"content"`
	AuthorId   int64  `json:"author_id"`
	AuthorName string `json:"author_name"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`

	// 状态转文字显示
	// 如果客户端是手机APP、小程序等涉及发版本的，可以考虑后端来处理，比如新增 StatusText 字段
	// 如果业务会发展到国际化，那么也可以考虑后端来处理
	Status uint8 `json:"status"`
}

type ArticleLikeReq struct {
}
