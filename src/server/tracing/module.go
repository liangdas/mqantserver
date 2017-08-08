/**
一定要记得在confin.json配置这个模块的参数,否则无法使用
*/
package tracing

import (
	"net"
	"github.com/liangdas/mqant/log"
	"github.com/liangdas/mqant/conf"
	"github.com/liangdas/mqant/module"
	"github.com/liangdas/mqant/module/base"
	"time"
	"fmt"
	"runtime"
)

var Module = func() module.Module {
	this := new(tracing)
	return this
}

type tracing struct {
	basemodule.BaseModule
}

func (self *tracing) GetType() string {
	//很关键,需要与配置文件中的Module配置对应
	return "Tracing"
}
func (self *tracing) Version() string {
	//可以在监控时了解代码版本
	return "1.0.0"
}
func (self *tracing) OnInit(app module.App, settings *conf.ModuleSettings) {
	self.BaseModule.OnInit(self, app, settings)

}


func (self *tracing) Run(closeSig chan bool) {
	// switch 类似 if 可以带上一个短语句
	StoreFile:="/tmp/appdash.gob"
	switch os := runtime.GOOS; os {
	case "darwin":

	case "linux":
	case "windows":
		StoreFile="c://appdash.gob"
	default:
		// freebsd, openbsd,
		// plan9, windows...
		fmt.Printf("%s.", os)
	}
	cmd:=&ServeCmd{
		URL           :"http://localhost:7700",
		CollectorAddr :":7701",
		HTTPAddr      :":7700",

		StoreFile     :StoreFile,
		PersistInterval	: time.Second*2,

		Debug 	:false,
		Trace	:false,

		DeleteAfter :time.Second*60*10,

		LimitMax :10,
	}
	l, err := net.Listen("tcp",cmd.HTTPAddr)
	if err != nil {
		log.Info("tracing server ",err.Error())
	}else{
		go func() {
			err:=cmd.Execute(l)
			if err!=nil{
				log.Error("tracing server ",err.Error())
			}
		}()
	}
	<-closeSig
	log.Info("tracing server Shutting down...")
	if l!=nil{
		l.Close()
	}
}

func (self *tracing) OnDestroy() {
	//一定别忘了关闭RPC
	self.GetServer().OnDestroy()
}
