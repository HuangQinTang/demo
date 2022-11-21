package parser

import (
	"crawler/engine"
	"crawler/model"
	"regexp"
	"strconv"
	"strings"
)

const (
	baseInfoRe = `<div class="m-btn purple" data-v-8b1eac0c>([^<]+)</div>`
	livingRe   = `<div class="m-btn pink" data-v-8b1eac0c>([^<]+)</div>`
	foodRe     = `<div class="question f-fl" data-v-8b1eac0c>喜欢的一道菜：</div> <div class="answer f-fl" data-v-8b1eac0c>([^<]+)</div>`
	idolRe     = `<div class="question f-fl" data-v-8b1eac0c>欣赏的一个名人：</div> <div class="answer f-fl" data-v-8b1eac0c>([^<]+)</div>`
	songRe     = `<div class="question f-fl" data-v-8b1eac0c>喜欢的一首歌：</div> <div class="answer f-fl" data-v-8b1eac0c>([^<]+)</div>`
	bookRe     = `<div class="question f-fl" data-v-8b1eac0c>喜欢的一本书：</div> <div class="answer f-fl" data-v-8b1eac0c>([^<]+)</div>`
	thingRe    = `<div class="question f-fl" data-v-8b1eac0c>喜欢做的事：</div> <div class="answer f-fl" data-v-8b1eac0c>([^<]+)</div>`
	mateRe     = `<div class="m-btn" data-v-8b1eac0c>([^<]+)</div>`
	signRe     = `<span data-v-8b1eac0c>([^<]+)</span>`
)

var (
	baseInfoMatches = regexp.MustCompile(baseInfoRe)
	livingMatches   = regexp.MustCompile(livingRe)
	foodMatch       = regexp.MustCompile(foodRe)
	idolMatch       = regexp.MustCompile(idolRe)
	songMatch       = regexp.MustCompile(songRe)
	bookMatch       = regexp.MustCompile(bookRe)
	thingMatch      = regexp.MustCompile(thingRe)
	mateMatch       = regexp.MustCompile(mateRe)
	signMatch       = regexp.MustCompile(signRe)
)

const (
	Marriage = iota
	Age
	Xinzuo
	Height
	Weight
	Workplace
	Income
	Occupation
	Education
)

const (
	Nationality = iota
	Hokou
	Stature
	Smoking
	Drink
	House
	Car
	Child
	likeChild
	MarriageTime
)

func ParseProfile(contents []byte, url, nickname, gender string) engine.ParseResult {
	profile := model.Profile{Url: url, NickName: nickname, Gender: gender}

	//基本信息
	profile.BaseInfo = getBaseInfo(contents)
	//生活状况
	profile.Living = getLiving(contents)
	//兴趣爱好
	profile.HobbyProfile = getHobby(contents)
	//择偶条件
	profile.ChooseMate = getMateCondition(contents)
	//内心独白
	profile.Sign = extractString(contents, signMatch)

	//fmt.Println("---成功抓取---")
	//fmt.Printf("%#v", profile)
	return engine.ParseResult{
		//Items:    []interface{}{profile},
		Requests: []engine.Request{},
	}
}

// getBaseInfo 基本信息
func getBaseInfo(contents []byte) model.BaseInfo {
	baseInfo := model.BaseInfo{}
	matches := baseInfoMatches.FindAllSubmatch(contents, -1)
	if len(matches) == 9 { //信息完整
		for k, val := range matches {
			switch k {
			case Marriage:
				baseInfo.Marriage = string(val[1])
			case Age:
				ageBts := val[1][0 : len(val[1])-3] //去除末尾中文"岁"
				age, err := strconv.Atoi(string(ageBts))
				if err == nil {
					baseInfo.Age = age
				}
			case Xinzuo:
				baseInfo.Xinzuo = string(val[1])
			case Height:
				heightBts := val[1][0 : len(val[1])-2] //去除末尾"cm"
				height, err := strconv.Atoi(string(heightBts))
				if err == nil {
					baseInfo.Height = height
				}
			case Weight:
				weightBts := val[1][0 : len(val[1])-2] //去除末尾"kg"
				weight, err := strconv.Atoi(string(weightBts))
				if err == nil {
					baseInfo.Weight = weight
				}
			case Workplace:
				workplace := val[1][10:] //去除开头"工作地:"
				baseInfo.Workplace = string(workplace)
			case Income:
				income := val[1][10:] //去除开头"月收入:"
				baseInfo.Income = string(income)
			case Occupation:
				baseInfo.Occupation = string(val[1])
			case Education:
				baseInfo.Education = string(val[1])
			}
		}
	} else { //信息不完整判断字符串解析
		for _, v := range matches {
			val := string(v[1])
			//婚姻状况
			if strings.Contains(val, "离异") || strings.Contains(val, "丧偶") || strings.Contains(val, "未婚") {
				baseInfo.Marriage = val
				continue
			}
			//年龄
			if strings.Contains(val, "岁") {
				ageBts := v[1][0 : len(v[1])-3] //去除末尾中文"岁"
				age, err := strconv.Atoi(string(ageBts))
				if err == nil {
					baseInfo.Age = age
				}
				continue
			}
			//星座
			if strings.Contains(val, "座") {
				baseInfo.Xinzuo = val
				continue
			}
			//身高
			if strings.Contains(val, "cm") {
				heightBts := v[1][0 : len(v[1])-2] //去除末尾"cm"
				height, err := strconv.Atoi(string(heightBts))
				if err == nil {
					baseInfo.Height = height
				}
				continue
			}
			//体重
			if strings.Contains(val, "kg") {
				weightBts := v[1][0 : len(v[1])-2] //去除末尾"kg"
				weight, err := strconv.Atoi(string(weightBts))
				if err == nil {
					baseInfo.Weight = weight
				}
				continue
			}
			//工作地
			if strings.Contains(val, "工作地") {
				res := v[1][10:] //去除开头"工作地:"
				baseInfo.Workplace = string(res)
				continue
			}
			//收入
			if strings.Contains(val, "收入") {
				income := v[1][10:] //去除开头"月收入:"
				baseInfo.Income = string(income)
				continue
			}
			//教育
			if strings.Contains(val, "小学") || strings.Contains(val, "初中") || strings.Contains(val, "高中") ||
				strings.Contains(val, "专") || strings.Contains(val, "本科") || strings.Contains(val, "硕") ||
				strings.Contains(val, "博士") {
				baseInfo.Education = val
				continue
			}
			//都没有匹配到认为是职业
			baseInfo.Occupation = val
		}
	}
	return baseInfo
}

// getLiving 生活状况
func getLiving(contents []byte) model.Living {
	living := model.Living{}
	matches := livingMatches.FindAllSubmatch(contents, -1)
	if len(matches) == 10 { //信息完整
		for k, val := range matches {
			switch k {
			case Nationality:
				living.Nationality = string(val[1])
			case Hokou:
				hokou := val[1][7:] //去除开头的"籍贯"
				living.Hokou = string(hokou)
			case Stature:
				stature := val[1][7:] //去除开头"体型"
				living.Stature = string(stature)
			case Smoking:
				living.Smoking = string(val[1])
			case Drink:
				living.Drink = string(val[1])
			case House:
				living.House = string(val[1])
			case Car:
				living.Car = string(val[1])
			case Child:
				living.Child = string(val[1])
			case likeChild:
				likeChildBts := val[1][19:] //去除开头"是否想要孩子:"
				living.LikeChild = string(likeChildBts)
			case MarriageTime:
				marriageTime := val[1][13:]
				living.MarriageTime = string(marriageTime)
			}
		}
	} else { //信息不完整时
		for _, v := range matches {
			val := string(v[1])
			//民族
			if strings.Contains(val, "族") {
				living.Nationality = val
				continue
			}
			//户口
			if strings.Contains(val, "籍贯") {
				living.Hokou = val
				continue
			}
			//体型
			if strings.Contains(val, "体型") {
				living.Stature = val
				continue
			}
			//是否吸烟
			if strings.Contains(val, "烟") {
				living.Smoking = val
				continue
			}
			//是否喝酒
			if strings.Contains(val, "酒") {
				living.Drink = val
				continue
			}
			//是否购房
			if strings.Contains(val, "房") || strings.Contains(val, "租") || strings.Contains(val, "住") || strings.Contains(val, "宿舍") {
				living.House = val
				continue
			}
			//是否购车
			if strings.Contains(val, "车") {
				living.Car = val
				continue
			}
			//是否有孩子
			if strings.Contains(val, "孩子") && !strings.Contains(val, "是否想要孩子") {
				living.Child = val
				continue
			}
			//是否想要孩子
			if strings.Contains(val, "是否想要孩子") {
				likeChildBts := v[1][19:] //去除开头"是否想要孩子:"
				living.LikeChild = string(likeChildBts)
				continue
			}
			//合适结婚
			if strings.Contains(val, "何时结婚") {
				living.MarriageTime = string(val)
				continue
			}
		}
	}
	return living
}

// getHobby 兴趣爱好
func getHobby(contents []byte) model.HobbyProfile {
	hoby := model.HobbyProfile{}
	hoby.Food = extractString(contents, foodMatch)   //喜欢的食物
	hoby.Song = extractString(contents, songMatch)   //喜欢的歌
	hoby.Idol = extractString(contents, idolMatch)   //喜欢的名人
	hoby.Book = extractString(contents, bookMatch)   //喜欢的书
	hoby.Hobby = extractString(contents, thingMatch) //喜欢做的事
	return hoby
}

func extractString(contents []byte, re *regexp.Regexp) string {
	match := re.FindSubmatch(contents)
	if len(match) >= 2 {
		return string(match[1])
	}
	return ""
}

// getMateCondition 解析择偶条件
func getMateCondition(contents []byte) model.ChooseMate {
	matches := mateMatch.FindAllSubmatch(contents, -1)
	mate := model.ChooseMate{}
	for _, v := range matches {
		val := string(v[1])
		//年龄
		if strings.Contains(val, "岁") {
			mate.AgeMate = val
			continue
		}
		//身高
		if strings.Contains(val, "cm") {
			mate.HeightMate = val
			continue
		}
		//工作地
		if strings.Contains(val, "工作地") {
			res := v[1][10:] //去除开头"工作地:"
			mate.WorkplaceMate = string(res)
			continue
		}
		//教育
		if strings.Contains(val, "小学") || strings.Contains(val, "初中") || strings.Contains(val, "高中") ||
			strings.Contains(val, "专") || strings.Contains(val, "本科") || strings.Contains(val, "硕") ||
			strings.Contains(val, "博士") {
			mate.EducationMate = val
			continue
		}
		//薪资
		if strings.Contains(val, "月薪") {
			mate.IncomeMate = val
			continue
		}
		//身材
		if strings.Contains(val, "体型") {
			mate.StatureMate = val
			continue
		}
		//是否吸烟
		if strings.Contains(val, "烟") {
			mate.SmokingMate = val
			continue
		}
		//是否喝酒
		if strings.Contains(val, "酒") {
			mate.DrinkMate = val
			continue
		}
		//是否有孩子
		if strings.Contains(val, "孩子") && !strings.Contains(val, "是否想要孩子") {
			mate.ChildMate = val
			continue
		}
		//是否想要孩子
		if strings.Contains(val, "是否想要孩子") {
			likeChildBts := v[1][19:] //去除开头"是否想要孩子:"
			mate.LikeChildMate = string(likeChildBts)
			continue
		}
	}
	return mate
}
