package main
import (
	"github.com/liangdas/mqant/conf"
	"github.com/liangdas/mqant"
	"github.com/liangdas/mqant/log"
	"server/chat"
	"server/login"
	"server/gate"
	"os"
	"path/filepath"
	"os/exec"
	"fmt"

	"flag"
)
//func ChatRoute( app module.App,moduleType string,serverId string,Type string) (*module.ServerSession){
//	//演示多个服务路由 默认使用第一个Server
//	log.Debug("Type:%s Id:%s 将要调用 type : %s",moduleType,serverId,Type)
//	servers:=app.GetServersByType(Type)
//	if len(servers)==0{
//		return nil
//	}
//	return servers[0]
//}
func main() {
	file, _ := exec.LookPath(os.Args[0])
	ApplicationPath, _ := filepath.Abs(file)
	ApplicationDir, _ := filepath.Split(ApplicationPath)
	defaultPath:= fmt.Sprintf("%s/conf/server.conf",ApplicationDir)
	confPath := flag.String("conf", defaultPath, "Server configuration file path")
	flag.Parse() //解析输入的参数
	log.Release("Server configuration file path [%s]",*confPath)
	f, err := os.Open(*confPath)
	if err!=nil{
		panic(err)
	}
	conf.LoadConfig(f.Name()) //加载配置文件
	app:=mqant.CreateApp()
	app.Configure(conf.Conf)  //配置信息
	//app.Route("Chat",ChatRoute)
	app.Run(gate.Module(),	//这是默认网关模块,是必须的支持 TCP,websocket,MQTT协议
		login.Module(), //这是用户登录验证模块
		chat.Module())  //这是聊天模块

}
