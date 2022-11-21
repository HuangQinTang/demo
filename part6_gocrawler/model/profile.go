package model

const Zhenai = "zhenai"

type Profile struct {
	Url          string
	NickName     string //昵称
	Gender       string //性别
	Sign         string //内心独白
	BaseInfo            //基本信息
	Living              //生活状况
	HobbyProfile        //爱好
	ChooseMate          //择偶条件
}

// BaseInfo 基本信息
type BaseInfo struct {
	Marriage   string //婚姻状况
	Age        int    //年龄
	Xinzuo     string //星座
	Height     int    //身高
	Weight     int    //体重
	Workplace  string //工作地
	Income     string //收入
	Occupation string //职业
	Education  string //教育
}

// Living 生活状况
type Living struct {
	Nationality  string //民族
	Hokou        string //户口
	Stature      string //身材
	Smoking      string //是否吸烟
	Drink        string //是否喝酒
	House        string //是否购房
	Car          string //是否购车
	Child        string //是否有孩子
	LikeChild    string //是否想要孩子
	MarriageTime string //结婚时机
}

// HobbyProfile 兴趣爱好
type HobbyProfile struct {
	Food  string //喜欢的食物
	Song  string //喜欢的歌
	Hobby string //爱好
	Idol  string //偶像
	Book  string //喜欢的书
}

// ChooseMate 择偶条件
type ChooseMate struct {
	AgeMate       string //年龄
	HeightMate    string //身高
	WorkplaceMate string //工作地
	IncomeMate    string //月薪
	EducationMate string //教育
	StatureMate   string //身材
	SmokingMate   string //是否吸烟
	DrinkMate     string //是否喝酒
	ChildMate     string //是否有孩子
	LikeChildMate string //是否想要孩子
}

// SimpleInfo 简单人物信息
type SimpleInfo struct {
	Id            int    //用户id
	Url           string //个人信息url
	Nickname      string //昵称
	Gender        string //性别
	Place         string //居住地
	Age           int    //年龄
	Income        string //收入（男性能抓到）
	EducationMate string //教育（女性能抓到）
	Marriage      string //婚姻状况
	Height        int    //身高
	Introduce     string //自我介绍
}
