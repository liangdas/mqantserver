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
	"fmt"
	"github.com/liangdas/mqant/utils"
	"math"
	"server/xaxb/objects"
)

var (
	VoidPeriod            = FSMState("空档期")
	IdlePeriod            = FSMState("空闲期")
	BettingPeriod         = FSMState("押注期")
	OpeningPeriod         = FSMState("开奖期")
	SettlementPeriod      = FSMState("结算期")
	VoidPeriodEvent       = FSMEvent("进入空档期")
	IdlePeriodEvent       = FSMEvent("进入空闲期")
	BettingPeriodEvent    = FSMEvent("进入押注期")
	OpeningPeriodEvent    = FSMEvent("进入开奖期")
	SettlementPeriodEvent = FSMEvent("进入结算期")
)

func (this *Table) InitFsm() {
	this.fsm = *NewFSM(VoidPeriod)
	this.VoidPeriodHandler = FSMHandler(func() FSMState {
		fmt.Println("已进入空档期")
		return VoidPeriod
	})
	this.IdlePeriodHandler = FSMHandler(func() FSMState {
		fmt.Println("已进入空闲期")
		this.step1 = this.current_frame
		this.NotifyIdle()

		for _, seat := range this.GetSeats() {
			player := seat.(*objects.Player)
			if player.Bind() {
				if player.Coin <= 0 {
					player.Session().Send("XaXb/Exit", []byte(`{"Info":"金币不足你被强制离开房间"}`))
					player.OnUnBind() //踢下线
				}
			}
		}

		return IdlePeriod
	})
	this.BettingPeriodHandler = FSMHandler(func() FSMState {
		fmt.Println("已进入押注期")
		this.step2 = this.current_frame
		this.NotifyBetting()
		return BettingPeriod
	})
	this.OpeningPeriodHandler = FSMHandler(func() FSMState {
		fmt.Println("已进入开奖期")
		this.step3 = this.current_frame
		this.NotifyOpening()
		return OpeningPeriod
	})
	this.SettlementPeriodHandler = FSMHandler(func() FSMState {
		fmt.Println("已进入结算期")
		var mixWeight int64 = math.MaxInt64
		var winer *objects.Player = nil
		Result := utils.RandInt64(0, 10)
		for _, seat := range this.GetSeats() {
			player := seat.(*objects.Player)
			if player.Stake {
				player.Weight = int64(math.Abs(float64(player.Target - Result)))
				if mixWeight > player.Weight {
					mixWeight = player.Weight
					winer = player
				}
			}
		}
		if winer != nil {
			winer.Coin += 800
		}

		this.step4 = this.current_frame
		this.NotifySettlement(Result)
		return SettlementPeriod
	})

	this.fsm.AddHandler(IdlePeriod, VoidPeriodEvent, this.VoidPeriodHandler)
	this.fsm.AddHandler(SettlementPeriod, VoidPeriodEvent, this.VoidPeriodHandler)
	this.fsm.AddHandler(BettingPeriod, VoidPeriodEvent, this.VoidPeriodHandler)
	this.fsm.AddHandler(OpeningPeriod, VoidPeriodEvent, this.VoidPeriodHandler)

	this.fsm.AddHandler(VoidPeriod, IdlePeriodEvent, this.IdlePeriodHandler)
	this.fsm.AddHandler(SettlementPeriod, IdlePeriodEvent, this.IdlePeriodHandler)

	this.fsm.AddHandler(IdlePeriod, BettingPeriodEvent, this.BettingPeriodHandler)
	this.fsm.AddHandler(BettingPeriod, OpeningPeriodEvent, this.OpeningPeriodHandler)
	this.fsm.AddHandler(OpeningPeriod, SettlementPeriodEvent, this.SettlementPeriodHandler)
}

/**
进入空闲期
*/
func (this *Table) StateSwitch() {
	switch this.fsm.getState() {
	case VoidPeriod:

	case IdlePeriod:
		if (this.current_frame - this.step1) > 5 {
			this.fsm.Call(BettingPeriodEvent)
		} else {
			//this.NotifyAxes()
		}
	case BettingPeriod:
		if (this.current_frame - this.step2) > 20 {
			this.fsm.Call(OpeningPeriodEvent)
		} else {
			ready := true
			for _, seat := range this.GetSeats() {
				player := seat.(*objects.Player)
				if player.SitDown() && !player.Stake {
					ready = false
				}
			}
			if ready {
				//都押注了直接开奖
				this.fsm.Call(OpeningPeriodEvent)
			}
		}
	case OpeningPeriod:
		if (this.current_frame - this.step3) > 5 {
			this.fsm.Call(SettlementPeriodEvent)
		} else {
			//this.NotifyAxes()
		}
	case SettlementPeriod:
		if (this.current_frame - this.step4) > 5 {
			this.fsm.Call(IdlePeriodEvent)
		} else {
			//this.NotifyAxes()
		}
	}
}
