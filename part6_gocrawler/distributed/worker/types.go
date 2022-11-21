package worker

import (
	"crawler/distributed/config"
	"crawler/engine"
	"crawler/model"
	"crawler/zhenai/parser"
	"errors"
	"fmt"
	"log"
)

type SerializedParser struct {
	Name string
	Args interface{}
}

// Request 可以在网络传播的engine.Request
type Request struct {
	Url    string
	Parser SerializedParser
}

type ParseResult struct {
	Items   []model.SimpleInfo
	Request []Request
}

// SerializeRequest 转化engine.Request序列化为可以在网络中传播的Request
func SerializeRequest(r engine.Request) Request {
	name, args := r.Parser.Serialize()
	return Request{
		Url: r.Url,
		Parser: SerializedParser{
			Name: name,
			Args: args,
		},
	}
}

// SerializeResult 转换engine.ParseResult序列化为可以在网络中传播的ParseResult
func SerializeResult(r engine.ParseResult) ParseResult {
	result := ParseResult{
		Items: r.Items,
	}
	for _, req := range r.Requests {
		result.Request = append(result.Request, SerializeRequest(req))
	}
	return result
}

// DeserializeRequest 将Request 反序列化得到engine.Request
func DeserializeRequest(r Request) (engine.Request, error) {
	parser, err := deserializeParser(r.Parser)
	if err != nil {
		return engine.Request{}, err
	}
	return engine.Request{
		Url:    r.Url,
		Parser: parser,
	}, nil
}

// DeserializeResult 将ParseResult 反序列化得到engine.ParseResult
func DeserializeResult(r ParseResult) engine.ParseResult {
	result := engine.ParseResult{
		Items: r.Items,
	}
	for _, req := range r.Request {
		engineReq, err := DeserializeRequest(req)
		if err != nil {
			log.Printf("error DeserializeResult request: %v", err)
			continue
		}
		result.Requests = append(result.Requests, engineReq)
	}
	return result
}

// deserializeParser 将SerializedParser 反序列化得到engine.Parser
func deserializeParser(p SerializedParser) (engine.Parser, error) {
	switch p.Name {
	case config.ParseCityList:
		if parseCityNum, ok := p.Args.(float64); ok {	//不晓得为什么传int会变为float64
			return parser.NewParseCityList(int(parseCityNum)), nil
		} else {
			return nil, fmt.Errorf("invalid arg: %v, required int type", p.Args)
		}
	case config.ParseSimpleInfo:
		return engine.NewFuncParser(parser.ParseSimpleInfo, config.ParseSimpleInfo), nil
	case config.NilParser:
		return engine.NilParser{}, nil
	default:
		return nil, errors.New("unknown parser name")
	}
}
