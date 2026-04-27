package transport

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
)

type MessageHandler interface {
	HandleMessage(ctx context.Context, value []byte) error
}

type KafkaConsumer struct {
	brokers []string
	topic   string
	handler MessageHandler
	cfg     *sarama.Config
}

func NewKafkaConsumer(brokers []string, topic string, handler MessageHandler, cfg *sarama.Config) *KafkaConsumer {

	return &KafkaConsumer{
		brokers: brokers,
		topic:   topic,
		handler: handler,
		cfg:     cfg,
	}
}

func (c *KafkaConsumer) Run(ctx context.Context) error {

	consumer, err := sarama.NewConsumer(c.brokers, c.cfg)
	if err != nil {
		return fmt.Errorf("failed to create consumer: %w", err)
	}
	defer consumer.Close()

	pc, err := consumer.ConsumePartition(c.topic, 0, sarama.OffsetNewest)
	if err != nil {
		return fmt.Errorf("failed to consume partition: %w", err)
	}
	defer pc.Close()

	for {
		select {
		case msg := <-pc.Messages():
			if msg == nil {
				continue
			}

			if err := c.handler.HandleMessage(ctx, msg.Value); err != nil {
				return err
			}

		case err := <-pc.Errors():
			if err != nil {
				return err
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
