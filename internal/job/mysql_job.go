package job

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"golang.org/x/sync/semaphore"

	"github.com/wx-up/go-book/internal/domain"

	"github.com/wx-up/go-book/pkg/logger"

	"github.com/wx-up/go-book/internal/service"
)

// Executor 任务执行器抽象
// 可以是本地函数、rpc 或者 http 调用等等
type Executor interface {
	Name() string
	Exec(context.Context, domain.Job) error
}

type HTTPExecutor struct{}

func (H *HTTPExecutor) Name() string {
	return "http_executor"
}

func (H *HTTPExecutor) Exec(ctx context.Context, job domain.Job) error {
	type Cfg struct {
		Endpoint string `json:"endpoint"`
		Method   string `json:"method"`
	}
	var cfg Cfg
	err := json.Unmarshal([]byte(job.Cfg), &cfg)
	if err != nil {
		return err
	}

	// job.Cfg 还可以增加请求参数的配置
	req, err := http.NewRequest(cfg.Method, cfg.Endpoint, nil)
	_ = req
	return nil
}

type LocalFunc func(ctx context.Context, job domain.Job) error

type LocalFuncExecutor struct {
	fs map[string]LocalFunc
}

func NewLocalFuncExecutor() *LocalFuncExecutor {
	return &LocalFuncExecutor{
		fs: make(map[string]LocalFunc),
	}
}

func (l *LocalFuncExecutor) Name() string {
	return "local_func_executor"
}

func (l *LocalFuncExecutor) RegisterFunc(name string, f LocalFunc) {
	l.fs[name] = f
}

func (l *LocalFuncExecutor) Exec(ctx context.Context, job domain.Job) error {
	fn, ok := l.fs[job.Name]
	if !ok {
		return errors.New("未找到对应的执行函数")
	}
	return fn(ctx, job)
}

type Scheduler struct {
	execs      map[string]Executor
	jobService service.JobService
	logger     logger.Logger

	// 信号量
	limiter *semaphore.Weighted
}

func NewScheduler(jobService service.JobService, logger logger.Logger) *Scheduler {
	return &Scheduler{
		execs:      make(map[string]Executor),
		jobService: jobService,
		logger:     logger,
		limiter:    semaphore.NewWeighted(100), // 信号量，限制并发，控制 Schedule 最大的 goroutine 数量
	}
}

func (s *Scheduler) RegisterExecutor(exec Executor) {
	s.execs[exec.Name()] = exec
}

func (s *Scheduler) Schedule(b context.Context) error {
	for {

		// 退出调度循环
		if b.Err() != nil {
			return b.Err()
		}

		// 进一步优化，可以考虑使用 goroutine 的池化技术
		// 信号量虽然可以控制最大的 goroutine 数量，但是每次都是新创建 goroutine
		err := s.limiter.Acquire(b, 1)
		if err != nil {
			s.logger.Error("获取信号量失败", logger.Error(err))
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		j, err := s.jobService.Preempt(ctx)
		cancel()
		if err != nil {
			s.logger.Error("抢占任务失败", logger.Error(err))
			continue
		}

		exec, ok := s.execs[j.Executor]
		if !ok {
			// 本地debug 的时候可以考虑直接panic
			// 线上则记录日志之后，继续执行下一个任务
			s.logger.Error("获取执行器失败", logger.Error(err), logger.String("executor", j.Executor))
			continue
		}

		// 执行
		go func() {
			defer func() {
				s.limiter.Release(1)
				// 释放任务
				err1 := j.CancelFunc()
				if err1 != nil {
					s.logger.Error("释放任务失败", logger.Error(err1))
				}
			}()

			// 这里执行错误了，也是需要更新 next_time 的
			err2 := exec.Exec(b, j)
			if err2 != nil {
				// 可以考虑在这里先重试，重试失败再记录日志
				s.logger.Error("执行任务失败", logger.Error(err2))
			}

			// 设置 next_time 的值，用于下一次调度
			err2 = s.jobService.ResetNextTime(context.Background(), j)
			if err2 != nil {
				s.logger.Error("设置下一次执行时间失败", logger.Error(err2))
			}
		}()

	}
}
