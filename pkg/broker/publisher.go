package broker

import (
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/mchmarny/webcr/pkg/commons"
	"golang.org/x/net/context"
)

var (
	logger = log.New(os.Stdout, "[broker] ", log.Lshortfile|log.Ldate|log.Ltime)
)

// Publisher represents the publisher object
type Publisher struct {
	topic   *pubsub.Topic
	client  *pubsub.Client
	context context.Context
}

// NewPublisher creates new Publisher
func NewPublisher(projectID, topic string) (publisher *Publisher, err error) {

	p := &Publisher{
		context: context.Background(),
	}

	client, err := pubsub.NewClient(p.context, projectID)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client: %v", err)
	}

	p.client = client
	p.topic = p.client.Topic(topic)

	return p, nil

}

// Publish publishes WebResource to configured PubSub
func (p *Publisher) Publish(wr *commons.WebResource) (success bool) {

	content, err := wr.ToJSON()
	if err != nil {
		logger.Printf("Error while serializing object: %v:%v", wr, err)
		return false
	}

	msg := &pubsub.Message{Data: content}
	result := p.topic.Publish(p.context, msg)
	id, err := result.Get(p.context)
	if err != nil {
		logger.Fatalf("Error while publishing message: %v:%v", err, id)
		return false
	}

	return true

}
