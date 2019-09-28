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
	"fmt"
	"github.com/liangdas/mqant/gate"
	"github.com/liangdas/mqant/log"
	"github.com/liangdas/mqant/module"
	"github.com/liangdas/mqant/rpc"
	"github.com/liangdas/mqant/rpc/pb"
	"github.com/liangdas/mqant/rpc/util"
	"github.com/pkg/errors"
)

type Listener struct {
	module module.RPCModule
}

func (l *Listener) BeforeHandle(fn string, callInfo *mqrpc.CallInfo) error {
	//放行
	for i, Type := range callInfo.RpcInfo.ArgsType {
		v, err := argsutil.Bytes2Args(l.module.GetApp(), Type, callInfo.RpcInfo.Args[i])
		if err != nil {
			log.Error("BeforeHandle %v", err)
			continue
		}
		switch v2 := v.(type) { //多选语句switch
		case gate.Session:
			//尝试加载Span
			if v2 != nil {
				if v2 == nil {
					return fmt.Errorf("session 不能为nil")
				}
				if v2.GetUserId() == "" {
					return fmt.Errorf("必须先登录账号")
				}
			}
		}
	}

	return nil
}
func (l *Listener) NoFoundFunction(fn string) (*mqrpc.FunctionInfo, error) {
	return nil, errors.Errorf("Remote function(%s) not found", fn)
}
func (l *Listener) OnTimeOut(fn string, Expired int64) {
	log.Error("请求(%s)超时了!", fn)
}
func (l *Listener) OnError(fn string, callInfo *mqrpc.CallInfo, err error) {
	log.Error("请求(%s)出现异常 error(%s)!", fn, err.Error())
}

/**
fn 		方法名
params		参数
result		执行结果
exec_time 	方法执行时间 单位为 Nano 纳秒  1000000纳秒等于1毫秒
*/
func (l *Listener) OnComplete(fn string, callInfo *mqrpc.CallInfo, result *rpcpb.ResultInfo, exec_time int64) {
	log.Info("请求(%s) 执行时间为:[%d 微妙]!", fn, exec_time/1000)
}
