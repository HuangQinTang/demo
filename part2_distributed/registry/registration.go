package registry

// Registration 服务注册请求结构
type Registration struct {
	ServiceName      ServiceName   `json:"service_name"` //服务名称
	ServiceUrl       string        `json:"service_url"`  //服务url
	RequiredServices []ServiceName //当前服务所依赖的其他服务
	ServiceUpdateURL string        //服务中心回调本服务地址
}

type ServiceName string

const (
	LogService     = ServiceName("LogService")     //日志服务名
	GradingService = ServiceName("GradingService") //学生信息服务名
)

//服务变更请求结构
type patchEntry struct {
	Name ServiceName
	URL  string
}

type patch struct {
	Added   []patchEntry	//增加的服务
	Removed []patchEntry	//减少的服务
}
