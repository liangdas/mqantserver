/**
一定要记得在confin.json配置这个模块的参数,否则无法使用
*/
package webapp

import (
	"github.com/liangdas/mqant/log"
	"github.com/liangdas/mqant/conf"
	"github.com/liangdas/mqant/module"
	"github.com/gorilla/mux"
	"net/http"
	"net"
	"github.com/liangdas/mqant/module/base"
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
func (self *Web) Run(closeSig chan bool) {
	//这里如果出现异常请检查8080端口是否已经被占用
	l, err := net.Listen("tcp", ":8080")
	if err!=nil{
		log.Error("webapp server error",err.Error())
		return
	}
	go func() {
		log.Info("webapp server Listen : %s", ":8080")
		root := mux.NewRouter()
		static:=root.PathPrefix("/mqant/")
		static.Handler(http.StripPrefix("/mqant/", http.FileServer(http.Dir(self.GetModuleSettings().Settings["StaticPath"].(string)))))
		//r.Handle("/static",static)
		ServeMux:=http.NewServeMux()
		ServeMux.Handle("/", root)
		http.Serve(l, ServeMux)
	}()
	<-closeSig
	log.Info("webapp server Shutting down...")
	l.Close()
}

func (self *Web) OnDestroy() {
	//一定别忘了关闭RPC
	self.GetServer().OnDestroy()
}

