package parser

import (
	"io/ioutil"
	"testing"
)

func TestParseCity(t *testing.T) {
	contentBts, err := ioutil.ReadFile("city.txt")
	if err != nil {
		panic(err)
	}
	correct := []string{
		"http://album.zhenai.com/u/1958903678",
		"http://album.zhenai.com/u/1396738538",
		"http://album.zhenai.com/u/1774343032",
		"http://album.zhenai.com/u/1529742129",
		"http://album.zhenai.com/u/1700301913",
		"http://album.zhenai.com/u/1398432707",
		"http://album.zhenai.com/u/1821376705",
		"http://album.zhenai.com/u/1491731990",
		"http://album.zhenai.com/u/1693313153",
		"http://album.zhenai.com/u/1761252184",
		"http://album.zhenai.com/u/1426975040",
		"http://album.zhenai.com/u/1677398857",
		"http://album.zhenai.com/u/1292525744",
		"http://album.zhenai.com/u/1834384219",
		"http://album.zhenai.com/u/1918103401",
		"http://album.zhenai.com/u/1122288128",
		"http://album.zhenai.com/u/1998372165",
		"http://album.zhenai.com/u/1839731206",
		"http://album.zhenai.com/u/1968078839",
		"http://album.zhenai.com/u/1977189312",
	}
	data := ParseCity(contentBts)
	for k, v := range data.Requests {
		if v.Url != correct[k] {
			t.Errorf("函数发生改变")
		}
	}
}
