package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/ecodeclub/ekit/syncx/atomicx"

	"github.com/redis/go-redis/v9"

	"github.com/wx-up/go-book/internal/domain"
)

type ArticleRankingCache interface {
	Set(context.Context, []domain.Article) error
	Get(context.Context) ([]domain.Article, error)
	ForceGet(context.Context) ([]domain.Article, error)
}

type RedisArticleRankingCache struct {
	cmd        redis.Cmdable
	key        string
	expiration time.Duration
}

func (r *RedisArticleRankingCache) ForceGet(ctx context.Context) ([]domain.Article, error) {
	// TODO implement me
	panic("implement me")
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

// Set
// TODO：这里有一个预加载的优化方案
// 我们可以预期榜单的数据，用户点击的概率很大，因此在这里的时候，可以循环 articles 拿到 ID 缓存文章的信息
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
	return &RedisArticleRankingCache{
		cmd: cmd,
		key: "article:ranking",
	}
}

type LocalArticleRankingCache struct {
	topN *atomicx.Value[[]domain.Article]
	ddl  *atomicx.Value[time.Time]
	exp  time.Duration
}

// ForceGet 直接获取，不考虑缓存是否过期
func (l *LocalArticleRankingCache) ForceGet(ctx context.Context) ([]domain.Article, error) {
	arts := l.topN.Load()
	if len(arts) <= 0 {
		return nil, ErrCacheExpired
	}
	return arts, nil
}

func NewLocalArticleRankingCache() *LocalArticleRankingCache {
	return &LocalArticleRankingCache{
		topN: atomicx.NewValue[[]domain.Article](),
		ddl:  atomicx.NewValue[time.Time](),
		exp:  10 * time.Minute,
	}
}

type item struct {
	arts []domain.Article
	exp  time.Time
}

func (l *LocalArticleRankingCache) Set(_ context.Context, articles []domain.Article) error {
	// 因为这里是两个原则操作，因此存在并发问题
	// 有一个优化手段，新定义一个结构体 item 包含 article 和过期时间两个字段
	// 原子操作只操作 item 那么就只有一个原子操作了
	l.topN.Store(articles)
	l.ddl.Store(time.Now().Add(l.exp))
	return nil
}

var ErrCacheExpired = errors.New("cache expired")

func (l *LocalArticleRankingCache) Get(_ context.Context) ([]domain.Article, error) {
	ddl := l.ddl.Load()
	if ddl.IsZero() || ddl.Before(time.Now()) {
		return nil, ErrCacheExpired
	}
	arts := l.topN.Load()
	if len(arts) <= 0 {
		return nil, ErrCacheExpired
	}
	return arts, nil
}
