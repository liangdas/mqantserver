/**
一定要记得在confin.json配置这个模块的参数,否则无法使用
*/
package helloworld

import (
	"fmt"
	"github.com/liangdas/mqant/conf"
	"github.com/liangdas/mqant/gate"
	"github.com/liangdas/mqant/module"
	"github.com/liangdas/mqant/module/base"
)

var Module = func() module.Module {
	this := new(HellWorld)
	return this
}

type HellWorld struct {
	basemodule.BaseModule
}

func (m *HellWorld) GetType() string {
	//很关键,需要与配置文件中的Module配置对应
	return "HelloWorld"
}
func (m *HellWorld) Version() string {
	//可以在监控时了解代码版本
	return "1.0.0"
}
func (m *HellWorld) OnInit(app module.App, settings *conf.ModuleSettings) {
	m.BaseModule.OnInit(m, app, settings)

	m.GetServer().RegisterGO("HD_Say", m.say) //我们约定所有对客户端的请求都以HD_开头
}

func (m *HellWorld) Run(closeSig chan bool) {

}

func (m *HellWorld) OnDestroy() {
	//一定别忘了关闭RPC
	m.GetServer().OnDestroy()
}
func (m *HellWorld) say(session gate.Session, msg map[string]interface{}) (result string, err string) {
	if msg["say"] == nil {
		result = "say cannot be nil"
		return
	}
	say := msg["say"].(string)
	return fmt.Sprintf("you say : %s", say), ""
}
