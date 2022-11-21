package parser

import (
	"crawler/distributed/config"
	"crawler/engine"
	"crawler/model"
	"regexp"
	"strconv"
)

//男士
const manProfile = `<a href="(http://album.zhenai.com/u/[0-9]+)" target="_blank">([^<]+)</a></th></tr> <tr><td width="180"><span class="grayL">性别：</span>男士</td> <td><span class="grayL">居住地：</span>([^<]+)</td></tr> <tr><td width="180"><span class="grayL">年龄：</span>([^<]+)</td> <!----> <td><span class="grayL">月   薪：</span>([^<]+)</td></tr> <tr><td width="180"><span class="grayL">婚况：</span>([^<]+)</td> <td width="180"><span class="grayL">身   高：</span>([^<]+)</td></tr></tbody></table> <div class="introduce">([^<]+)</div>`

//女士
const womanProfile = `<a href="(http://album.zhenai.com/u/[0-9]+)" target="_blank">([^<]+)</a></th></tr> <tr><td width="180"><span class="grayL">性别：</span>女士</td> <td><span class="grayL">居住地：</span>([^<]+)</td></tr> <tr><td width="180"><span class="grayL">年龄：</span>([^<]+)</td> <td><span class="grayL">学   历：</span>([^<]+)</td> <!----></tr> <tr><td width="180"><span class="grayL">婚况：</span>([^<]+)</td> <td width="180"><span class="grayL">身   高：</span>([^<]+)</td></tr></tbody></table> <div class="introduce">([^<]+)</div>`

//下一页
const cityListReNextPage = `<a href="(http://www.zhenai.com/zhenghun/[a-z]+/[2-6]+)">下一页</a>`

const userIdRe = `[0-9]+`

var (
	manMatches      = regexp.MustCompile(manProfile)
	womanMatches    = regexp.MustCompile(womanProfile)
	cityListMatches = regexp.MustCompile(cityListReNextPage)
	userIdMatche    = regexp.MustCompile(userIdRe)
)

// ParseSimpleInfo 简单信息解析器
func ParseSimpleInfo(contents []byte) engine.ParseResult {
	result := engine.ParseResult{}

	//男士信息
	men := manMatches.FindAllSubmatch(contents, -1)
	for _, val := range men {
		age, _ := strconv.Atoi(string(val[4]))
		height, _ := strconv.Atoi(string(val[7]))
		url := string(val[1])
		result.Items = append(result.Items, model.SimpleInfo{
			Id:        parseUserId(userIdMatche, url),
			Gender:    "男士",
			Url:       url,
			Nickname:  string(val[2]),
			Place:     string(val[3]),
			Age:       age,
			Income:    string(val[5]),
			Marriage:  string(val[6]),
			Height:    height,
			Introduce: string(val[8]),
		})
	}

	//女士信息
	women := womanMatches.FindAllSubmatch(contents, -1)
	for _, val := range women {
		age, _ := strconv.Atoi(string(val[4]))
		height, _ := strconv.Atoi(string(val[7]))
		url := string(val[1])
		result.Items = append(result.Items, model.SimpleInfo{
			Id:            parseUserId(userIdMatche, url),
			Gender:        "女士",
			Url:           string(val[1]),
			Nickname:      string(val[2]),
			Place:         string(val[3]),
			Age:           age,
			EducationMate: string(val[5]),
			Marriage:      string(val[6]),
			Height:        height,
			Introduce:     string(val[8]),
		})
	}

	//下一页信息解析
	pageNext := cityListMatches.FindSubmatch(contents)
	if len(pageNext) >= 2 {
		result.Requests = append(result.Requests, engine.Request{
			Url:    string(pageNext[1]),
			Parser: engine.NewFuncParser(ParseSimpleInfo, config.ParseSimpleInfo),
		})
	}

	return result
}

func parseUserId(matche *regexp.Regexp, url string) int {
	result := matche.FindSubmatch([]byte(url))
	userId, err := strconv.Atoi(string(result[0]))
	if err != nil {
		return 0
	}
	return userId
}
