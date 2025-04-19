package main

import (
	"context"
	"log"

	"github.com/twmb/franz-go/pkg/kgo"
)

func main() {
	opts := []kgo.Opt{
		kgo.SeedBrokers("192.168.1.203:9092"), // Kafka broker address
		kgo.DefaultProduceTopic("MyUsers"),    // Default topic to produce to
		kgo.ClientID("my-client-id"),          // Optional: client identifier
	}

	client, err := kgo.NewClient(opts...)
	if err != nil {
		log.Fatalf("Failed to create Kafka client: %v", err)
	}
	defer client.Close()

	record := &kgo.Record{
		Value: []byte("Hello World"),
		Topic: "MyUsers",
	}

	if err := client.ProduceSync(context.Background(), record).FirstErr(); err != nil {
		log.Fatalf("Failed to produce message: %v", err)
	}

	log.Println("Message produced successfully")
}
