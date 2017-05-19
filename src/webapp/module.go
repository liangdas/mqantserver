/**
一定要记得在confin.json配置这个模块的参数,否则无法使用
*/
package webapp

import (
	"fmt"
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

	self.GetServer().RegisterGO("mongodb", self.mongodb) //演示后台模块间的rpc调用
}
func someHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Normal Handler")
}
func (self *Web) Run(closeSig chan bool) {
	l, _ := net.Listen("tcp", ":8090")
	go func() {
		log.Info("webapp server Listen : %s", ":8090")
		root := mux.NewRouter()
		static:=root.PathPrefix("/static/")
		static.Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("/work/go/mqantserver/bin"))))
		sub:=root.PathPrefix("/sub/").Subrouter()
		sub.HandleFunc("/", someHandler)
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



func (self *Web) mongodb() (rpc_result string, rpc_err string) {


	return fmt.Sprintf("My is Login Module"), ""
}
