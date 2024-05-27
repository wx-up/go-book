package ioc

import (
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
	article "github.com/wx-up/go-book/interactive/events/articles"
	"github.com/wx-up/go-book/internal/events"
)

func InitKafka() sarama.Client {
	type Config struct {
		Addr []string `json:"addr"`
	}
	saramaCfg := sarama.NewConfig()
	saramaCfg.Producer.Return.Successes = true
	var cfg Config
	err := viper.UnmarshalKey("kafka", &cfg)
	if err != nil {
		panic(err)
	}
	client, err := sarama.NewClient(cfg.Addr, saramaCfg)
	if err != nil {
		panic(err)
	}
	return client
}

func NewSyncProducer(client sarama.Client) sarama.SyncProducer {
	res, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}
	return res
}

func CreateConsumers(ac *article.ReadEventKafkaConsumer) []events.Consumer {
	return []events.Consumer{ac}
}
