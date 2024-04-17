package main

import (
	"log"
	"testing"

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
