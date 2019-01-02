package main

import (
	"net/http"
	"flag"
	"net/url"
	"net/http/httputil"
	"log"
	"multi-proxy/service"
	"fmt"
	"runtime"
)

var cmd Cmd
var srv http.Server

func StartServer(proxyConfigs []service.ProxyConfig) {
	for _, proxy := range proxyConfigs {
		service.GetWebLoggerInstance().Info(fmt.Sprintf("Listening on %s, forwarding to %s", proxy.From, proxy.To))
		h := &handle{reverseProxy: proxy.To}
		srv := http.Server{}
		srv.Addr = proxy.From
		srv.Handler = h
		go func() {
			if err := srv.ListenAndServe(); err != nil {
				service.GetWebLoggerInstance().Info("ListenAndServe: ", err)
			} else {
				panic(err)
			}
		}()
	}
	select {}
}

func StopServer() {
	if err := srv.Shutdown(nil); err != nil {
		log.Println(err)
	}
}

func main() {
	defer Recover()
	cmd = parseCmd()
	service.SetConfigPath(cmd.configPath)
	proxys := service.LoadProxysConfig()
	StartServer(proxys.Proxy)
}

type Cmd struct {
	configPath string
}

func parseCmd() Cmd {
	var cmd Cmd
	flag.StringVar(&cmd.configPath, "c", "./config/dev/config.tml", "load config path")
	flag.Parse()
	return cmd
}

type handle struct {
	reverseProxy string
}

func (this *handle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	service.GetWebLoggerInstance().Info(r.RemoteAddr + " " + r.Method + " " + r.URL.String() + " " + r.Proto + " " + r.UserAgent())
	remote, err := url.Parse(this.reverseProxy)
	if err != nil {
		service.GetWebLoggerInstance().Error(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(w, r)
}

func Recover() {
	if err := recover(); err != nil {
		var stacktrace string
		for i := 1; ; i++ {
			_, f, l, got := runtime.Caller(i)
			if !got {
				break
			}
			stacktrace += fmt.Sprintf("%s:%d\n", f, l)
		}

		// when stack finishes
		logMessage := fmt.Sprintf("Trace: %s\n", err)
		logMessage += fmt.Sprintf("%s\n", stacktrace)
		service.GetWebLoggerInstance().Error(logMessage)
	}
}
