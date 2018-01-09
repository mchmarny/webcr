package commons

import (
	"fmt"
	"strings"
)

const (
	// ExcludeDomainFile holds a name of the exclude domain file
	ExcludeDomainFile = "exclude-domains.txt"
	// ExcludeTitlesFile holds a name of the exclude title file
	ExcludeTitlesFile = "exclude-titles.txt"
)

// Filter represents search filter
type Filter struct {
	parts []string
}

// NewFilter returns new filter instance
func NewFilter(filePath string) (filter *Filter, err error) {

	if !PathExists(filePath) {
		return nil, fmt.Errorf("File does not existst: %s", filePath)
	}

	f := &Filter{}

	lines, err := GetFileLines(filePath)
	if err != nil {
		return f, fmt.Errorf("Error while reading file: %s -> %v", filePath, err)
	}
	f.parts = lines

	return f, nil
}

// ShouldExclude validates passed string against exclude list
func (f *Filter) ShouldExclude(s string) bool {

	// empty
	if s == "" {
		return true
	}

	// in title array
	for _, part := range f.parts {
		if strings.Contains(strings.ToLower(s), part) {
			return true
		}
	}

	return false

}
