package commons

import (
	"strings"
	"testing"
)

func TestURLID(t *testing.T) {

	u1 := "https://sub.domain.com/folder/image.jpg?q=1&v=a"
	id1 := GetMD5(u1)
	id2 := GetMD5(u1)

	if id1 != id2 {
		t.Fatal("Inconsistent IDs generated from same URL")
	}

	u3 := strings.ToUpper(u1)
	id3 := GetMD5(u3)

	if id1 != id3 {
		t.Fatal("Inconsistent IDs generated from same URL in different case")
	}

}
