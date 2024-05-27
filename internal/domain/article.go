package domain

import (
	"time"
)

type Article struct {
	Id         int64
	Title      string
	Content    string
	Author     Author
	Status     ArticleStatus
	CreateTime time.Time
	UpdateTime time.Time
}

type Author struct {
	Id   int64
	Name string
}

// Abstract 文章摘要
// 如果文章内容是保存到 OSS 的话，需要在创建文章的时候生成一份摘要存储到数据库中
// 然后从数据库中直接读取摘要
func (a Article) Abstract() string {
	// 摘要
	// 方案一：截取文章内容的前几个字符
	// 方案二：使用AI来生成
	cs := []rune(a.Content)
	if len(cs) > 100 {
		return string(cs[:100])
	}
	return a.Content
}

// ArticleStatus 定义衍生类型
// 衍生类型的好处就是可以在上面定义一些方法
type ArticleStatus uint8

const (
	// ArticleStatusUnknown 为了避免零值问题，
	// 因为有时候如果前端传递了 status 字段的话，如果把零值当成有意义的状态，你就区分不出来
	// 前端是传递了 status=0 还是没有传递 status 值
	ArticleStatusUnknown     ArticleStatus = iota // 未知
	ArticleStatusUnpublished                      // 未发表
	ArticleStatusPublished                        // 已发表
	ArticleStatusPrivate                          // 仅自己可见
)

func (s ArticleStatus) ToUint8() uint8 {
	return uint8(s)
}

func (s ArticleStatus) Valid() bool {
	return s > 0 && s < 4
}

func (s ArticleStatus) String() string {
	switch s {
	case ArticleStatusPublished:
		return "已发表"
	case ArticleStatusUnpublished:
		return "未发表"
	case ArticleStatusPrivate:
		return "仅自己可见"
	default:
		return "未知"
	}
}

type ArticleStatusV1 struct {
	Val  uint8
	Name string
}

var ArticleStatusV1Unknown = ArticleStatusV1{
	Val:  0,
	Name: "未知",
}
