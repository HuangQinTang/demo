package registry

// Registration 服务注册请求结构
type Registration struct {
	ServiceName      ServiceName   `json:"service_name"` //服务名称
	ServiceUrl       string        `json:"service_url"`  //服务url
	RequiredServices []ServiceName //当前服务所依赖的其他服务
	ServiceUpdateURL string        //服务中心回调本服务地址
	HeartbeatURL     string        //心跳地址
}

type ServiceName string

// 服务名称，向服务提供者注册服务名称时，请定义在这；获取依赖服务时，也从这儿取
const (
	LogService     = ServiceName("LogService")     //日志服务名
	GradingService = ServiceName("GradingService") //学生信息服务名
	PortalService  = ServiceName("PortalService")  //学生信息的web服务
)

//服务变更请求结构
type patchEntry struct {
	Name ServiceName
	URL  string
}

type patch struct {
	Added   []patchEntry //增加的服务
	Removed []patchEntry //减少的服务
}
