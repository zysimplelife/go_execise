package main

import (
	"log"

	"github.com/Shopify/sarama"
)

func main() {
	// For the data collector, we are looking for strong consistency semantics.
	// Because we don't change the flush settings, sarama will try to produce messages
	// as fast as possible to keep latency low.
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	config.Producer.Return.Successes = true
	tlsConfig := createTlsConfiguration()
	if tlsConfig != nil {
		config.Net.TLS.Config = tlsConfig
		config.Net.TLS.Enable = true
	}

	// On the broker side, you may want to change the following settings to get
	// stronger consistency guarantees:
	// - For your broker, set `unclean.leader.election.enable` to false
	// - For the topic, you could increase `min.insync.replicas`.

	producer, err := sarama.NewSyncProducer("127.0.0.1:9092", config)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
	}

	partition, offset, err := producer.SendMessage(&sarama.ProducerMessage{
		Topic: "important",
		Value: sarama.StringEncoder("I have send a message"),
	})

	if err != nil {
		fmt.printf("Failed to store your data:, %s", err)
	} else {
		// The tuple (topic, partition, offset) can be used as a unique identifier
		// for a message in a Kafka cluster.
		fmt.printf("Your data is stored with unique identifier important/%d/%d", partition, offset)
	}
}
