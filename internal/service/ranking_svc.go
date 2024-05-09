package service

import (
	"context"
	"math"
	"time"

	"github.com/wx-up/go-book/internal/repository"

	"github.com/ecodeclub/ekit/queue"

	"github.com/wx-up/go-book/internal/domain"
	"github.com/wx-up/go-book/pkg/slice"
)

// RankingService 排名服务
type RankingService interface {
	// TopN 默认取排名100的文章
	TopN(ctx context.Context) error

	// 接口定义上可以将 N 的值由调用来指定
	// TopN(ctx context.Context, n int64) error
	// 因为 top 的文章数据是存储在 redis 中，因此可以不返回 []domain.Article 但是为了方便测试，还是可以考虑返回 []domain.Article
	// TopN(ctx context.Context, n int64) ([]domain.Article, error)
}

// ArticleRankingService 文章排名
type ArticleRankingService struct {
	// 单体服务可以考虑组合 repository 而不是 svc
	// 这里将 BatchRankingService 当作聚合服务来定位，因此组合了 svc
	articleSvc     ArticleService
	interactiveSvc InteractiveService

	batchSize int64
	n         int64

	scoreFunc func(t time.Time, likeCnt int64) float64

	repo repository.ArticleRankingRepo
}

func NewArticleRankingService(
	articleSvc ArticleService,
	interactiveSvc InteractiveService,
	repo repository.ArticleRankingRepo,
) *ArticleRankingService {
	return &ArticleRankingService{
		articleSvc:     articleSvc,
		interactiveSvc: interactiveSvc,
		batchSize:      100,
		n:              100,
		repo:           repo,
		// hacknews 模型
		scoreFunc: func(t time.Time, likeCnt int64) float64 {
			return float64(likeCnt+1) / math.Pow(time.Since(t).Seconds()+2, 1.5)
		},
	}
}

// TopN 直接测试它不好测，因为它没有返回值
// 因此可以考虑，拆分出一个方法，比如 topN 它有返回值，然后先测试它，再测试 TopN
// 原则：大方法不好测试时，将大方法拆分成小方法进行测试（ 分而治之 ）
func (b *ArticleRankingService) TopN(ctx context.Context) error {
	ids, err := b.topN(ctx)
	if err != nil {
		return err
	}
	_ = ids
	return b.repo.Set(ctx, nil)
}

// topN 对于这种复杂的函数，可以考虑 TDD，测试驱动编写
func (b *ArticleRankingService) topN(ctx context.Context) ([]int64, error) {
	offset := int64(0)
	type Score struct {
		art   domain.Article
		score float64
	}

	// NewPriorityQueue 优先级队列
	topN := queue.NewPriorityQueue[Score](int(b.n), func(a, b Score) int {
		if a.score > b.score {
			return 1
		} else if a.score < b.score {
			return -1
		} else {
			return 0
		}
	})

	// 全量的文章计算热榜
	// 这里其实考虑业务折中，比如7天之前的文章不进入文章热榜计算
	now := time.Now()
	for {
		arts, err := b.articleSvc.ListPub(ctx, now, offset, b.batchSize)
		if err != nil {
			return nil, err
		}
		ids := slice.Map[domain.Article, int64](arts, func(idx int, a domain.Article) int64 {
			return a.Id
		})

		inters, err := b.interactiveSvc.GetByIds(ctx, "article", ids)
		if err != nil {
			return nil, err
		}

		// 转成 map 方便查找
		idInterMap := make(map[int64]domain.Interactive)
		for _, inter := range inters {
			idInterMap[inter.BizId] = inter
		}

		for _, art := range arts {
			inter, ok := idInterMap[art.Id]
			if !ok {
				continue
			}

			score := b.scoreFunc(art.UpdateTime, inter.LikeCnt)

			err := topN.Enqueue(Score{
				art:   art,
				score: score,
			})

			// 表示优先级队列满了，这时候需要淘汰
			if err == queue.ErrOutOfCapacity {
				v, _ := topN.Peek()
				if v.score < score {
					_, _ = topN.Dequeue()
					_ = topN.Enqueue(Score{
						art:   art,
						score: score,
					})
				}
			}
		}

		// 一批没有取够，说明已经取完了
		// 如果业务折中是7天之前的文章不进入文章热榜计算，这里增加一个或的判断条件 arts[len(arts)-1].Update 是否在7天之前即可
		// 是的话，则 break
		if len(arts) < int(b.batchSize) {
			break
		}

		offset = offset + int64(len(arts))
	}

	// 取出优先级队列中的数据
	// 需要注意：topN.Dequeue 得到的值是 score 最小的，因此需要倒序取出
	ids := make([]int64, b.n)
	for i := b.n - 1; i >= 0; i-- {
		v, err := topN.Dequeue()
		// 优先队列中不足 n 个
		if err != nil {
			break
		}
		ids[i] = v.art.Id
	}
	return ids, nil
}
