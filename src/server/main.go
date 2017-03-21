package main

import (
	"github.com/liangdas/mqant"
	"github.com/liangdas/mqant/module"
	"server/chat"
	"server/gate"
	"server/login"
	"server/hitball"
)

//func ChatRoute( app module.App,Type string,hash string) (*module.ServerSession){
//	//演示多个服务路由 默认使用第一个Server
//	log.Debug("Hash:%s 将要调用 type : %s",hash,Type)
//	servers:=app.GetServersByType(Type)
//	if len(servers)==0{
//		return nil
//	}
//	return servers[0]
//}
func main() {
	app := mqant.CreateApp()
	//app.Route("Chat",ChatRoute)
	app.Run(true, //只有是在调试模式下才会在控制台打印日志, 非调试模式下只在日志文件中输出日志
		module.MasterModule(),
		hitball.Module(),
		gate.Module(),  //这是默认网关模块,是必须的支持 TCP,websocket,MQTT协议
		login.Module(), //这是用户登录验证模块
		chat.Module())  //这是聊天模块

}
