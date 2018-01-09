package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mchmarny/webcr/pkg/broker"
	"github.com/mchmarny/webcr/pkg/process"
	"github.com/mchmarny/webcr/pkg/search"
	"golang.org/x/net/context"
)

var (
	logger = log.New(os.Stdout, "[app] ", log.Lshortfile|log.Ldate|log.Ltime)
)

const (
	maskedValue      = "******"
	linkTopic        = "links"
	imgTopic         = "images"
	linkSubcritption = "links-sub"
)

func main() {

	// FLAGS
	projectID := flag.String("project", os.Getenv("SEARCH_PROJECT"),
		"Google Cloud Platform projectID")
	ctx := flag.String("cx", os.Getenv("SEARCH_CTX"),
		"Google search context (https://cse.google.com/cse/setup/basic)")
	key := flag.String("key", os.Getenv("SEARCH_KEY"),
		"Google API key (https://pantheon.corp.google.com/apis/credentials?project=xxx)")
	query := flag.String("q", "boarding pass airline ticket",
		"Search query [boarding pass airline ticket]")
	confDir := flag.String("conf", "./config", "Path to configuration directory [./config]")

	flag.Parse()

	//TODO: split it into four specific actions
	//      if search then what's required

	if *ctx == "" || *key == "" || *query == "" || *projectID == "" {
		logger.Panicf("Missing required arguments: ctx=%s, key=%s, query=%s, projectID=%s",
			maskedValue, maskedValue, *query, *projectID)
	}
	// END FLAGAS

	// CMD
	appContext, cancel := context.WithCancel(context.Background())
	go func() {
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		log.Println(<-ch)
		cancel()
		os.Exit(0)
	}()
	// END CMD

	// PUBLISHER
	linkPub, err := broker.NewPublisher(*projectID, linkTopic)
	imgPub, err := broker.NewPublisher(*projectID, imgTopic)
	// END PUBLISHER

	// SUBSCRIBERS
	linkSub, err := broker.NewSubscriber(*projectID, linkSubcritption)
	// END SUBSCRIBERS

	// PROCESSOR
	p, err := process.NewProcessor(*projectID, linkSub, imgPub)
	if err != nil {
		logger.Panicf("Error while creating processor: %v", err)
	}
	go p.Process()
	// END PROCESSOR

	// SEARCH
	s, err := search.NewSearch(*ctx, *key, *query, linkPub, *confDir)
	if err != nil {
		logger.Panicf("Error while creating searcher: %v", err)
	}
	go s.Do()
	// END SEARCH

	// LOOP
	for {
		select {
		case <-appContext.Done():
			break
		}
	}
	// END LOOP

}
