package commons

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"
)

func TestDomainFilter(t *testing.T) {

	wd, _ := os.Getwd()
	wd = strings.Replace(wd, "/pkg/commons", "/config", 1)

	filter, err := NewFilter(path.Join(wd, ExcludeDomainFile))
	if err != nil {
		t.Fatalf("Error while reading config file: %s -> %v", ExcludeDomainFile, err)
	}

	if filter.ShouldExclude("https://sub.domain.com/folder/image.jpg?q=1&v=a") {
		t.Fatal("Filter fails on valid domain (sans port")
	}

	if filter.ShouldExclude("http://domain.co:8080/image.jpg") {
		t.Fatal("Filter fails on valid domain (with port")
	}

	if !filter.ShouldExclude(fmt.Sprintf("https://%s/folder/image.jpg?q=1&v=a", filter.parts[0])) {
		t.Fatal("Filter fails to exclude invalid domain")
	}

}

func TestTitleFilter(t *testing.T) {

	wd, _ := os.Getwd()
	wd = strings.Replace(wd, "/pkg/commons", "/config", 1)

	filter, err := NewFilter(path.Join(wd, ExcludeTitlesFile))
	if err != nil {
		t.Fatalf("Error while reading config file: %s -> %v", ExcludeTitlesFile, err)
	}

	if !filter.ShouldExclude(fmt.Sprintf("My boarding %s pass", filter.parts[0])) {
		t.Fatal("Filter fails to exclude invalid title")
	}

	if filter.ShouldExclude("My boarding pass") {
		t.Fatal("Filter fails to exclude invalid title")
	}

}
