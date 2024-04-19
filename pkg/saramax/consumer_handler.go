package saramax

import (
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/wx-up/go-book/pkg/logger"
)

// 如果不需要记录日志的话可以将 Handler 定义为：
// type Handler[T any] func(message *sarama.ConsumerMessage, t T) error

type Handler[T any] struct {
	logger      logger.Logger
	handlerFunc func(message *sarama.ConsumerMessage, t T) error
}

func NewHandler[T any](logger logger.Logger, handlerFunc func(message *sarama.ConsumerMessage, t T) error) *Handler[T] {
	return &Handler[T]{
		logger:      logger,
		handlerFunc: handlerFunc,
	}
}

func (h *Handler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	messages := claim.Messages()
	for message := range messages {
		var t T
		err := json.Unmarshal(message.Value, &t)
		if err != nil {
			h.logger.Error("反序列化消息失败",
				logger.Error(err),
				logger.String("topic", message.Topic),
				logger.Int32("partition", message.Partition),
				logger.Int64("offset", message.Offset),
			)
			continue
		}

		// 处理消息
		// 为什么把消息传递给 handlerFunc 是因为 handlerFunc 可能需要用到消息的一些属性，比如 headers、topic、partition、offset 等
		// 可以在这里增加重试
		// for i := 0; i < 3; i++ {
		// 	err = h.handlerFunc(message, t)
		// 	if err == nil {
		// 		break
		// 	}
		// 	h.logger.Error
		// }
		err = h.handlerFunc(message, t)
		if err != nil {
			h.logger.Error("处理消息失败",
				logger.Error(err),
				logger.String("topic", message.Topic),
				logger.Int32("partition", message.Partition),
				logger.Int64("offset", message.Offset),
			)
		} else {
			session.MarkMessage(message, "")
		}

	}
	return nil
}
