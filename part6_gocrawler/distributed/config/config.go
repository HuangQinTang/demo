package config

const (
	// Service ports
	ItemSaverPort = 8889
	WorkerPort0   = 9000

	// ElasticSearch
	Zhenai = "zhenai"

	// RPC Endpoints
	ItemSaverRpc    = "ItemSaverService.Save"
	CrawlServiceRpc = "CrawlService.Process"

	// Parser names
	ParseCityList   = "ParseCityList"
	ParseSimpleInfo = "ParseSimpleInfo"
	NilParser       = "NilParser"

	//Fetch Qps
	Qps = 10
)
