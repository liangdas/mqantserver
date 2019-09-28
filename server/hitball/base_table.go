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
package hitball

import (
	"fmt"
	"github.com/liangdas/mqant/log"
	"github.com/yireyun/go-queue"
	"reflect"
	"runtime"
	"sync"
)

type QueueMsg struct {
	Func   string
	Params []interface{}
}
type QueueReceive interface {
	Receive(msg *QueueMsg, index int)
}
type BaseTable struct {
	functions       map[string]interface{}
	receive         QueueReceive
	queue0          *queue.EsQueue
	queue1          *queue.EsQueue
	current_w_queue int //当前写的队列
	lock            *sync.RWMutex
}

func (self *BaseTable) Init() {
	self.functions = map[string]interface{}{}
	self.queue0 = queue.NewQueue(256)
	self.queue1 = queue.NewQueue(256)
	self.current_w_queue = 0
	self.lock = new(sync.RWMutex)
}
func (self *BaseTable) SetReceive(receive QueueReceive) {
	self.receive = receive
}
func (self *BaseTable) Register(id string, f interface{}) {

	if _, ok := self.functions[id]; ok {
		panic(fmt.Sprintf("function id %v: already registered", id))
	}

	self.functions[id] = f
}

/**
协成安全,任意协成可调用
*/
func (self *BaseTable) PutQueue(_func string, params ...interface{}) error {
	ok, quantity := self.wqueue().Put(&QueueMsg{
		Func:   _func,
		Params: params,
	})
	if !ok {
		return fmt.Errorf("Put Fail, quantity:%v\n", quantity)
	} else {
		return nil
	}

}

/**
切换并且返回读的队列
*/
func (self *BaseTable) switchqueue() *queue.EsQueue {
	self.lock.Lock()
	if self.current_w_queue == 0 {
		self.current_w_queue = 1
		self.lock.Unlock()
		return self.queue0
	} else {
		self.current_w_queue = 0
		self.lock.Unlock()
		return self.queue1
	}

}
func (self *BaseTable) wqueue() *queue.EsQueue {
	self.lock.Lock()
	if self.current_w_queue == 0 {
		self.lock.Unlock()
		return self.queue0
	} else {
		self.lock.Unlock()
		return self.queue1
	}

}

/**
【每帧调用】执行队列中的所有事件
*/
func (self *BaseTable) ExecuteEvent(arge interface{}) {
	ok := true
	queue := self.switchqueue()
	index := 0
	for ok {
		val, _ok, _ := queue.Get()
		index++
		if _ok {
			if self.receive != nil {
				self.receive.Receive(val.(*QueueMsg), index)
			} else {
				msg := val.(*QueueMsg)
				function, ok := self.functions[msg.Func]
				if !ok {
					fmt.Println(fmt.Sprintf("Remote function(%s) not found", msg.Func))
					continue
				}
				f := reflect.ValueOf(function)
				in := make([]reflect.Value, len(msg.Params))
				for k, _ := range in {
					in[k] = reflect.ValueOf(msg.Params[k])
				}
				_runFunc := func() {
					defer func() {
						if r := recover(); r != nil {
							var rn = ""
							switch r.(type) {

							case string:
								rn = r.(string)
							case error:
								rn = r.(error).Error()
							}
							buf := make([]byte, 1024)
							l := runtime.Stack(buf, false)
							errstr := string(buf[:l])
							log.Error("table qeueu event(%s) exec fail error:%s \n ----Stack----\n %s", msg.Func, rn, errstr)
						}
					}()
					f.Call(in)
				}
				_runFunc()
			}
		}
		ok = _ok
	}
}
