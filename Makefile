# Parameters
APP_NAME=webcr
PROJECT_NAME=$(SEARCH_PROJECT)
TOPIC_LINKS=links
TOPIC_LINKS_SUB=$(TOPIC_LINKS)-sub
TOPIC_IMAGES=images
TOPIC_IMAGES_SAVE_SUB=$(TOPIC_IMAGES)-save

all: test

gcp-setup:
	gcloud pubsub topics create $(TOPIC_LINKS)
	gcloud pubsub subscriptions create $(TOPIC_LINKS_SUB) --topic=$(TOPIC_LINKS)
	gcloud pubsub topics create $(TOPIC_IMAGES)
	gcloud pubsub subscriptions create $(TOPIC_IMAGES_SAVE_SUB) --topic=$(TOPIC_LINKS)
	bq mk $(APP_NAME)
	bq mk -t $(APP_NAME).$(TOPIC_IMAGES) ./config/images-schema.json 
	gsutil mb gs://$(PROJECT_NAME)-$(APP_NAME)-tmp
	gcloud dataflow jobs run $(APP_NAME)-JOB \
		--gcs-location gs://dataflow-templates/pubsub-to-bigquery/template_file \
		--parameters="topic=projects/$(PROJECT_NAME)/topics/$(TOPIC_IMAGES)","table=$(PROJECT_NAME):$(APP_NAME).$(TOPIC_IMAGES)","stagingLocation=gs://$(PROJECT_NAME)-$(APP_NAME)-tmp"

gcp-cleanup:
	gcloud pubsub topics delete $(TOPIC_LINKS)
	gcloud pubsub subscriptions delete $(TOPIC_LINKS_SUB)
	gcloud pubsub topics delete $(TOPIC_IMAGES)
	gcloud pubsub subscriptions delete $(TOPIC_IMAGES_SAVE_SUB)
	bq rm -f -t $(APP_NAME).$(TOPIC_IMAGES)
	bq rm -f $(APP_NAME)
	# gcloud dataflow jobs cancel [JOB_ID] <- use gcloud dataflow jobs list to find

build:
	go build -o ./bin/webcr -v

test:
	go test -v ./...

run:
	go run *.go

clean:
	go clean
	rm -f ./bin/*

deps:
	go get github.com/tools/godep
	godep restore


