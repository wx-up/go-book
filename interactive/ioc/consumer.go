package ioc

import (
	"github.com/wx-up/go-book/interactive/events/articles"
	"github.com/wx-up/go-book/pkg/saramax"
)

func CreateConsumers(ReadEventKafkaConsumer *articles.ReadEventKafkaConsumer) []saramax.Consumer {
	return []saramax.Consumer{
		ReadEventKafkaConsumer,
	}
}
