/**
一定要记得在confin.json配置这个模块的参数,否则无法使用
 */
package chat
import (
	"github.com/liangdas/mqant/module"
	"encoding/json"
	"github.com/liangdas/mqant/gate"
	"github.com/liangdas/mqant/conf"
	"github.com/liangdas/mqant/log"
)
var Module = func() (module.Module){
	chat := new(Chat)
	return chat
}


type Chat struct {
	app	module.App
	server *module.Server
	chats  map[string]map[string]*gate.Session
}
func (m *Chat) GetType()(string){
	//很关键,需要与配置文件中的Module配置对应
	return "Chat"
}
func (m *Chat) GetServer() (*module.Server){
	if m.server==nil{
		m.server = new(module.Server)
	}
	return m.server
}

func (m *Chat) OnInit(app module.App,settings *conf.ModuleSettings) {
	//初始化模块
	m.app=app
	m.chats=map[string]map[string]*gate.Session{}

	//创建一个远程调用的RPC
	m.GetServer().OnInit(app,settings)
	//注册远程调用的函数
	m.GetServer().RegisterGO("HD_JoinChat",m.joinChat) //我们约定所有对客户端的请求都以Handler_开头
	m.GetServer().RegisterGO("HD_Say",m.say) //我们约定所有对客户端的请求都以Handler_开头

}

func (m *Chat) Run(closeSig chan bool) {
	//运行模块
}

func (m *Chat) OnDestroy() {
	//注销模块
	//一定别忘了关闭RPC
	m.GetServer().OnDestroy()
}

func (m *Chat) joinChat(s map[string]interface{},msg map[string]interface{})(result map[string]interface{},err string) {
	if msg["roomName"]==""{
		err="roomName cannot be nil"
		return
	}
	session:=gate.NewSession(m.app,s)
	log.Debug("session %v",session.ExportMap())
	if session.Userid==""{
		err="Not Logined"
		return
	}
	roomName:=msg["roomName"].(string)

	r,_:=m.app.RpcInvoke("Login","getRand",roomName)

	log.Debug("演示模块间RPC调用 :",r)

	userList:=m.chats[roomName]
	if userList==nil{
		//添加一个新的房间
		userList=map[string]*gate.Session{session.Userid:session}
		m.chats[roomName]=userList
	}else{
		//user:=userList[session.Userid]
		//if user!=nil{
			//已经加入过这个聊天室了 不过这里还是替换一下session 因此用户可能是重连的
			//err="Already in this chat room"
			//userList[session.Userid]=session
			//return
		//}
		//添加这个用户进入聊天室
		userList[session.Userid]=session
	}

	rmsg:=map[string]string{}
	rmsg["roomName"]=roomName
	rmsg["user"]=session.Userid
	b,_:=json.Marshal(rmsg)

	userL:=make([]string,len(userList))
	//广播添加用户信息到该房间的所有用户
	i:=0
	for _,user:=range userList{
		if user.Userid!=session.Userid{
			//给其他用户发送消息
			err:=user.Send("Chat/OnJoin",b)
			if err!=""{
				//信息没有发送成功
				m.onLeave(roomName,user.Userid)
			}
		}
		userL[i]=user.Userid
		i++

	}
	result=map[string]interface{}{
		"users":userL,
	}
	return
}

func (m *Chat) say(s map[string]interface{},msg map[string]interface{})(result string,err string){
	if msg["roomName"]==nil||msg["content"]==nil{
		err="roomName or say cannot be nil"
		return
	}
	session:=gate.NewSession(m.app,s)
	if session.Userid==""{
		err="Not Logined"
		return
	}
	roomName:=msg["roomName"].(string)
	//from:=msg["from"].(string)
	target:=msg["target"].(string)
	content:=msg["content"].(string)
	userList:=m.chats[roomName]
	if userList==nil{
		err="No room"
		return
	}else{
		user:=userList[session.Userid]
		if user==nil{
			err="You haven't been in the room yet"
			return
		}
		rmsg:=map[string]string{}
		rmsg["roomName"]=roomName
		rmsg["from"]=session.Userid
		rmsg["target"]=target
		rmsg["msg"]=content
		b,_:=json.Marshal(rmsg)
		if target=="*"{
			//广播添加用户信息到该房间的所有用户
			for _,user:=range userList{
				err:=user.Send("Chat/OnChat",b)
				if err!=""{
					//信息没有发送成功
					m.onLeave(roomName,user.Userid)
				}
			}
		}else{
			user:=userList[target]
			if user==nil{
				err="This user haven't been in the room yet"
				return
			}
			e:=user.Send("Chat/OnChat",b)
			if e!=""{
				//信息没有发送成功
				m.onLeave(roomName,user.Userid)
				err="The user has left the room"
				return
			}
		}


	}
	result="say success"
	return
}
/**
用户 断开连接 广播离线消息
 */
func (m *Chat) onLeave(roomName string,Userid string){
	userList:=m.chats[roomName]
	if userList==nil{
		return
	}
	delete(userList,Userid) //从列表中删除
	rmsg:=map[string]string{}
	rmsg["roomName"]=roomName
	rmsg["user"]=Userid
	b,_:=json.Marshal(rmsg)
	for _,user:=range userList{
		user.SendNR("Chat/OnLeave",b)
	}
}

