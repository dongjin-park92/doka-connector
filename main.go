package main

import (
	"context"
	"fmt"
	"kafka-connector/consumers"
	"kafka-connector/loggers"
	"kafka-connector/streamers"
)

func main() {
	fmt.Println(banner)

	streamerIsntance := streamers.NewStreamer()
	consumeInstance := consumers.NewConsumer(streamerIsntance)

	consumerGroup, topic := consumeInstance.CreateConsumer()
	ctx := context.Background()

	loggers.GlobalLogger.Println("start msk consume and kinesis stream putrecord data.")
	err := consumerGroup.Consume(ctx, []string{topic}, consumeInstance)
	if err != nil {
		loggers.GlobalLogger.Fatal(err)
	}

	defer func() {
		if err := consumerGroup.Close(); err != nil {
			loggers.GlobalLogger.Fatal(err)
		}
	}()
}

const (
	banner = `
	██████╗  ██████╗ ██╗  ██╗ █████╗ 
	██╔══██╗██╔═══██╗██║ ██╔╝██╔══██╗
	██║  ██║██║   ██║█████╔╝ ███████║
	██║  ██║██║   ██║██╔═██╗ ██╔══██║
	██████╔╝╚██████╔╝██║  ██╗██║  ██║
	╚═════╝  ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝

		`
)
