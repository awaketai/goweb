package command

import (
	"errors"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/awaketai/goweb/framework"
	"github.com/awaketai/goweb/framework/contract"
)

// devConfig 调试模式配置信息
type devConfig struct {
	Port string // 调式模式监听的端口，默认8070
	// Backend 后端调式模式配置
	Backend struct {
		RefreshTime   int    // 调式模式后端更新时间，单位s
		Port          string // 后端监听端口，默认8072
		MonitorFolder string // 监听文件夹
	}
	// Frontend 前端调工模式配置
	Frontend struct {
		Port string // 前端启动端口，默认8071
	}
}

const (
	devConfigDefaultPort = "8087"
	backendDefaultPort   = "8072"
	frontendDefaultPort  = "8071"
)

func initDevConfig(container framework.Container) *devConfig {
	devConfig := &devConfig{
		Port: devConfigDefaultPort,
		Backend: struct {
			RefreshTime   int
			Port          string
			MonitorFolder string
		}{
			1, backendDefaultPort, "",
		},
		Frontend: struct{ Port string }{
			frontendDefaultPort,
		},
	}

	configer := container.MustMake(contract.ConfigKey).(contract.Config)
	if configer.IsExists("app.dev.port") {
		devConfig.Port = configer.GetString("app.dev.port")
	}
	if configer.IsExists("app.dev.backend.refresh_time") {
		devConfig.Backend.RefreshTime = configer.GetInt("app.dev.backend.refresh_time")
	}
	if configer.IsExists("app.dev.backend.port") {
		devConfig.Backend.Port = configer.GetString("app.dev.backend.port")
	}
	monitorFolder := configer.GetString("app.dev.backend.monitor_folder")
	if monitorFolder == "" {
		appService := container.MustMake(contract.AppKey).(contract.App)
		devConfig.Backend.MonitorFolder = appService.AppFolder()
	}
	if configer.IsExists("app.dev.frontend.port") {
		devConfig.Frontend.Port = configer.GetString("app.dev.frontend.port")
	}
	return devConfig
}

// Proxy serve启动的服务代理
type Proxy struct {
	devConfig   *devConfig
	proxyServer *http.Server
	backendPid  int
	frontendPid int
}

func NewProxy(container framework.Container) *Proxy {
	devConfig := initDevConfig(container)
	return &Proxy{devConfig: devConfig}
}

func (p *Proxy) newProxyReverseProxy(frontend, backend *url.URL) *httputil.ReverseProxy {
	if p.frontendPid == 0 && p.backendPid == 0 {
		log.Println("backend and frontend server not exists both")
		return nil
	}

	if p.frontendPid == 0 && p.backendPid != 0 {
		return httputil.NewSingleHostReverseProxy(backend)
	}

	if p.backendPid == 0 && p.frontendPid != 0 {
		return httputil.NewSingleHostReverseProxy(frontend)
	}

	// 两个都有进程
	director := func(req *http.Request) {
		req.URL.Scheme = backend.Scheme
		req.URL.Host = backend.Host
	}
	NotFoundErr := errors.New("response is 404,need to redirect")
	return &httputil.ReverseProxy{
		Director: director,
		ModifyResponse: func(r *http.Response) error {
			if r.StatusCode == 404 {
				return NotFoundErr
			}
			return nil
		},
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			if errors.Is(err, NotFoundErr) {
				httputil.NewSingleHostReverseProxy(frontend).ServeHTTP(w, r)
			}
		},
	}
}
