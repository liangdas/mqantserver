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
package xaxb

import (
	"encoding/json"
	"mqantserver/src/server/xaxb/objects"
)

/**
定期刷新所有玩家的位置
*/
func (self *Table) NotifyAxes() {
	seats := []map[string]interface{}{}
	for _, player := range self.seats {
		if player.Bind() {
			seats = append(seats, player.SerializableMap())
		}
	}
	b, _ := json.Marshal(map[string]interface{}{
		"State":     self.State(),
		"StateGame": self.fsm.getState(),
		"Seats":     seats,
	})
	self.NotifyCallBackMsg("XaXb/OnSync", b)
}

/**
通知所有玩家有新玩家加入
*/
func (self *Table) NotifyJoin(player *objects.Player) {
	b, _ := json.Marshal(player.SerializableMap())
	self.NotifyCallBackMsg("XaXb/OnEnter", b)
}

/**
通知所有玩家开始游戏了
*/
func (self *Table) NotifyResume() {
	b, _ := json.Marshal(self.getSeatsMap())
	self.NotifyCallBackMsg("XaXb/OnResume", b)
}

/**
通知所有玩家开始游戏了
*/
func (self *Table) NotifyPause() {
	b, _ := json.Marshal(self.getSeatsMap())
	self.NotifyCallBackMsg("XaXb/OnPause", b)
}

/**
通知所有玩家开始游戏了
*/
func (self *Table) NotifyStop() {
	b, _ := json.Marshal(self.getSeatsMap())
	self.NotifyCallBackMsg("XaXb/OnStop", b)
}

/**
通知所有玩家进入空闲期了
*/
func (self *Table) NotifyIdle() {
	b, _ := json.Marshal(map[string]interface{}{
		"Coin": 500,
	})
	self.NotifyCallBackMsg("XaXb/Idle", b)
}

/**
通知所有玩家开始押注了
*/
func (self *Table) NotifyBetting() {
	b, _ := json.Marshal(map[string]interface{}{
		"Coin": 500,
	})
	self.NotifyCallBackMsg("XaXb/Betting", b)
}

/**
通知所有玩家开始开奖了
*/
func (self *Table) NotifyOpening() {
	b, _ := json.Marshal(map[string]interface{}{
		"Coin": 500,
	})
	self.NotifyCallBackMsg("XaXb/Opening", b)
}

/**
通知所有玩家开奖结果出来了
*/
func (self *Table) NotifySettlement(Result int64) {
	seats := []map[string]interface{}{}
	for _, player := range self.seats {
		if player.Bind() {
			seats = append(seats, player.SerializableMap())
		}
	}
	b, _ := json.Marshal(map[string]interface{}{
		"Result": Result,
		"Seats":  seats,
	})
	self.NotifyCallBackMsg("XaXb/Settlement", b)
}
