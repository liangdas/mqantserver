/**
一定要记得在confin.json配置这个模块的参数,否则无法使用
*/
package gate

import (
	"fmt"
	"github.com/liangdas/mqant/conf"
	"github.com/liangdas/mqant/gate/base"
	"github.com/liangdas/mqant/log"
	"github.com/liangdas/mqant/module"
	"github.com/liangdas/mqant/gate"
	"github.com/liangdas/mqant/gate/base/mqtt"
)

var Module = func() module.Module {
	gate := new(Gate)
	return gate
}

type Gate struct {
	basegate.Gate //继承
}

func (gate *Gate) GetType() string {
	//很关键,需要与配置文件中的Module配置对应
	return "Gate"
}
func (gate *Gate) Version() string {
	//可以在监控时了解代码版本
	return "1.0.0"
}
func (gate *Gate) OnInit(app module.App, settings *conf.ModuleSettings) {
	//注意这里一定要用 gate.Gate 而不是 module.BaseModule
	gate.Gate.OnInit(gate, app, settings)
	gate.Gate.SetStorageHandler(gate) //设置持久化处理器
	gate.Gate.SetTracingHandler(gate) //设置分布式跟踪系统处理器
}

/**
是否需要对本次客户端请求进行跟踪
*/
func (gate *Gate)OnRequestTracing(session gate.Session,msg *mqtt.Publish)bool{
	return true
}

/**
存储用户的Session信息
Session Bind Userid以后每次设置 settings都会调用一次Storage
*/
func (gate *Gate) Storage(Userid string, settings map[string]string) (err error) {
	log.Info("需要处理对Session的持久化")
	return nil
}

/**
强制删除Session信息
*/
func (gate *Gate) Delete(Userid string) (err error) {
	log.Info("需要删除Session持久化数据")
	return nil
}

/**
获取用户Session信息
用户登录以后会调用Query获取最新信息
*/
func (gate *Gate) Query(Userid string) (settings map[string]string, err error) {
	log.Info("查询Session持久化数据")
	return nil, fmt.Errorf("no redis")
}

/**
用户心跳,一般用户在线时60s发送一次
可以用来延长Session信息过期时间
*/
func (gate *Gate) Heartbeat(Userid string) {
	log.Info("用户在线的心跳包")
}
