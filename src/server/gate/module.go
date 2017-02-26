/**
一定要记得在confin.json配置这个模块的参数,否则无法使用
 */
package gate
import (
	"github.com/liangdas/mqant/gate"
	"github.com/liangdas/mqant/conf"
	"github.com/liangdas/mqant/module"
	"fmt"
	"github.com/liangdas/mqant/log"
)
var Module = func() (module.Module){
	gate := new(Gate)
	return gate
}

type Gate struct {
	gate.Gate
}
func (gate *Gate) GetType()(string){
	//很关键,需要与配置文件中的Module配置对应
	return "Gate"
}

func (gate *Gate) OnInit(app module.App,settings *conf.ModuleSettings) {
	//注意这里一定要用 gate.Gate 而不是 gate.BaseModule
	gate.Gate.OnInit(gate,app,settings)
	gate.Gate.SetStorageHandler(gate)	//设置持久化处理器
}

/**
	存储用户的Session信息
	Session Bind Userid以后每次设置 settings都会调用一次Storage
	 */
func (gate *Gate) Storage(Userid string,settings map[string]interface{})(err error){
	log.Debug("对Session持久化")
	return nil
}
/**
强制删除Session信息
 */
func (gate *Gate) Delete(Userid string)(err error){
	log.Debug("删除Session持久化数据")
	return nil
}
/**
获取用户Session信息
用户登录以后会调用Query获取最新信息
 */
func (gate *Gate) Query(Userid string)(settings map[string]interface{},err error){
	log.Debug("查询Session持久化数据")
	return nil,fmt.Errorf("no redis")
}
/**
用户心跳,一般用户在线时1s发送一次
可以用来延长Session信息过期时间
 */
func (gate *Gate) Heartbeat(Userid string){
	log.Debug("用户在线心跳包")
}
