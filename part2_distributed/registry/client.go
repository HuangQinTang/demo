package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
)

// RegisterService 服务注册
func RegisterService(r Registration) error {
	//解析url
	ServiceUpdateURL, err := url.Parse(r.ServiceUpdateURL)
	if err != nil {
		return err
	}
	//
	http.Handle(ServiceUpdateURL.Path, &serviceUpdateHanlder{})

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err = enc.Encode(r)
	if err != nil {
		return err
	}

	// 服务注册
	res, err := http.Post(ServicesUrl, "application/json", buf)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to register service. Registry service responded with code %s", res.StatusCode)
	}
	return nil
}

// ShutdownService 取消注册
func ShutdownService(url string) error {
	req, err := http.NewRequest(http.MethodDelete, ServicesUrl, bytes.NewBuffer([]byte(url)))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "text/plain")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to deregister service. Registry service responded with code %v", res.StatusCode)
	}
	return nil
}

// serviceUpdateHanlder providers的Hanlder
type serviceUpdateHanlder struct{}
func (suh serviceUpdateHanlder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	dec := json.NewDecoder(r.Body)
	var p patch
	err := dec.Decode(&p)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Printf("Updated received %v\n", p)
	prov.Update(p)	//当前服务提供的api维护进服务提供者
}

// providers 服务提供者
type providers struct {
	services map[ServiceName][]string //map[服务名称][]string{该服务下api的url...}
	mutex    *sync.RWMutex
}

var prov = providers{
	services: make(map[ServiceName][]string),
	mutex:    new(sync.RWMutex),
}

// Update 更新服务信息
func (p *providers) Update(pat patch) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	//添加服务
	for _, patchEntry := range pat.Added {
		if _, ok := p.services[patchEntry.Name]; !ok {
			p.services[patchEntry.Name] = make([]string, 0)
		}
		p.services[patchEntry.Name] = append(p.services[patchEntry.Name],
			patchEntry.URL)
	}

	//移除服务
	for _, patchEntry := range pat.Removed {
		if providerURLs, ok := p.services[patchEntry.Name]; ok {
			for i := range providerURLs {
				if providerURLs[i] == patchEntry.URL {
					p.services[patchEntry.Name] = append(providerURLs[:i],
						providerURLs[i+1:]...)
				}
			}
		}
	}
}

// 通过服务名称获取该服务下api url数组
func (p providers) get(name ServiceName) ([]string, error) {
	providers, ok := p.services[name]
	if !ok {
		return []string{}, fmt.Errorf("No providers available for service %v", name)
	}
	return providers, nil
}

func GetProvider(name ServiceName) ([]string, error) {
	return prov.get(name)
}