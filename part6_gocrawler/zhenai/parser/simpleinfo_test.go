package parser

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestParseSimpleInfo(t *testing.T) {
	contents, err := ioutil.ReadFile("simpleinfo.txt")
	if err != nil {
		panic(err)
	}

	result := ParseSimpleInfo(contents)
	if len(result.Items) < 1 {
		t.Errorf("解析失败")
	}
	for _, info := range result.Items {
		if info.Id != 1396738538 {
			t.Errorf("数据解析不正确")
		}
		fmt.Printf("%+v", info)
		break
	}
}
