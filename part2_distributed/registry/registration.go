package registry

type Registration struct {
	ServiceName ServiceName `json:"service_name"`
	ServiceUrl  string      `json:"service_url"`
}

type ServiceName string

const (
	LogService     = ServiceName("LogService")
	GradingService = ServiceName("GradingService")
)
