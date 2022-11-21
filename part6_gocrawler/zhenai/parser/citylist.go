package parser

import (
	"crawler/distributed/config"
	"crawler/engine"
	"regexp"
)

const cityListRe = `<a href="(http://www.zhenai.com/zhenghun/[0-9a-z]+)" data-v-1573aa7c>([^<]+)</a>`

// ParseCityList 城市列表解析器解析器
// num 爬取的城市数量
func parseCityList(contents []byte, num int) engine.ParseResult {
	re := regexp.MustCompile(cityListRe)
	matches := re.FindAllSubmatch(contents, -1) //-1,匹配所有

	result := engine.ParseResult{}
	i := 0
	for _, m := range matches {
		i++

		//result.Items = append(result.Items, "City "+string(m[2]))
		result.Requests = append(result.Requests, engine.Request{
			Url:    string(m[1]), //url
			Parser: engine.NewFuncParser(ParseSimpleInfo, config.ParseSimpleInfo),
		})

		if i >= num {
			break
		}
	}
	return result
}

type ParseCityList struct {
	ParseCityNum int
}

func (p *ParseCityList) Parse(contents []byte) engine.ParseResult {
	return parseCityList(contents, p.ParseCityNum)
}

func (p *ParseCityList) Serialize() (name string, args interface{}) {
	return "ParseCityList", p.ParseCityNum
}

// NewParseCityList 城市列表解析器工场，传入要拉取的城市数
func NewParseCityList(parseCityNum int) *ParseCityList {
	return &ParseCityList{
		ParseCityNum: parseCityNum,
	}
}
