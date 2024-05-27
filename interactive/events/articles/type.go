package articles

const TopicReadEvent = "article_read"

type ReadEvent struct {
	Uid int64
	Aid int64
}
