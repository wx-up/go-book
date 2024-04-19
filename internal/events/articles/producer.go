package article

import (
	"context"

	"github.com/wx-up/go-book/internal/events"

	"github.com/IBM/sarama"
)

const TopicReadEvent = "article_read"

// Producer article 领域的所有事件通过这个接口来发送
type Producer interface {
	ProduceReadEvent(ctx context.Context, evt ReadEvent) error
	BatchProduceReadEvent(ctx context.Context, evt BatchReadEvent) error
}

type ReadEvent struct {
	Uid int64
	Aid int64
}

type BatchReadEvent struct {
	Us []int64
	As []int64
}

type KafkaProducer struct {
	producer sarama.SyncProducer
	topic    string
}

func (k *KafkaProducer) BatchProduceReadEvent(ctx context.Context, evt BatchReadEvent) error {
	// TODO implement me
	panic("implement me")
}

// NewKafkaProducer 保持依赖注入
func NewKafkaProducer(producer sarama.SyncProducer) *KafkaProducer {
	return &KafkaProducer{
		producer: producer,
		topic:    TopicReadEvent,
	}
}

// ProduceReadEvent 如果重试逻辑很复杂就使用装饰器
// 如果很简单的话，则直接在 ProduceReadEvent 中实现
// 消费者那端需要保证幂等性
func (k *KafkaProducer) ProduceReadEvent(ctx context.Context, evt ReadEvent) error {
	_, _, err := k.producer.SendMessage(&sarama.ProducerMessage{
		Topic: k.topic,
		Value: events.NewJsonEncoder(evt),
	})
	return err
}
