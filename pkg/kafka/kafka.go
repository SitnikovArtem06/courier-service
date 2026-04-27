package kafka

import (
	"fmt"
	"os"
	"strings"

	"github.com/IBM/sarama"
)

type KafkaEnvConfig struct {
	Brokers []string
	Topic   string
}

func InitKafka() (KafkaEnvConfig, *sarama.Config, error) {
	brokersRaw := os.Getenv("KAFKA_BROKERS")
	topic := os.Getenv("KAFKA_ORDER_TOPIC")

	if brokersRaw == "" || topic == "" {
		return KafkaEnvConfig{}, nil, fmt.Errorf("KAFKA_BROKERS or KAFKA_ORDER_TOPIC is empty")
	}

	parts := strings.Split(brokersRaw, ",")
	brokers := make([]string, 0, len(parts))
	for _, p := range parts {
		if p != "" {
			brokers = append(brokers, p)
		}
	}
	if len(brokers) == 0 {
		return KafkaEnvConfig{}, nil, fmt.Errorf("KAFKA_BROKERS is empty after parsing")
	}

	cfg := sarama.NewConfig()
	cfg.Consumer.Return.Errors = true

	return KafkaEnvConfig{
		Brokers: brokers,
		Topic:   topic,
	}, cfg, nil
}
