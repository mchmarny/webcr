# WebCR

OCR fro the Web - Find and extract text from social media images (e.g. airline boarding passes)

> Note, DO NOT post images of your boarding pass or airline tickets, even with just a barcode as that can be easilly converted to data as well and it includes a more data than you think

## What does it do

#### Aquire Images

1. Search Google for images for a specific search query
2. Filter out images based on provided criteria (domains and titles)
2. Publish valida to images to PubSub

#### Process Images

1. Subscribe to the PubSub topic, for each event:
2. Process image throiugh Google Vision API for OCR to get text
3. Publish extracted data back to PubSub 
 
#### Analyze Data

1. Subscribe to the PubSub topic using Cloud Dataflow and drain (using template) all events into BogQuery table
2. Use SQL to find interesting insight ;)

## Setup 

#### Setup Google API

* Go to [Google API Console](https://console.developers.google.com/) to create a new project
* Go to `Library` and search `Custom Search API` and enable it
* Go to `Credentials` and click `Create Credentials` and choose `API Key` (do restrict the key if you can). The value under `Key` is your API key
* Define that value as `$SEARCH_KEY` in your envirnemnt variables or pass it into `webcr`  as `--key` argument on each execution

#### Setup Google Custom Search

* Go to [Google Cusomt Search](https://cse.google.com/cse/all) and create a new `engine` (enter *.google.com and from the `hunting grounds` list below as `Sites to Seatch`)
* Click on the newly created search engine and change the `Sites to Search` option to `Search the entire web but emphasize included sites`
* Copy the Engine Context ID (value of the `cx=xxx` argument in URL)
* Make sure the `Image Search` option is enabled and Update button at the bottom when you are done
* Define that value as `$SEARCH_CTX` in your envirnemnt variables or pass into the `webcr` as `--ctx` argument on each execution

Also good hunting grounds:

* *.flickr.com
* *.snapchat.com
* *.facebook.com
* *.instagram.com
* *.twitter.com
* *.medium.com

> Note: `Free` version of Google Custom Search has a limit of `100` searches per day



#### Google Vision API

> Default free quota includes 1,000 queries and 600 queries per min limit

#### PubSub Topic and Subscriptions 

You can create the necessary PubSub topics and subscription by hand using `gcloud` commands (see below) or just use the provided `make` file (`make gcp-setup`) which will create all the GCP dependancies.

```
	gcloud pubsub topics create $(TOPIC_LINKS)
	gcloud pubsub subscriptions create $(TOPIC_LINKS_SUB) --topic=$(TOPIC_LINKS)
	gcloud pubsub topics create $(TOPIC_IMAGES)
	gcloud pubsub subscriptions create $(TOPIC_IMAGES_SAVE_SUB) --topic=$(TOPIC_LINKS)
```

#### BigQuery Table

Create a BigQuery table. You can do this using the `bq` command or just use the provided `make` file (`make gcp-setup`) which will create all the GCP dependancies.

```
	bq mk $(APP_NAME)
	bq mk -t $(APP_NAME).$(TOPIC_IMAGES) ./config/images-schema.json 
```

#### Dataflow Job

```
	gsutil mb gs://$(PROJECT_NAME)-$(APP_NAME)-tmp
	gcloud dataflow jobs run $(APP_NAME)-JOB \
		--gcs-location gs://dataflow-templates/pubsub-to-bigquery/template_file \
		--parameters="topic=projects/$(PROJECT_NAME)/topics/$(TOPIC_IMAGES)","table=$(PROJECT_NAME):$(APP_NAME).$(TOPIC_IMAGES)","stagingLocation=gs://$(PROJECT_NAME)-$(APP_NAME)-tmp"
```