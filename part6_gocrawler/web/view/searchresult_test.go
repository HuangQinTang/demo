package view

import (
	"crawler/web/model"
	common "crawler/model"
	"os"
	"testing"
)

func TestSearchResultView_Render(t *testing.T) {
	//读取模板数据
	view := CreateSearchResultView("template.html")
	page := model.SearchResult{
		Hits: 2,
		Start: 0,
		Items: []common.SimpleInfo{
			{Url: "http://album.zhenai.com/u/1693313153",Nickname: "非诚勿扰",Gender: "男士",Place: "四川阿坝",Age: 34,Income: "20001-50000元", EducationMate: "大专", Marriage: "未婚", Height: 180,Introduce: "从头到尾一定要诚！只有真诚换真爱，能遇见就且行且珍惜。"},
			{Url: "http://album.zhenai.com/u/1693313153",Nickname: "非诚勿扰",Gender: "男士",Place: "四川阿坝",Age: 34,Income: "20001-50000元", EducationMate: "大专", Marriage: "未婚", Height: 180,Introduce: "从头到尾一定要诚！只有真诚换真爱，能遇见就且行且珍惜。"},
		},
	}

	//创建文件
	out, err := os.Create("template.test.html")
	defer out.Close()
	if err != nil {
		t.Errorf(err.Error())
	}

	//写渲染后的模板数据至创建的文件
	err = view.Render(out, page)
	if err != nil {
		t.Errorf(err.Error())
	}
}
