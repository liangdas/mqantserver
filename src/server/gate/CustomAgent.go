// Copyright 2014 loolgame Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//与客户端通信的自定义粘包示例，需要mqant v1.6.4版本以上才能运行
//该示例只用于简单的演示，并没有实现具体的粘包协议
package mgate

import (
	"github.com/liangdas/mqant/gate"
	"github.com/liangdas/mqant/module"
	"github.com/liangdas/mqant/network"
	"bufio"
	"github.com/liangdas/mqant/log"
	"time"
)

func NewAgent(module module.RPCModule)*CustomAgent{
	a := &CustomAgent{
		module:module,
	}
	return a
}

type CustomAgent struct {
	gate.Agent
	module 		module.RPCModule
	session                          gate.Session
	conn                             network.Conn
	r                                *bufio.Reader
	w                                *bufio.Writer
	gate                             gate.Gate
	rev_num                          int64
	send_num                         int64
	last_storage_heartbeat_data_time int64 //上一次发送存储心跳时间
	isclose                          bool
}
func (this *CustomAgent) OnInit(gate gate.Gate,conn network.Conn)error{
	log.Info("CustomAgent","OnInit")
	this.conn=conn
	this.gate=gate
	this.r=bufio.NewReader(conn)
	this.w=bufio.NewWriter(conn)
	this.isclose=false
	this.rev_num=0
	this.send_num=0
	return nil
}
/**
给客户端发送消息
 */
func (this *CustomAgent) WriteMsg(topic string, body []byte) error{
	this.send_num++
	//粘包完成后调下面的语句发送数据
	//this.w.Write()
	return nil
}

func (this *CustomAgent)Run() (err error){
	log.Info("CustomAgent","开始读数据了")

	this.session, err = this.gate.NewSessionByMap( map[string]interface{}{
		"Sessionid": "生成一个随机数",
		"Network":   this.conn.RemoteAddr().Network(),
		"IP":        this.conn.RemoteAddr().String(),
		"Serverid":  this.module.GetServerId(),
		"Settings":  make(map[string]string),
	})


	//这里可以循环读取客户端的数据


	//这个函数返回后连接就会被关闭
	return nil
}
/**
接收到一个数据包
 */
func (this *CustomAgent) OnRecover(topic string,msg []byte) {
	//通过解析的数据得到
	moduleType:=""
	_func:=""

	//如果要对这个请求进行分布式跟踪调试,就执行下面这行语句
	//a.session.CreateRootSpan("gate")

	//然后请求后端模块，第一个参数为session
	result, e :=this.module.RpcInvoke(moduleType , _func , this.session,msg)
	log.Info("result",result)
	log.Info("error",e )

	//回复客户端
	this.WriteMsg(topic,[]byte("请求成功了谢谢"))

	this.heartbeat()
}

func (this *CustomAgent)heartbeat(){
	//自定义网关需要你自己设计心跳协议
	if this.GetSession().GetUserid() != "" {
		//这个链接已经绑定Userid
		interval := time.Now().UnixNano()/1000000/1000 - this.last_storage_heartbeat_data_time //单位秒
		if interval > this.gate.GetMinStorageHeartbeat() {
			//如果用户信息存储心跳包的时长已经大于一秒
			if this.gate.GetStorageHandler() != nil {
				this.gate.GetStorageHandler().Heartbeat(this.GetSession().GetUserid())
				this.last_storage_heartbeat_data_time = time.Now().UnixNano() / 1000000 / 1000
			}
		}
	}
}

func (this *CustomAgent)Close(){
	log.Info("CustomAgent","主动断开连接")
	this.conn.Close()
}
func (this *CustomAgent)OnClose() error{
	this.isclose = true
	log.Info("CustomAgent","连接断开事件")
	//这个一定要调用，不然gate可能注销不了,造成内存溢出
	this.gate.GetAgentLearner().DisConnect(this) //发送连接断开的事件
	return nil
}
func (this *CustomAgent)Destroy(){
	this.conn.Destroy()
}
func (this *CustomAgent)RevNum() int64{
	return this.rev_num
}
func (this *CustomAgent)SendNum() int64{
	return this.send_num
}
func (this *CustomAgent)IsClosed() bool{
	return this.isclose
}
func (this *CustomAgent)GetSession() gate.Session{
	return this.session
}