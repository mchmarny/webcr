package commons

import (
	"encoding/json"
	"fmt"
	"strings"
)

// WebResource represents Web content
type WebResource struct {
	ID    string   `json:"id"`
	Title string   `json:"title"`
	Link  string   `json:"link"`
	Parts []string `json:"parts"`
}

// String returns string representation of the WebResource object
func (r *WebResource) String() string {
	return fmt.Sprintf("WebResource(ID:%s, Title:%s, Link:%s, Parts:%s)",
		r.ID, r.Title, r.Link, strings.Join(r.Parts, ","))
}

// ToWebResource converts bytes into instance reference of WebResource
func ToWebResource(b []byte) (r *WebResource, err error) {
	o := &WebResource{}
	e := json.Unmarshal(b, &o)
	return o, e
}

// ToJSON Marshal current instance into JSON bytes
func (r *WebResource) ToJSON() (content []byte, err error) {
	return json.Marshal(r)
}
