package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/require"
)

func TestSarama_Producer(t *testing.T) {
	addr := []string{"localhost:9094"}
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.Partitioner = sarama.NewHashPartitioner
	producer, err := sarama.NewSyncProducer(addr, cfg)
	require.NoError(t, err)
	// 第一个返回值是消息发送到的分区
	// 第二个返回值是消息在分区中的偏移量
	_, _, err = producer.SendMessage(&sarama.ProducerMessage{
		Topic: "test_topic",
		// 可以实现一些 Encoder 接口，比如 JsonEncoder
		// 有些为了消息更加紧凑占用空间更小，可能会使用 protobuf
		Value: sarama.StringEncoder("这是一条消息"),
		// 会在生产者和消费者之间传递
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("trace_id"),
				Value: []byte("123456"),
			},
		},

		// 只用于发送过程，一般对 sarama 进行二次封装的时候会使用
		Metadata: "metadata",
	})
	require.NoError(t, err)
}

func TestSarama_Async_Producer(t *testing.T) {
	config := sarama.NewConfig()
	// 异步生产者不建议把 Errors 和 Successes 都开启，一般开启 Errors 就行
	// 同步生产者就必须都开启，因为会同步返回发送成功或者失败
	config.Producer.Return.Errors = true    // 设定是否需要返回错误信息
	config.Producer.Return.Successes = true // 设定是否需要返回成功信息

	producer, err := sarama.NewAsyncProducer([]string{"localhost:9094"}, config)
	if err != nil {
		log.Fatal("NewAsyncProducer err:", err)
	}
	defer producer.AsyncClose()

	// Input() 获取 channel 用于发送消息
	producer.Input() <- &sarama.ProducerMessage{
		Topic: "test_topic",
		// 可以实现一些 Encoder 接口，比如 JsonEncoder
		// 有些为了消息更加紧凑占用空间更小，可能会使用 protobuf
		Value: sarama.StringEncoder("这是一条消息"),
		// 会在生产者和消费者之间传递
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("trace_id"),
				Value: []byte("123456"),
			},
		},

		// 只用于发送过程，一般对 sarama 进行二次封装的时候会使用
		Metadata: "metadata",
	}

	// 消息是否发送成功，通过 Successes() 和 Errors() channel 获取
	select {
	case suc := <-producer.Successes():
		if suc != nil {
		}
	case fail := <-producer.Errors():
		if fail != nil {
			log.Printf("[Producer] Errors: err:%v msg:%+v \n", fail.Err, fail.Msg)
		}
	}
}

func Test_Sarama_Consumer(t *testing.T) {
	cfg := sarama.NewConfig()
	addr := []string{"localhost:9094"}
	consumer, err := sarama.NewConsumerGroup(addr, "test_group", cfg)
	require.NoError(t, err)

	// 超时控制
	// 一般不会做超时控制，因为程序退出了，自然消费者组也退出了
	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(time.Second*10, func() {
		cancel()
	})
	consumer.Consume(ctx, []string{"test_topic"}, testConsumerGroupHandler{})
}

func Test_Sarama_Async_Consumer(t *testing.T) {
	cfg := sarama.NewConfig()
	addr := []string{"localhost:9094"}
	consumer, err := sarama.NewConsumerGroup(addr, "test_group", cfg)
	require.NoError(t, err)

	// 超时控制
	// 一般不会做超时控制，因为程序退出了，自然消费者组也退出了
	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(time.Second*10, func() {
		cancel()
	})
	consumer.Consume(ctx, []string{"test_topic"}, testConsumerGroupHandler{})
}

type testAsyncConsumerGroupHandler struct{}

func (t testAsyncConsumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	// TODO implement me
	panic("implement me")
}

func (t testAsyncConsumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	// TODO implement me
	panic("implement me")
}

func (t testAsyncConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	ms := claim.Messages()
	const batchSize = 10
	for {
		var eg errgroup.Group
		var last *sarama.ConsumerMessage
		for i := 0; i < batchSize && len(ms) > 0; i++ {
			// 如果1秒中都收取不到消息，则超时退出
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			done := false
			select {
			case <-ctx.Done():
				done = true
			case msg := <-ms:
				cancel()
				last = msg
				eg.Go(func() error {
					// 消费消息
					fmt.Println(msg.Value)
					return nil
				})
			}
			if done {
				break
			}
		}
		err := eg.Wait()
		if err != nil {
			// 说明存在消息消费失败，需要重试，这里重试的话是整个批次重试，也可以分散到每个消息中去重试
			// 如果重试失败就记录日志，这里必须记录日志，否则位移提交之后，就消费不到这个消息了
			continue
		}
		// 提交最后一个消息即可
		session.MarkMessage(last, "")
	}
}

type testConsumerGroupHandler struct{}

func (t testConsumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	// 一个 topic 有多个分区，因此返回的是一个切片
	partitions := session.Claims()["test_topic"]
	for _, part := range partitions {
		session.ResetOffset("test_topic", part, sarama.OffsetOldest, "")
		session.ResetOffset("test_topic", part, 11, "")
	}
	return nil
}

func (t testConsumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	// TODO implement me
	panic("implement me")
}

// ConsumeClaim 消费逻辑
// session 是你和kafka 的会话
func (t testConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	ms := claim.Messages()
	for msg := range ms {
		// 字节
		bs := msg.Value
		// 反序列化
		err := json.Unmarshal(bs, nil)
		if err != nil {
			// 重试
			// 重试还是失败就记录日志（ 记录日志是为了人工处理 ）
			continue
		}
		// 标记消费成功
		session.MarkMessage(msg, "")
	}
	return nil
}
