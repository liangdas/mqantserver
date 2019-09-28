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
package objects

import (
	"encoding/json"
	"github.com/liangdas/mqant-modules/room"
)

type Player struct {
	room.BasePlayerImp
	SeatIndex  int
	Coin       int //金币数量
	timeToMove int64
	Target     int64 //押注目标
	Stake      bool  //是否已押注
	Weight     int64 //计算后权重
}

func NewPlayer(SeatIndex int) *Player {
	this := new(Player)
	this.SeatIndex = SeatIndex
	this.Coin = 1000
	return this
}

func (this *Player) Serializable() ([]byte, error) {

	return json.Marshal(this.SerializableMap())
}

func (this *Player) SerializableMap() map[string]interface{} {
	rid := ""
	if this.Session() != nil {
		rid = this.Session().GetUserId()
	}
	return map[string]interface{}{
		"SeatIndex": this.SeatIndex,
		"Rid":       rid,
		"Coin":      this.Coin,
		"Stake":     this.Stake,
		"Target":    this.Target,
		"Weight":    this.Weight,
		"SitDown":   this.SitDown(),
	}
}
