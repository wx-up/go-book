package article

import (
	"context"
	"time"

	"github.com/IBM/sarama"
	"github.com/wx-up/go-book/internal/repository"
	"github.com/wx-up/go-book/pkg/logger"
	"github.com/wx-up/go-book/pkg/saramax"
)

type ReadEventKafkaConsumer struct {
	l      logger.Logger
	repo   repository.InteractiveRepository
	client sarama.Client
}

func NewReadEventKafkaConsumer(l logger.Logger, repo repository.InteractiveRepository, client sarama.Client) *ReadEventKafkaConsumer {
	return &ReadEventKafkaConsumer{
		l:      l,
		repo:   repo,
		client: client,
	}
}

func (k *ReadEventKafkaConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive_read_event", k.client)
	if err != nil {
		return err
	}
	go func() {
		er := cg.Consume(
			context.Background(),
			[]string{TopicReadEvent},
			saramax.NewHandler[ReadEvent](k.l, k.Consume))
		if er != nil {
			k.l.Error("退出消费", logger.Error(er))
		}
	}()
	return nil
}

// Consume 一般 consumer 都要处理幂等，当然业务场景可以不处理幂等
func (k *ReadEventKafkaConsumer) Consume(message *sarama.ConsumerMessage, t ReadEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return k.repo.IncrReadCnt(ctx, "articles", t.Aid)
}
