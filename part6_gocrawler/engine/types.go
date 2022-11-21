package engine

import "crawler/model"

// Request 种子结构体
type Request struct {
	Url    string //目标url
	Parser Parser //解析器对象(内置解析器方法，解析url所返回数据)
}

// ParseResult 解析器返回结构体
type ParseResult struct {
	Items    []model.SimpleInfo //所解析的目标对象(这里是简单用户信息）
	Requests []Request          //目标对象下还需要处理的种子
}

// ParserFun 解析器方法
type ParserFun func(contents []byte) ParseResult

// Parser 解析器对象接口
type Parser interface {
	Parse(contents []byte) ParseResult          //解析方法
	Serialize() (name string, args interface{}) //序列化当前解析对象，用于rpc传递
}

type NilParser struct {
}

func (NilParser) Parse(_ []byte) ParseResult {
	return ParseResult{}
}

func (NilParser) Serialize() (name string, args interface{}) {
	return "NilParser", nil
}

// FuncParser 解析器对象
type FuncParser struct {
	parser ParserFun
	name   string
}

// Parse 解析contents（网页流），返回ParseResult（解析结果）
func (f *FuncParser) Parse(contents []byte) ParseResult {
	return f.parser(contents)
}

// Serialize 序列化方法
func (f *FuncParser) Serialize() (name string, args interface{}) {
	return f.name, nil
}

// NewFuncParser 解析器对象工场
func NewFuncParser(p ParserFun, name string) *FuncParser {
	return &FuncParser{
		parser: p,
		name:   name,
	}
}
