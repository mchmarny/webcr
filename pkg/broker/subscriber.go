package broker

import (
	"encoding/json"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/mchmarny/webcr/pkg/commons"
	"golang.org/x/net/context"
)

// Subscriber represents the GCP subscriber object
type Subscriber struct {
	sub     *pubsub.Subscription
	client  *pubsub.Client
	context context.Context
}

// NewSubscriber creates new Subscriber
func NewSubscriber(projectID, subscription string) (subscriber *Subscriber, err error) {

	s := &Subscriber{
		context: context.Background(),
	}

	client, err := pubsub.NewClient(s.context, projectID)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client: %v", err)
	}

	s.client = client
	s.sub = s.client.Subscription(subscription)

	return s, nil

}

// Subscribe pushes PubSub events into a local channel
func (s *Subscriber) Subscribe(out chan *commons.WebResource) {

	err := s.sub.Receive(s.context, func(ctx context.Context, msg *pubsub.Message) {
		item := &commons.WebResource{}
		if err := json.Unmarshal(msg.Data, &item); err != nil {
			logger.Printf("Error while decoding PubSub message: %#v", msg)
			msg.Nack()
		} else {
			//logger.Printf("Event -> %s", item.String())
			out <- item
			msg.Ack()
		}
	})
	if err != nil {
		log.Fatal(err)
	}

}
