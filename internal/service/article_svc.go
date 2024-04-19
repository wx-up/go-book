package service

import (
	"context"
	"time"

	"github.com/wx-up/go-book/pkg/logger"

	events "github.com/wx-up/go-book/internal/events/articles"

	"github.com/wx-up/go-book/internal/repository"

	"github.com/wx-up/go-book/internal/domain"
)

var ErrArticleNotFound = repository.ErrArticleNotFound

//go:generate mockgen -destination=./mocks/article_svc.mock.go -package=svcmocks -source=article_svc.go
type ArticleService interface {
	Save(ctx context.Context, article domain.Article) (int64, error)
	Publish(ctx context.Context, article domain.Article) (int64, error)
	List(ctx context.Context, uid int64, page, size int64) ([]domain.Article, error)
	Detail(ctx context.Context, id int64) (domain.Article, error)

	PublishedDetail(ctx context.Context, id int64) (domain.Article, error)
}

type articleService struct {
	repo     repository.ArticleRepository
	producer events.Producer
	l        logger.Logger
	ch       chan readInfo
}

func NewArticleService(repo repository.ArticleRepository) ArticleService {
	return &articleService{
		repo: repo,
		ch:   make(chan readInfo, 10),
	}
}

type readInfo struct {
	aid int64
	uid int64
}

func (a *articleService) PublishedDetail(ctx context.Context, id int64) (domain.Article, error) {
	// 在这里调用 userService 进行 组装
	art, err := a.repo.GetPublishedById(ctx, id)
	if err == nil {
		er := a.producer.ProduceReadEvent(ctx, events.ReadEvent{
			// 即便消费者需要 art 的其他字段信息，也是需要消费者拿到id自己去查询的
			// 不要在这里查询之后放在消息体里
			// 除非是消费者需要关心当下的信息（ 快照信息 ），则需要写到消息体中
			Aid: id,
			Uid: 0,
		})
		if er != nil {
			a.l.Error("发送读者阅读事件失败", logger.Error(er))
		}
	}
	return art, err
}

func NewArticleServiceV1(repo repository.ArticleRepository) ArticleService {
	ch := make(chan readInfo, 10)
	svc := &articleService{
		repo: repo,
		ch:   ch,
	}

	// 如果程序退出了，这里还要考虑优雅退出的问题
	go func() {
		for {
			us := make([]int64, 0, 10)
			as := make([]int64, 0, 10)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			for i := 0; i < 10; i++ {
				select {
				case <-ctx.Done():
					goto label
				case info, ok := <-ch:
					cancel()
					if !ok {
						if len(us) > 0 {
							goto label
						} else {
							return
						}
					}
					us = append(us, info.uid)
					as = append(as, info.aid)
				}
			}
		label:
			cancel()

			// 批量发送消息
			ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
			err := svc.producer.BatchProduceReadEvent(ctx, events.BatchReadEvent{
				Us: us,
				As: as,
			})
			cancel()
			if err != nil {
				// TODO 记录日志
			}
		}
	}()
	return svc
}

func (a *articleService) PublishedDetailV1(ctx context.Context, id int64) (domain.Article, error) {
	// 在这里调用 userService 进行 组装
	art, err := a.repo.GetPublishedById(ctx, id)
	if err == nil {
		go func() {
			a.ch <- readInfo{
				aid: id,
				uid: 0,
			}
		}()
	}
	return art, err
}

func (a *articleService) Detail(ctx context.Context, id int64) (domain.Article, error) {
	// TODO implement me
	panic("implement me")
}

func (a *articleService) List(ctx context.Context, uid, page, size int64) ([]domain.Article, error) {
	return a.repo.GetByAuthorId(ctx, uid, page, size)
}

func (a *articleService) Publish(ctx context.Context, article domain.Article) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (a *articleService) Save(ctx context.Context, d domain.Article) (int64, error) {
	if d.Id <= 0 {
		// 新增
		return a.repo.Create(ctx, d)
	}
	// 修改
	// 有一种做法：先根据id查询文档，不存在就报错，存在就更新，相当于查询数据库两次，性能会比较差
	// 另一种做法就是直接更新，根据 update 语句返回的受影响行数来判断是否更新成功，性能会比较好
	return d.Id, a.repo.Update(ctx, d)
}
