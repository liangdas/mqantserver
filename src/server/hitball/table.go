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
package hitball
import (
	"math/rand"
	"time"
	"encoding/json"
	"github.com/liangdas/mqant/gate"
	"github.com/yireyun/go-queue"
	"fmt"
	"github.com/liangdas/mqant/module"
	"server/hitball/objects"
)
var(
	friction float64=0.99 	 // friction affects ball speed 速度递减因子
	degToRad float64=0.0174532925 // degrees-radians conversion
	minPower float64= 50 // minimum power applied to ball
	maxPower float64=200 // maximum power applied to ball
	rotateSpeed = 3; // arrow rotation speed
	ballRadius float64=20;
)

type CallBackMsg struct {
	notify		bool	//是否是广播
	players 	[]string	//如果不是广播就指定session
	topic 		*string
	body 		*[]byte
}

func (self *Table)sendCallBackMsg(players []string,topic string,body []byte) error {
	ok, quantity := self.queue_callback.Put(&CallBackMsg{
		notify:false,
		players:players,
		topic:&topic,
		body:&body,
	})
	if !ok {
		return	fmt.Errorf("Put Fail, quantity:%v\n", quantity)
	} else {
		return nil
	}
}

func (self *Table)notifyCallBackMsg(topic string,body []byte) error {
	ok, quantity := self.queue_callback.Put(&CallBackMsg{
		notify:true,
		players:nil,
		topic:&topic,
		body:&body,
	})
	if !ok {
		return	fmt.Errorf("Put Fail, quantity:%v\n", quantity)
	} else {
		return nil
	}
}

/**
【每帧调用】统一发送所有消息给各个客户端
 */
func (self *Table)executeCallBackMsg()  {
	ok := true
	queue:=self.queue_callback
	index:=0
	for ok {
		val, _ok, _ := queue.Get()
		index++
		if _ok{
			msg:=val.(*CallBackMsg)
			if msg.notify{
				for _,role:=range self.player{
					e := role.Session.Send(*msg.topic, *msg.body)
					if e != "" {
						self.remove(role.Rid)
					}
				}
			}else{
				for _,rid:=range msg.players{
					if player,ok:=self.player[rid];ok{
						e := player.Session.Send(*msg.topic, *msg.body)
						if e != "" {
							self.remove(rid)
						}
					}
				}
			}
		}
		ok=_ok
	}
}

type Table struct {
	BaseTable
	module 		module.Module
	player		map[string]*objects.Player
	coins		map[int]*objects.Coins
	queue_callback		*queue.EsQueue
	tableId		int
	current_id	int
	current_frame	int		//当前帧
	sync_frame	int		//上一次同步数据的帧
	stoped		bool
	world_width	float64
	world_height	float64
	proTime		int64
}

func NewTable(module module.Module,tableId int)(*Table){
	table:=&Table{
		module:module,
		player:map[string]*objects.Player{},
		coins:map[int]*objects.Coins{},
		queue_callback:queue.NewQueue(256),
		stoped:true,
		current_id:0,
		tableId:tableId,
		current_frame:0,
		sync_frame:0,
		world_width:1280,
		world_height:1280,
	}
	table.Init()
	table.Register("Join",table.join)
	table.Register("EatCoins",table.eatCoins)
	table.Register("Remove",table.remove)
	table.Register("Fire",table.fire)
	return table
}

func (self *Table)Empty() bool{
	return self.stoped
}

func (self *Table)Full() bool{
	return false
}

func (self *Table)TableId()int{
	return self.tableId
}

func (self *Table)Start(){
	if self.stoped{
		self.stoped=false
		go func() {
			//这里设置为22ms但实际上每次循环大概是23-25ms左右，根据机器定,客户端设置的帧为40
			tick := time.NewTicker(22 * time.Millisecond)
			defer func() {
				tick.Stop()
			}()
			for !self.stoped{
				select {
				case <-tick.C:
					self.Update(nil)
				}
			}
		}()
	}
}


func (self *Table)Stop(){
	self.stoped=true
}

/**
定帧计算所有玩家的位置
 */
func (self *Table)Update(arge interface{}){
	self.current_frame++
	self.ExecuteEvent(arge)	//执行这一帧客户端发送过来的消息

	//位置计算
	//startTime := time.Now().UnixNano()
	for _,role:=range self.player{
		self.wallBounce(role)
		role.Move(friction)
		role.Rotate()
	}
	//s :=(time.Now().UnixNano()-self.proTime)/1000000
	//self.proTime = time.Now().UnixNano()
	//if s>26{
	//	log.Debug("exct time %d et %d ns",s ,(time.Now().UnixNano()-startTime)/1000)
	//}

	if self.current_frame-self.sync_frame>3{
		//每四帧同步一次
		self.sync_frame=self.current_frame
		self.NotifyAxes(nil)
	}

	if len(self.coins)<8{
		self.addCoins()
	}

	self.executeCallBackMsg()	//统一发送数据到客户端
}


/**
添加一个金币
 */
func (self *Table)addCoins(){
	randomX:=rand.Float64()*(self.world_width-2*ballRadius)+ballRadius
	randomY:=rand.Float64()*(self.world_height-2*ballRadius)+ballRadius
	self.current_id++
	coins:=&objects.Coins{
		X:	randomX,
		Y:	randomY,
		Id:	self.current_id,
		Wid:	0,
		Type:	0,
	}
	self.coins[coins.Id]=coins
	self.NotifyAddCoins(coins) //广播给所有其他玩家
}
/**
玩家吃了金币
 */
func (self *Table)eatCoins(session gate.Session,Id int){
	if coins,ok:=self.coins[Id];ok{
		delete(self.coins,Id)
		self.NotifyEatCoins(coins) //广播给所有其他玩家
	}
}
/**
玩家加入场景
 */
func (self *Table)join(session gate.Session) {
	if player,ok:=self.player[session.GetUserId()];ok{
		//这个玩家已经在游戏中了
		player.OnRequest(session)
		return
	}
	randomX:=rand.Float64()*(self.world_width-2*ballRadius)+ballRadius
	randomY:=rand.Float64()*(self.world_height-2*ballRadius)+ballRadius
	player:=&objects.Player{
		X:	randomX,
		Y:	randomY,
		Wid:	0,
		Rid:	session.GetUserId(),
		RotateDirection:1,
		XSpeed: 0,
		YSpeed: 0,
		BallRadius:ballRadius,
		Session:session,

		RotateSpeed:rotateSpeed,
		DegToRad	:degToRad,
		MinPower 	:minPower,
		MaxPower 	:maxPower,
	}


	self.player[player.Rid]=player
	self.NotifyJoin(player) //广播给所有其他玩家
	result:=map[string]interface{}{
		"Rid":player.Rid,
		"Player":self.player,
		"Coins":self.coins,
	}
	b, _ := json.Marshal(result)
	self.sendCallBackMsg([]string{player.Rid},"Hitball/OnEnter", b)
}

func (self *Table)remove(Rid string)(error) {
	delete(self.player,Rid)
	return nil
}
/**
玩家点击屏幕开始奔跑
 */
func (self *Table)fire(session gate.Session,X float64, Y float64,angle float64,power float64) {
	//发射
	if player,ok:=self.player[session.GetUserId()];ok{
		player.Fire(X,Y,angle,power)
		player.OnRequest(session)
	}
}

/**
定期刷新所有玩家的位置
 */
func (self *Table)NotifyAxes(arge interface{}){
	b, _ := json.Marshal(self.player)
	self.notifyCallBackMsg("Hitball/OnMove", b)
}

/**
通知所有玩家有新玩家加入
 */
func (self *Table)NotifyJoin(player *objects.Player){
	b, _ := json.Marshal(player)
	self.notifyCallBackMsg("Hitball/OnJoin", b)
}

/**
通知所有玩家新加了金币
 */
func (self *Table)NotifyAddCoins(coins *objects.Coins){
	b, _ := json.Marshal(coins)
	self.notifyCallBackMsg("Hitball/OnAddCoins", b)
}
/**
通知所有玩家金币已经被吃掉
 */
func (self *Table)NotifyEatCoins(coins *objects.Coins){
	b, _ := json.Marshal(coins)
	self.notifyCallBackMsg("Hitball/OnEatCoins", b)
}

/**
检查小球是否超出世界范围
 */
func (self *Table) wallBounce (player *objects.Player){
	if(player.X<player.BallRadius){
		player.X=player.BallRadius;
		player.XSpeed*=-1
	}
	if(player.Y<player.BallRadius){
		player.Y=player.BallRadius;
		player.YSpeed*=-1
	}

	if(player.X>self.world_width-player.BallRadius){
		player.X=self.world_width-player.BallRadius;
		player.XSpeed*=-1
	}
	if(player.Y>self.world_height-player.BallRadius){
		player.Y=self.world_height-player.BallRadius;
		player.YSpeed*=-1
	}
}