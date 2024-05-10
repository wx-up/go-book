package job

import (
	"context"
	"sync"
	"time"

	"github.com/wx-up/go-book/pkg/logger"

	rlock "github.com/gotomicro/redis-lock"
	"github.com/wx-up/go-book/internal/service"
)

type ArticleRankingJob struct {
	client  *rlock.Client
	svc     service.ArticleRankingService
	timeout time.Duration // job 执行的超时时间
	key     string
	logger  logger.Logger

	lock *rlock.Lock

	localLock sync.Mutex
}

func NewArticleRankingJob(
	client *rlock.Client,
	svc service.ArticleRankingService,
	timeout time.Duration,
	logger logger.Logger,
) *ArticleRankingJob {
	return &ArticleRankingJob{
		client:    client,
		svc:       svc,
		timeout:   timeout,
		logger:    logger,
		key:       "lock:job:article_ranking",
		localLock: sync.Mutex{},
	}
}

func (a *ArticleRankingJob) Name() string {
	return "article_ranking_job"
}

func (a *ArticleRankingJob) Close() error {
	a.localLock.Lock()
	lock := a.lock
	a.lock = nil
	a.localLock.Unlock()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return lock.Unlock(ctx)
}

func (a *ArticleRankingJob) CloseV1() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	a.localLock.Lock()
	defer func() {
		a.lock = nil
		a.localLock.Unlock()
	}()
	return a.lock.Unlock(ctx)
}

func (a *ArticleRankingJob) Run() error {
	a.localLock.Lock()
	if a.lock == nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		lock, err := a.client.Lock(ctx, a.key, a.timeout, &rlock.FixIntervalRetry{
			Interval: time.Millisecond * 100,
			Max:      3,
		}, time.Second)
		if err != nil {
			// 获取锁失败，可能是锁已经被其他进程获取
			return nil
		}
		a.lock = lock
		a.localLock.Unlock()
		go func() {
			// 保证实例始终持有分布式锁
			// 续约机制
			// 第一个参数是续约的间隔时间，第二个参数是调用redis的超时时间
			// 时间间隔要小于锁过期的时间，也要考虑进去续约本身的执行时间
			er := a.lock.AutoRefresh(a.timeout/2, time.Second)
			if er != nil {
				// 续约失败，可能是锁已经过期或者调用redis失败
				a.localLock.Lock()
				// 将 a.lock=nil 用于下次执行时再尝试抢锁
				a.lock = nil
				a.localLock.Unlock()
			}
		}()

	}

	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()
	return a.svc.TopN(ctx)
}
