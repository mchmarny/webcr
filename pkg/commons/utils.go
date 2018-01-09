package commons

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"net/url"
	"os"
	"regexp"
	"strings"
)

var (
	validURLPattern = regexp.MustCompile(`[^.]*\.[^.]{2,3}(?:\.[^.]{2,3})?$`)
)

// GetMD5 returns MD5 hash of the string
func GetMD5(s string) string {
	if s == "" {
		return s
	}
	hasher := md5.New()
	hasher.Write([]byte(strings.ToLower(s)))
	return hex.EncodeToString(hasher.Sum(nil))
}

// GetFileLines returns array of file lines
func GetFileLines(path string) (lines []string, err error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	rows := []string{}
	for scanner.Scan() {
		rows = append(rows, scanner.Text())
	}

	lineErr := scanner.Err()
	return rows, lineErr

}

// PathExists checks if file or dir exeists
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return !os.IsNotExist(err)
	}
	return true
}

// ParseDomain parses 2 or 3 octate domains (some.com or some.com.uk)
func ParseDomain(s string) string {

	// empty
	if s == "" {
		return s
	}

	// valid
	u, err := url.Parse(s)
	if err != nil {
		return s
	}

	// domain
	return validURLPattern.FindString(u.Host)

}
