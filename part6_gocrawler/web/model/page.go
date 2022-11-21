package model

import "crawler/model"

type SearchResult struct {
	Hits  int64  //总共找到多少
	Start int    //从第几个开始
	Query string //搜索条件
	PrevFrom int //上一页
	NextFrom int //下一页
	Items []model.SimpleInfo
}
