package articles

import (
	"context"
	"github.com/wx-up/go-book/interactive/repository"
	"time"

	"github.com/IBM/sarama"
	"github.com/wx-up/go-book/pkg/logger"
	"github.com/wx-up/go-book/pkg/saramax"
)

type BatchReadEventKafkaConsumer struct {
	l      logger.Logger
	repo   repository.InteractiveRepository
	client sarama.Client
}

func NewBatchReadEventKafkaConsumer(l logger.Logger, repo repository.InteractiveRepository, client sarama.Client) *BatchReadEventKafkaConsumer {
	return &BatchReadEventKafkaConsumer{
		l:      l,
		repo:   repo,
		client: client,
	}
}

func (k *BatchReadEventKafkaConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive_read_event", k.client)
	if err != nil {
		return err
	}
	go func() {
		er := cg.Consume(
			context.Background(),
			[]string{TopicReadEvent},
			saramax.NewBatchHandler[ReadEvent](k.l, k.Consume))
		if er != nil {
			k.l.Error("退出消费", logger.Error(er))
		}
	}()
	return nil
}

// Consume 一般 consumer 都要处理幂等，当然业务场景可以不处理幂等
func (k *BatchReadEventKafkaConsumer) Consume(message []*sarama.ConsumerMessage, t []ReadEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	ids := make([]int64, 0, len(t))
	for _, m := range t {
		ids = append(ids, m.Aid)
	}
	err := k.repo.BatchIncrReadCnt(ctx, []string{"articles"}, ids)
	if err != nil {
		// 记录日志
	}
	return err
}
