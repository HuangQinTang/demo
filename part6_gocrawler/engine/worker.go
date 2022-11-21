package engine

import (
	"crawler/fetcher"
	"log"
)

// Worker 拉取网页数据并解析
func Worker(r Request) (ParseResult, error) {
	//判断是否已拉取过
	//unique, err := rediservice.RedisServer.ZhenaiDuplicate(r.Url)
	//if err != nil {
	//	log.Printf("Fetcher: error fetching url %s: %v", r.Url, err)
	//	return ParseResult{}, err
	//}
	//
	//if !unique {
	//	log.Printf("Fetcher: error fetching url %s: %s", r.Url, "已经拉取")
	//	return ParseResult{}, errors.New("已拉取过")
	//}

	log.Printf("Fetching %s", r.Url)
	body, err := fetcher.Fetch(r.Url)
	if err != nil {
		log.Printf("Fetcher: error fetching url %s: %v", r.Url, err)
		return ParseResult{}, err
	}

	return r.Parser.Parse(body), nil
}
