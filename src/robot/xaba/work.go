// Copyright 2014 hey Author. All Rights Reserved.
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
package xaba_task

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/liangdas/armyant/task"
	"github.com/liangdas/armyant/work"
	"github.com/liangdas/mqant/utils"
	"io/ioutil"
	"time"
)

func NewWork(manager *Manager) *Work {
	this := new(Work)
	this.manager = manager
	//opts:=this.GetDefaultOptions("tls://127.0.0.1:3563")
	//opts := this.GetDefaultOptions("tcp://127.0.0.1:3563")
	opts := this.GetDefaultOptions("ws://127.0.0.1:3653")
	opts.SetConnectionLostHandler(func(client MQTT.Client, err error) {
		fmt.Println("ConnectionLost", err.Error())
	})
	opts.SetOnConnectHandler(func(client MQTT.Client) {
		fmt.Println("OnConnectHandler")
	})
	// load root ca
	// 需要一个证书，这里使用的这个网站提供的证书https://curl.haxx.se/docs/caextract.html
	caData, err := ioutil.ReadFile("/work/go/gopath/src/github.com/liangdas/armyant/mqtt_task/caextract.pem")
	if err != nil {
		fmt.Println(err.Error())
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caData)

	config := &tls.Config{
		RootCAs:            pool,
		InsecureSkipVerify: true,
	}
	opts.SetTLSConfig(config)
	err = this.Connect(opts)
	if err != nil {
		fmt.Println(err.Error())
	}

	this.On("XaXb/OnEnter", func(client MQTT.Client, msg MQTT.Message) {
		//服务端主动下发玩家加入事件
		fmt.Println(msg.Topic(), string(msg.Payload()))
	})
	this.On("XaXb/Exit", func(client MQTT.Client, msg MQTT.Message) {
		fmt.Println(msg.Topic(), string(msg.Payload()))
		this.GetClient().Disconnect(250)
	})
	this.On("XaXb/OnSync", func(client MQTT.Client, msg MQTT.Message) {
		fmt.Println(msg.Topic(), string(msg.Payload()))
	})
	this.On("XaXb/Idle", func(client MQTT.Client, msg MQTT.Message) {
		fmt.Println(msg.Topic(), string(msg.Payload()))
	})
	this.On("XaXb/Betting", func(client MQTT.Client, msg MQTT.Message) {
		//服务端通知可以押注了
		fmt.Println(msg.Topic(), string(msg.Payload()))
		time.Sleep(time.Millisecond * time.Duration(utils.RandInt64(100, 1000)))
		//开始押注
		msg, err = this.Request("XaXb/HD_Stake", []byte(fmt.Sprintf(`{"Target":%d}`, utils.RandInt64(0, 10))))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		//押注完成
		fmt.Println(msg.Topic(), string(msg.Payload()))
	})
	this.On("XaXb/Opening", func(client MQTT.Client, msg MQTT.Message) {
		fmt.Println(msg.Topic(), string(msg.Payload()))
	})
	this.On("XaXb/Settlement", func(client MQTT.Client, msg MQTT.Message) {
		fmt.Println(msg.Topic(), string(msg.Payload()))
	})
	return this
}

/**
Work 代表一个协程内具体执行任务工作者
*/
type Work struct {
	work.MqttWork
	manager *Manager
}

func (this *Work) UnmarshalResult(payload []byte) map[string]interface{} {
	rmsg := map[string]interface{}{}
	json.Unmarshal(payload, &rmsg)
	return rmsg["Result"].(map[string]interface{})
}

/**
每一次请求都会调用该函数,在该函数内实现具体请求操作

task:=task.Task{
		N:1000,	//一共请求次数，会被平均分配给每一个并发协程
		C:100,		//并发数
		//QPS:10,		//每一个并发平均每秒请求次数(限流) 不填代表不限流
}

N/C 可计算出每一个Work(协程) RunWorker将要调用的次数
*/
func (this *Work) RunWorker(t task.Task) {
	//登陆
	//s := `{"phone":"1880000000", "password":"123456"}`
	//msg,err:=this.Request("User/HD_LoginWithPassword",[]byte(s))
	//if err!=nil{
	//	return
	//}
	//fmt.Println(msg.Topic(),string(msg.Payload()))
	//申请牌桌
	msg, err := this.Request("XaXb/HD_GetUsableTable", []byte(`{"gameName":"xaxb"}`))
	if err != nil {
		return
	}
	fmt.Println(msg.Topic(), string(msg.Payload()))
	//进入牌桌
	BigRoomId := this.UnmarshalResult(msg.Payload())["BigRoomId"].(string)
	msg, err = this.Request("XaXb/HD_Enter", []byte(fmt.Sprintf(`{"BigRoomId":"%s"}`, BigRoomId)))
	if err != nil {
		return
	}
	fmt.Println(msg.Topic(), string(msg.Payload()))
	//坐下
	msg, err = this.Request("XaXb/HD_SitDown", []byte(fmt.Sprintf(`{"BigRoomId":"%s"}`, BigRoomId)))
	if err != nil {
		return
	}
	fmt.Println(msg.Topic(), string(msg.Payload()))

}
func (this *Work) Init(t task.Task) {

}
func (this *Work) Close(t task.Task) {

}