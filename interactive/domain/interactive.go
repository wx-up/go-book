package domain

type Interactive struct {
	BizId      int64
	ReadCnt    int64
	LikeCnt    int64
	CollectCnt int64
	Liked      bool // 是否点赞
	Collected  bool // 是否收藏
}
