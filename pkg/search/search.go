package search

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/mchmarny/webcr/pkg/commons"
)

var (
	logger = log.New(os.Stdout, "[search] ", log.Lshortfile|log.Ldate|log.Ltime)
)

const (
	// Speck: https://developers.google.com/custom-search/json-api/v1/reference/cse/list
	baseURL          = "https://www.googleapis.com/customsearch/v1?cx=%s&key=%s&fields=items(link,title)&searchType=image&q=%s&num=%d"
	pagerPart        = "&start=%d"
	pageSize         = 10
	maxQueriesPerMin = 500
)

// Result holds google search result object
type Result struct {
	Itemes []*commons.WebResource `json:"items"`
}

// NewSearch returns configured search object
func NewSearch(ctx, key, query string, publisher commons.Publisher, configDirPath string) (s *Search, e error) {

	if ctx == "" || key == "" || query == "" {
		return nil, fmt.Errorf("Missing required arguments: ctx=%d, key=%d, query=%s",
			len(ctx), len(key), query)
	}

	if !commons.PathExists(configDirPath) {
		return nil, fmt.Errorf("Config directory does not existst: %s", configDirPath)
	}

	domainFilter, err := commons.NewFilter(path.Join(configDirPath, commons.ExcludeDomainFile))
	if err != nil {
		return nil, fmt.Errorf("Error while reading config file: %s -> %v", commons.ExcludeDomainFile, err)
	}

	titleFilter, err := commons.NewFilter(path.Join(configDirPath, commons.ExcludeTitlesFile))
	if err != nil {
		return nil, fmt.Errorf("Error while reading config file: %s -> %v", commons.ExcludeTitlesFile, err)
	}

	search := &Search{
		query:        query,
		queryURL:     fmt.Sprintf(baseURL, ctx, key, query, pageSize),
		page:         1,
		publisher:    publisher,
		domainFilter: domainFilter,
		titleFilter:  titleFilter,
	}

	return search, nil
}

// Search represents search object
type Search struct {
	query        string
	queryURL     string
	page         int
	publisher    commons.Publisher
	domainFilter *commons.Filter
	titleFilter  *commons.Filter
}

// Do executes search
func (s *Search) Do() {

	logger.Printf("Starting search for: %s", s.query)

	counter := 0
	q := s.queryURL

	for {

		// QUERY
		request, err := http.Get(q)
		if err != nil {
			logger.Printf("Error while getting URL: %s -> %v", q, err)
			return
		}

		if request.StatusCode != http.StatusOK {
			logger.Printf("Invalid http status code: %d", request.StatusCode)
			return
		}
		// END QUERY

		// CONTENT
		defer request.Body.Close()
		content, _ := ioutil.ReadAll(request.Body)

		buf := bytes.NewBuffer(content)
		reader, err := gzip.NewReader(buf)
		if err != nil {
			logger.Printf("Error while decompressing result: %v", err)
			return
		}
		// END CONTENT

		// JSON
		result := Result{}
		dec := json.NewDecoder(reader)
		err = dec.Decode(&result)
		if err != nil && err != io.EOF {
			logger.Printf("Error while parsing JSON from result: %v", err)
			return
		}
		// END JSON

		// OUTPUT URLS
		for _, item := range result.Itemes {
			if !s.titleFilter.ShouldExclude(item.Title) && !s.domainFilter.ShouldExclude(item.Link) {
				item.ID = commons.GetMD5(item.Link)
				s.publisher.Publish(item)
				//logger.Printf("Published: Page:%d Item:%d -> %s", s.page, counter, item.Link)
			}
			counter++
		}
		// END OUTPUT URLS

		// CHECK FOR NO MORE RESULTS
		if len(result.Itemes) < pageSize {
			logger.Printf("Done on page: %d [%d items]", s.page, counter)
			return
		}
		// END CHECK FOR NO MORE RESULTS

		// PAGING
		s.page = s.page + pageSize
		q = s.queryURL + fmt.Sprintf(pagerPart, s.page)
		//logger.Printf("Paged[%d] query URL: %s", s.page, q)
		// END PAGING

	}

}
