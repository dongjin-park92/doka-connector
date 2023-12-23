package consumers

import (
	"kafka-connector/configs"
	"kafka-connector/loggers"
	"kafka-connector/streamers"

	"github.com/Shopify/sarama"
)

type Consumer struct {
	streamer *streamers.Streamer
	conf     *configs.ViperConfig
}

func NewConsumer(streamer *streamers.Streamer) *Consumer {
	conf := configs.GetConfig()
	consumer := &Consumer{
		streamer: streamer,
		conf:     conf,
	}
	go consumer.streamer.RenewClient()
	return consumer
}

func (c *Consumer) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for messege := range claim.Messages() {
		messageValues := string(messege.Value)
		// FIXME: Separate another goroutine for kinesis putcord
		c.recordMessage(messageValues)
		session.MarkMessage(messege, "")
	}
	return nil
}

func (c Consumer) CreateConsumer() (sarama.ConsumerGroup, string) {
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	client, err := sarama.NewClient(c.conf.GetStringSlice("msk.brokers"), config)
	if err != nil {
		loggers.GlobalLogger.Fatal("fail to create msk topic consume:", err)
	}

	consumerGroup, err := sarama.NewConsumerGroupFromClient(c.conf.GetString("msk.groupId"), client)
	if err != nil {
		loggers.GlobalLogger.Fatal("can not create msk consumergroup:", err)
	}

	return consumerGroup, c.conf.GetString("msk.topic")
}

func (c Consumer) recordMessage(message string) {
	c.streamer.PutStreamData(message)
}
