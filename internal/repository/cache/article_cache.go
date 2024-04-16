package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/wx-up/go-book/internal/domain"
)

type ArticleCache interface {
	GetFirstPage(ctx context.Context, authorId int64) ([]domain.Article, error)
	SetFirstPage(ctx context.Context, authorId int64, articles []domain.Article) error
	DelFirstPage(ctx context.Context, authorId int64) error

	Set(ctx context.Context, obj domain.Article, expire time.Duration) error

	SetPub(ctx context.Context, obj domain.Article, expire time.Duration) error
}

type RedisArticleCache struct {
	client redis.Cmdable
}

func (r *RedisArticleCache) SetPub(ctx context.Context, obj domain.Article, expire time.Duration) error {
	bs, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.pubKey(obj.Id), bs, expire).Err()
}

func (r *RedisArticleCache) Set(ctx context.Context, obj domain.Article, expire time.Duration) error {
	bs, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.key(obj.Id), bs, expire).Err()
}

func (r *RedisArticleCache) DelFirstPage(ctx context.Context, authorId int64) error {
	// TODO implement me
	panic("implement me")
}

func (r *RedisArticleCache) SetFirstPage(ctx context.Context, authorId int64, articles []domain.Article) error {
	// 不需要将 content 缓存，因为列表页也不需要，其次比较大还占内存
	for index := range articles {
		articles[index].Content = ""
	}
	bs, err := json.Marshal(articles)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.firstPageKey(authorId), bs, time.Minute*10).Err()
}

func (r *RedisArticleCache) firstPageKey(authorId int64) string {
	return fmt.Sprintf("article:first_page:%d", authorId)
}

func (r *RedisArticleCache) pubKey(id int64) string {
	return fmt.Sprintf("article:pub:detail:%d", id)
}

func (r *RedisArticleCache) key(id int64) string {
	return fmt.Sprintf("article:detail:%d", id)
}

func NewRedisArticleCache(client redis.Cmdable) *RedisArticleCache {
	return &RedisArticleCache{client: client}
}

func (r *RedisArticleCache) GetFirstPage(ctx context.Context, authorId int64) ([]domain.Article, error) {
	// TODO implement me
	panic("implement me")
}
