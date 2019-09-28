/**
一定要记得在confin.json配置这个模块的参数,否则无法使用
*/
package user

import (
	"fmt"
	"github.com/liangdas/mqant/conf"
	"github.com/liangdas/mqant/module"
	"github.com/liangdas/mqant/module/base"
)

var Module = func() module.Module {
	user := new(User)
	return user
}

type User struct {
	basemodule.BaseModule
}

func (self *User) GetType() string {
	//很关键,需要与配置文件中的Module配置对应
	return "User"
}
func (self *User) Version() string {
	//可以在监控时了解代码版本
	return "1.0.0"
}
func (self *User) OnInit(app module.App, settings *conf.ModuleSettings) {
	self.BaseModule.OnInit(self, app, settings)

	self.GetServer().RegisterGO("mongodb", self.mongodb) //演示后台模块间的rpc调用
}

func (self *User) Run(closeSig chan bool) {
}

func (self *User) OnDestroy() {
	//一定别忘了关闭RPC
	self.GetServer().OnDestroy()
}

func (self *User) mongodb() (rpc_result string, rpc_err string) {

	return fmt.Sprintf("My is Login Module"), ""
}
