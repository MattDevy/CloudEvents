package main

import (
	"context"
	"log"
	"os"
	"time"

	cepubsub "github.com/cloudevents/sdk-go/protocol/pubsub/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type Example struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
}

func main() {
	// Connect to localhost if not running inside docker
	if os.Getenv("PUBSUB_EMULATOR_HOST") == "" {
		os.Setenv("PUBSUB_EMULATOR_HOST", "localhost:8085")
	}

	t, err := cepubsub.New(context.Background(),
		cepubsub.AllowCreateTopic(true),
		cepubsub.WithProjectID("test"),
		cepubsub.WithTopicID("test"),
	)
	if err != nil {
		panic(err)
	}

	c, err := cloudevents.NewClient(t, cloudevents.WithTimeNow(), cloudevents.WithUUIDs())
	if err != nil {
		panic(err)
	}

	for i := 1; i <= 5; i++ {
		event := cloudevents.NewEvent()
		event.SetType("com.synack.sample.sent")
		event.SetSource("github.com/MattDevy/CloudEvents/cmd/sender/")
		_ = event.SetData("application/json", &Example{
			ID:      i,
			Message: "hello",
		})
		if result := c.Send(context.Background(), event); cloudevents.IsUndelivered(result) {
			log.Printf("failed to send: %v", err)
			os.Exit(1)
		} else {
			log.Printf("sent, accepted: %t", cloudevents.IsACK(result))
		}
		<-time.After(2 * time.Second)
	}

	event := cloudevents.NewEvent()
	event.SetType("com.synack.random.other")
	event.SetSource("github.com/MattDevy/CloudEvents/cmd/sender/")
	_ = event.SetData("application/json", &map[string]string{
		"random": "something",
	})
	if result := c.Send(context.Background(), event); cloudevents.IsUndelivered(result) {
		log.Printf("failed to send: %v", err)
		os.Exit(1)
	} else {
		log.Printf("sent, accepted: %t", cloudevents.IsACK(result))
	}

	os.Exit(0)
}
