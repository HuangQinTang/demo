package parser

import (
	"io/ioutil"
	"testing"
)

func TestParseCityList(t *testing.T) {
	contents, err := ioutil.ReadFile("citylist.txt")
	if err != nil {
		panic(err)
	}

	const resultSize = 470
	result := parseCityList(contents, resultSize)

	//验证城市列表是否470个 2022-06-10
	if len(result.Requests) != resultSize {
		t.Errorf("result should have %d requests; but had %d", resultSize, len(result.Requests))
	}

	expectedUrl := []string{"http://www.zhenai.com/zhenghun/aba", "http://www.zhenai.com/zhenghun/akesu", "http://www.zhenai.com/zhenghun/alashanmeng"}

	//遍历前三个url和city，判断与预期值是否相等 2022-06-14
	for i, url := range expectedUrl {
		if result.Requests[i].Url != url {
			t.Errorf("expected url #%d: %s; but was %s", i, url, result.Requests[i].Url)
		}
	}
	for i, Url := range expectedUrl {
		if result.Requests[i].Url != Url {
			t.Errorf("expected url %s; but was %s", Url, result.Requests[i].Url)
		}
	}
}
