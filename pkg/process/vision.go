package process

import (
	"fmt"
	"log"
	"os"
	"strings"

	vision "cloud.google.com/go/vision/apiv1"
	"github.com/mchmarny/webcr/pkg/commons"
	"golang.org/x/net/context"
)

var (
	logger = log.New(os.Stdout, "[processor] ", log.Lshortfile|log.Ldate|log.Ltime)
)

const (
	minNumOfParts = 20
)

// Processor represents the publisher object
type Processor struct {
	client     *vision.ImageAnnotatorClient
	context    context.Context
	subscriber commons.Subscriber
	publisher  commons.Publisher
}

// NewProcessor creates new Processor
func NewProcessor(projectID string, sub commons.Subscriber, pub commons.Publisher) (processor *Processor, err error) {

	p := &Processor{
		context:    context.Background(),
		subscriber: sub,
		publisher:  pub,
	}

	client, err := vision.NewImageAnnotatorClient(p.context)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client: %v", err)
	}

	p.client = client

	return p, nil

}

// Process processes WebResource against Google Vision API
func (p *Processor) Process() {

	itemChan := make(chan *commons.WebResource, 1)
	go p.subscriber.Subscribe(itemChan)

	for {
		select {
		case wr := <-itemChan:
			p.getText(wr)
			break
		}
	}
}

// getText actually processes WebResource against Google Vision API
func (p *Processor) getText(wr *commons.WebResource) (success bool) {

	//logger.Printf("Processing -> %s", wr.String())
	parts, err := p.client.DetectTexts(p.context, vision.NewImageFromURI(wr.Link), nil, 100)
	if err != nil {
		logger.Printf("Error white detecting text from %s: %v", wr.Link, err)
		return false
	}

	if len(parts) >= minNumOfParts {
		list := []string{}
		for _, s := range parts {
			list = append(list, strings.Replace(s.Description, "\n", ",", -1))
		}
		wr.Parts = list
		//logger.Printf("Resource: %s", item.String())
		p.publisher.Publish(wr)
		logger.Printf("Published -> %s", wr.String())
	} else {
		//logger.Printf("Skipping Resource: %s", wr.String())
		return false
	}

	return true

}
