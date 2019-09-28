/**
一定要记得在confin.json配置这个模块的参数,否则无法使用
*/
package webapp

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/liangdas/mqant/conf"
	"github.com/liangdas/mqant/log"
	"github.com/liangdas/mqant/module"
	"github.com/liangdas/mqant/module/base"
	"net"
	"net/http"
	"time"
)

var Module = func() *Web {
	web := new(Web)
	return web
}

type Web struct {
	basemodule.BaseModule
}

func (self *Web) GetType() string {
	//很关键,需要与配置文件中的Module配置对应
	return "Webapp"
}
func (self *Web) Version() string {
	//可以在监控时了解代码版本
	return "1.0.0"
}
func (self *Web) OnInit(app module.App, settings *conf.ModuleSettings) {
	self.BaseModule.OnInit(self, app, settings)
}

func loggingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		//[26/Oct/2017:19:07:04 +0800]`-`"GET /g/c HTTP/1.1"`"curl/7.51.0"`502`[127.0.0.1]`-`"-"`0.006`166`-`-`127.0.0.1:8030`-`0.000`xd
		log.Info("%s %s %s [%s] in %v", r.Method, r.URL.Path, r.Proto, r.RemoteAddr, time.Since(start))
	})
}
func Statushandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
	})
}
func (self *Web) Run(closeSig chan bool) {
	//这里如果出现异常请检查8080端口是否已经被占用
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Error("webapp server error", err.Error())
		return
	}
	go func() {
		log.Info("webapp server Listen : %s", ":8080")
		root := mux.NewRouter()
		status := root.PathPrefix("/status")
		status.HandlerFunc(Statushandler)

		static := root.PathPrefix("/mqant/")
		static.Handler(http.StripPrefix("/mqant/", http.FileServer(http.Dir(self.GetModuleSettings().Settings["StaticPath"].(string)))))
		//r.Handle("/static",static)
		ServeMux := http.NewServeMux()
		ServeMux.Handle("/", root)
		http.Serve(l, loggingHandler(ServeMux))
	}()
	<-closeSig
	log.Info("webapp server Shutting down...")
	l.Close()
}

func (self *Web) OnDestroy() {
	//一定别忘了关闭RPC
	self.GetServer().OnDestroy()
}
