package saramax

import (
	"context"
	"encoding/json"
	"time"

	"github.com/IBM/sarama"
	"github.com/wx-up/go-book/pkg/logger"
)

type BatchHandler[T any] struct {
	l            logger.Logger
	fn           func(msg []*sarama.ConsumerMessage, ts []T) error
	batchSize    int
	batchTimeout time.Duration
}

type BatchHandlerOption[T any] func(*BatchHandler[T])

func WithBatchSize[T any](size int) BatchHandlerOption[T] {
	return func(h *BatchHandler[T]) {
		h.batchSize = size
	}
}

func WithBatchTimeout[T any](timeout time.Duration) BatchHandlerOption[T] {
	return func(b *BatchHandler[T]) {
		b.batchTimeout = timeout
	}
}

func NewBatchHandler[T any](l logger.Logger, fn func(msg []*sarama.ConsumerMessage, ts []T) error, opts ...BatchHandlerOption[T]) *BatchHandler[T] {
	b := &BatchHandler[T]{
		l:            l,
		fn:           fn,
		batchSize:    10,
		batchTimeout: time.Second * 5,
	}
	for _, opt := range opts {
		opt(b)
	}
	return b
}

func (b *BatchHandler[T]) Setup(session sarama.ConsumerGroupSession) error {
	// TODO implement me
	panic("implement me")
}

func (b *BatchHandler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	// TODO implement me
	panic("implement me")
}

func (b *BatchHandler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msg := claim.Messages()
	const batchSize = 10
	for {
		batch := make([]*sarama.ConsumerMessage, 0, batchSize)
		ts := make([]T, 0, batchSize)

		// 超时控制，5秒钟需要有一个批次的消息
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		done := false
		for i := 0; i < batchSize; i++ {
			select {
			case <-ctx.Done():
				done = true
			case m, ok := <-msg:
				if !ok {
					cancel()
					return nil
				}
				var t T
				err := json.Unmarshal(m.Value, &t)
				if err != nil {
					b.l.Error("反序列消息体失败",
						logger.String("topic", m.Topic),
						logger.Int32("partition", m.Partition),
						logger.Int64("offset", m.Offset),
						logger.Error(err))
					continue
				}
				batch = append(batch, m)
				ts = append(ts, t)
			}

			// 如果是超时了，就直接退出批次循环了
			if done {
				break
			}
		}
		cancel()

		// 凑够一批，然后交给业务处理
		err := b.fn(batch, ts)
		if err != nil {
			// 这里可以尝试把整个批次的消息都记录下来，不记录也没事
			// 因为业务在处理消息的时候，自己应该要重试，重试失败之后自己应该要会记录
			b.l.Error("处理消息失败",
				logger.Error(err))
		}

		// 提交消息，继续消费下一批
		for _, msg := range batch {
			session.MarkMessage(msg, "")
		}
	}
}
