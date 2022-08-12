package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	cepubsub "github.com/cloudevents/sdk-go/protocol/pubsub/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type Example struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
}

func receive(ctx context.Context, event cloudevents.Event) {
	switch event.Type() {
	case "com.synack.sample.sent":
		data, _ := json.MarshalIndent(&event, "", "\t")
		fmt.Printf("Got event\n%v\n", string(data))

		// unmarshalling example
		example := &Example{}
		if err := event.DataAs(example); err != nil {
			panic(err)
		}
		fmt.Printf("Got example %#v\n", *example)
	default:
		fmt.Printf("Unknown event type: %v\n", event.Type())
		data, _ := json.MarshalIndent(&event, "", "\t")
		fmt.Println(string(data))
	}

}

func main() {
	// Connect to localhost if not running inside docker
	if os.Getenv("PUBSUB_EMULATOR_HOST") == "" {
		os.Setenv("PUBSUB_EMULATOR_HOST", "localhost:8085")
	}

	t, err := cepubsub.New(context.Background(),
		cepubsub.AllowCreateSubscription(true),
		cepubsub.WithProjectID("test"),
		cepubsub.WithSubscriptionAndTopicID("test", "test"),
	)
	if err != nil {
		panic(err)
	}

	c, err := cloudevents.NewClient(t)
	if err != nil {
		panic(err)
	}

	if err := c.StartReceiver(context.Background(), receive); err != nil {
		log.Fatalf("failed to start pubsub receiver, %s", err.Error())
	}
}
