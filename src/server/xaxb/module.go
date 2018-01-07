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
// 多人猜数字游戏
// 游戏周期:
//第一阶段 空闲期，什么都不做
//第一阶段 空档期，等待押注
//第二阶段 押注期  可以押注
//第三阶段 开奖期  开奖
//第四阶段 结算期  结算
//玩法:
//玩家可以押 0-9 中的一个数字,每押一次需要消耗500金币
//押注完成后牌桌内随机开出一个 0-9 的数字
//哪个玩家猜中的数字与系统开出的数字最近就算赢,可以赢取所有玩家本局押注的金币*80%
//玩家金币用完后将被踢出房间
package xaxb

import (
	"errors"
	"fmt"
	"github.com/liangdas/mqant/conf"
	"github.com/liangdas/mqant/gate"
	"github.com/liangdas/mqant/module"
	"github.com/liangdas/mqant/module/base"
	"github.com/liangdas/mqant-modules/room"
)

var Module = func() module.Module {
	this := new(xaxb)
	return this
}

type xaxb struct {
	basemodule.BaseModule
	room    *room.Room
	proTime int64
	gameId  int
}

func (self *xaxb) GetType() string {
	//很关键,需要与配置文件中的Module配置对应
	return "XaXb"
}
func (self *xaxb) Version() string {
	//可以在监控时了解代码版本
	return "1.0.0"
}
func (self *xaxb) GetFullServerId() string {
	return self.GetType() + "@" + self.GetServerId()
}
func (self *xaxb) usableTable(table room.BaseTable) bool {
	return table.AllowJoin()
}
func (self *xaxb) newTable(module module.RPCModule, tableId int) (room.BaseTable, error) {
	table := NewTable(module, tableId)
	return table, nil
}
func (self *xaxb) OnInit(app module.App, settings *conf.ModuleSettings) {
	self.BaseModule.OnInit(self, app, settings)
	self.gameId = 13
	self.room = room.NewRoom(self, self.gameId, self.newTable, self.usableTable)
	self.GetServer().Register("GetUsableTable", self.getUsableTable)
	self.GetServer().Register("HD_GetUsableTable", self.HDGetUsableTable)
	self.GetServer().Register("HD_Enter", self.enter)
	self.GetServer().RegisterGO("HD_Exit", self.exit)
	self.GetServer().RegisterGO("HD_SitDown", self.sitdown)
	self.GetServer().RegisterGO("HD_StartGame", self.startGame)
	self.GetServer().RegisterGO("HD_PauseGame", self.pauseGame)
	self.GetServer().RegisterGO("HD_Stake", self.stake)
}

func (self *xaxb) Run(closeSig chan bool) {

}

func (self *xaxb) OnDestroy() {
	//一定别忘了关闭RPC
	self.GetServer().OnDestroy()
}

/**
检查参数是否存在
*/
func (self *xaxb) ParameterCheck(msg map[string]interface{}, paras ...string) error {
	for _, v := range paras {
		if _, ok := msg[v]; !ok {
			return fmt.Errorf("No %s found", v)
		}
	}
	return nil
}

/**
检查参数是否存在
*/
func (self *xaxb) GetTableByBigRoomId(bigRoomId string) (*Table, error) {
	_, tableid, _, err := room.ParseBigRoomId(bigRoomId)
	if err != nil {
		return nil, err
	}
	table := self.room.GetTable(tableid)
	if table != nil {
		tableimp := table.(*Table)
		return tableimp, nil
	} else {
		return nil, errors.New("No table found")
	}
}
/**
创建一个房间
*/
func (self *xaxb) HDGetUsableTable(session gate.Session,msg map[string]interface{}) (map[string]interface{}, string) {
	return self.getUsableTable(session)
}
/**
创建一个房间
*/
func (self *xaxb) getUsableTable(session gate.Session) (map[string]interface{}, string) {
	table, err := self.room.GetUsableTable()
	if err == nil {
		table.Create()
		tableInfo := map[string]interface{}{
			"BigRoomId": room.BuildBigRoomId(self.GetFullServerId(), table.TableId(), table.TransactionId()),
		}
		return tableInfo, ""
	} else {
		return nil, "There is no available table"
	}
}

func (self *xaxb) enter(session gate.Session, msg map[string]interface{}) (string, string) {
	if BigRoomId, ok := msg["BigRoomId"]; !ok {
		return "", "No BigRoomId found"
	} else {
		bigRoomId := BigRoomId.(string)

		moduleId, tableid, _, err := room.ParseBigRoomId(bigRoomId)
		if err != nil {
			return "", err.Error()
		}
		if session.Get("BigRoomId") != "" {
			//用户当前已经加入过一个BigRoomId
			if session.Get("BigRoomId") != bigRoomId {
				//先从上一个桌子退出
				_, e := self.RpcInvoke(moduleId, "HD_Exit", session, map[string]interface{}{
					"BigRoomId": session.Get("BigRoomId"),
				})
				if e != "" {
					return "", e
				}
			}
		}
		table := self.room.GetTable(tableid)
		if table != nil {
			tableimp := table.(*Table)
			if table.VerifyAccessAuthority(session.GetUserid(), bigRoomId) == false {
				return "", "Access rights validation failed"
			}
			erro := tableimp.Join(session)
			if erro == nil {
				bigRoomId = room.BuildBigRoomId(self.GetFullServerId(), table.TableId(), table.TransactionId())
				session.Set("BigRoomId", bigRoomId) //设置到session
				session.Push()
				return bigRoomId, ""
			}
			return "", erro.Error()
		} else {
			return "", "No room found"
		}
	}

}

func (self *xaxb) exit(session gate.Session, msg map[string]interface{}) (string, string) {
	if BigRoomId, ok := msg["BigRoomId"]; !ok {
		return "", "No BigRoomId found"
	} else {
		bigRoomId := BigRoomId.(string)
		table, err := self.GetTableByBigRoomId(bigRoomId)
		if err != nil {
			return "", err.Error()
		}
		err = table.Exit(session)
		if err == nil {
			bigRoomId = room.BuildBigRoomId(self.GetFullServerId(), table.TableId(), table.TransactionId())
			session.Set("BigRoomId", "") //设置到session
			session.Push()
			return bigRoomId, ""
		}
		return "", err.Error()
	}

}

func (self *xaxb) sitdown(session gate.Session, msg map[string]interface{}) (string, string) {
	bigRoomId := session.Get("BigRoomId")
	if bigRoomId == "" {
		return "", "fail"
	}
	table, err := self.GetTableByBigRoomId(bigRoomId)
	if err != nil {
		return "", err.Error()
	}
	err = table.PutQueue("SitDown", session)
	if err != nil {
		return "", err.Error()
	}
	return "success", ""
}
func (self *xaxb) startGame(session gate.Session, msg map[string]interface{}) (string, string) {
	bigRoomId := session.Get("BigRoomId")
	if bigRoomId == "" {
		return "", "fail"
	}
	table, err := self.GetTableByBigRoomId(bigRoomId)
	if err != nil {
		return "", err.Error()
	}
	err = table.PutQueue("StartGame", session)
	if err != nil {
		return "", err.Error()
	}
	return "success", ""
}
func (self *xaxb) pauseGame(session gate.Session, msg map[string]interface{}) (string, string) {
	bigRoomId := session.Get("BigRoomId")
	if bigRoomId == "" {
		return "", "fail"
	}
	table, err := self.GetTableByBigRoomId(bigRoomId)
	if err != nil {
		return "", err.Error()
	}
	err = table.PutQueue("PauseGame", session)
	if err != nil {
		return "", err.Error()
	}
	return "success", ""
}

func (self *xaxb) stake(session gate.Session, msg map[string]interface{}) (string, string) {
	if Target, ok := msg["Target"]; !ok {
		return "", "No Target found"
	} else {
		bigRoomId := session.Get("BigRoomId")
		if bigRoomId == "" {
			return "", "fail"
		}
		table, err := self.GetTableByBigRoomId(bigRoomId)
		if err != nil {
			return "", err.Error()
		}
		err = table.PutQueue("Stake", session, int64(Target.(float64)))
		if err != nil {
			return "", err.Error()
		}
		return "success", ""
	}
}
