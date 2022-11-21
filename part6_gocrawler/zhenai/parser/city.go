package parser

import (
	"crawler/engine"
	"regexp"
)

const cityRe = `<a href="(http://album.zhenai.com/u/[0-9]+)"[^>]*>([^<]+)</a></th></tr> <tr><td width="180"><span class="grayL">性别：</span>(.{2})</td>`

// ParseCity 城市解析器
func ParseCity(contents []byte) engine.ParseResult {
	re := regexp.MustCompile(cityRe)
	matches := re.FindAllSubmatch(contents, -1) //-1,匹配所有
	result := engine.ParseResult{}
	for _, m := range matches {
		url := string(m[1])
		nickname := string(m[2])
		gender := string(m[3])
		//result.Items = append(result.Items, string(m[2]))
		result.Requests = append(result.Requests, engine.Request{
			Url: string(m[1]),
			Parser: engine.NewFuncParser(func(bytes []byte) engine.ParseResult {
				return ParseProfile(contents, url, nickname, gender)
			}, "ParseProfile"),
		})
	}
	return result
}
