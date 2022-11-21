package fetcher

import (
	"fmt"
	"testing"
)

func TestFetch(t *testing.T) {
	contents, err := Fetch("http://album.zhenai.com/u/1958903678")
	if err != nil {
		t.Errorf("%s", err)
	}
	fmt.Println(string(contents))
}
