/**
一定要记得在confin.json配置这个模块的参数,否则无法使用
*/
package login

import (
	"fmt"
	"github.com/liangdas/mqant/conf"
	"github.com/liangdas/mqant/gate"
	"github.com/liangdas/mqant/module"
)

var Module = func() module.Module {
	gate := new(Login)
	return gate
}

type Login struct {
	module.BaseModule
}

func (m *Login) GetType() string {
	//很关键,需要与配置文件中的Module配置对应
	return "Login"
}
func (m *Login) Version() string {
	//可以在监控时了解代码版本
	return "1.0.0"
}
func (m *Login) OnInit(app module.App, settings *conf.ModuleSettings) {
	m.BaseModule.OnInit(m, app, settings)

	m.GetServer().RegisterGO("HD_Login", m.login)  //我们约定所有对客户端的请求都以Handler_开头
	m.GetServer().RegisterGO("getRand", m.getRand) //演示后台模块间的rpc调用
}

func (m *Login) Run(closeSig chan bool) {
}

func (m *Login) OnDestroy() {
	//一定别忘了关闭RPC
	m.GetServer().OnDestroy()
}

func (m *Login) login(s map[string]interface{}, msg map[string]interface{}) (result string, err string) {
	if msg["userName"] == nil || msg["passWord"] == nil {
		result = "userName or passWord cannot be nil"
		return
	}
	userName := msg["userName"].(string)
	//passWord:=msg["passWord"].(string)

	session := gate.NewSession(m.App, s)
	err = session.Bind(userName)
	if err != "" {
		return
	}
	session.Set("login", true)
	session.Push() //推送到网关
	return fmt.Sprintf("login success %s", userName), ""
}

func (m *Login) getRand(by []byte,mp map[string]interface{},f float64,i int,b bool) (result string, err string) {
	//演示后台模块间的rpc调用
	return fmt.Sprintf("My is Login Module %s", by,mp,f,i,b), ""
}
