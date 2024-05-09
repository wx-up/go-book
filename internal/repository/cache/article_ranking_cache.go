package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/wx-up/go-book/internal/domain"
)

type ArticleRankingCache interface {
	Set(context.Context, []domain.Article) error
	Get(context.Context) ([]domain.Article, error)
}

type RedisArticleRankingCache struct {
	cmd        redis.Cmdable
	key        string
	expiration time.Duration
}

func (r *RedisArticleRankingCache) Get(ctx context.Context) ([]domain.Article, error) {
	bs, err := r.cmd.Get(ctx, r.key).Bytes()
	if err != nil {
		return nil, err
	}
	var articles []domain.Article
	if err = json.Unmarshal(bs, &articles); err != nil {
		return nil, err
	}
	return articles, nil
}

func (r *RedisArticleRankingCache) Set(ctx context.Context, articles []domain.Article) error {
	// 避免在缓存中存储大字段
	for i, article := range articles {
		article.Content = article.Abstract()
		articles[i] = article
	}
	bs, err := json.Marshal(articles)
	if err != nil {
		return err
	}

	// r.expiration 需要久一点，最少为一次热榜的计算时间（ 包括重试 ）
	// 其实不设置有效期也是可以的，这样就保证数据库有问题时（ 来不及刷新榜单数据 ）
	// 榜单接口仍旧可用
	return r.cmd.Set(ctx, r.key, bs, r.expiration).Err()
}

func NewRedisArticleRankingCache(cmd redis.Cmdable) *RedisArticleRankingCache {
	return &RedisArticleRankingCache{cmd: cmd}
}
