// Copyright 2014 mqantserver Author. All Rights Reserved.
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
package chat

import (
	"github.com/liangdas/mqant/rpc"
	"github.com/liangdas/mqant/log"
)

type Listener struct {
	server		*mqrpc.RPCServer
	moduleType 	string
	serverId	string
}

func (l *Listener)OnTimeOut(fn string ,Expired int64) {
	log.Error("[%s]的请求(%s)超时了!",l.moduleType,fn)
}
func (l *Listener)OnError(fn string,params []interface{},err error) {
	log.Error("[%s]的请求(%s)出现异常 error(%s)!",l.moduleType,fn,err.Error())
}
/**
fn 		方法名
params		参数
result		执行结果
exec_time 	方法执行时间 单位为 Nano 纳秒  1000000纳秒等于1毫秒
 */
func (l *Listener)OnComplete(fn string,params []interface{},result *mqrpc.ResultInfo,exec_time int64) {
	log.Error("[%s]的请求(%s) 执行时间为:[%d 微妙] 当前正在执行的请求数:[%d]!",l.moduleType,fn,exec_time/1000,l.server.GetExecuting())
}


