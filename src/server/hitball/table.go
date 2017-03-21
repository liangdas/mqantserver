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
	"math"
	"time"
	"encoding/json"
	"github.com/liangdas/mqant/gate"
	"github.com/liangdas/mqant/log"
)
var(
	friction float64=0.99 	 // friction affects ball speed 速度递减因子
	degToRad float64=0.0174532925 // degrees-radians conversion
	minPower float64= 50 // minimum power applied to ball
	maxPower float64=200 // maximum power applied to ball
	rotateSpeed = 3; // arrow rotation speed
	ballRadius float64=20;
)
type player struct {
	session 	*gate.Session
	X		float64
	Y       	float64
	Wid		int
	Rid             string
	XSpeed 		float64
	YSpeed 		float64
	RotateDirection int     // rotate direction: 1-clockwise, 2-counterclockwise
	ballRadius	float64
	Angle		float64
	Power		float64
}

type coins struct {
	Id             	int
	X		float64
	Y       	float64
	Wid            	int
	Type		int
}

func (self *player)Move(){
	self.X=self.X+self.XSpeed;
	self.Y=self.Y+self.YSpeed;
	// reduce ball speed using friction 速度递减
	self.XSpeed*=friction;
	self.YSpeed*=friction;
}

func (self *player)Rotate () {
	self.Angle+=float64(rotateSpeed*self.RotateDirection);
}
func (self *player)Fire(X float64, Y float64,angle float64,power float64) {
	//发射
	self.XSpeed += math.Cos(angle*degToRad)*power/20
	self.YSpeed += math.Sin(angle*degToRad)*power/20
	self.Power = minPower
	//self.Angle=angle  //这里不同步客户端发过来的角速度
	self.X = X
	self.Y = Y
	self.Power=power
	self.RotateDirection*=-1;
}

type table struct {
	player		map[string]*player
	coins		map[int]*coins
	current_id	int
	stoped		bool
	world_width	float64
	world_height	float64
	proTime		int64
}

func NewTable()(*table){
	table:=&table{
		player:map[string]*player{},
		coins:map[int]*coins{},
		stoped:true,
		current_id:0,
		world_width:1280,
		world_height:1280,
	}
	return table
}
/**
添加一个金币
 */
func (self *table)AddCoins(){
	randomX:=rand.Float64()*(self.world_width-2*ballRadius)+ballRadius
	randomY:=rand.Float64()*(self.world_height-2*ballRadius)+ballRadius
	self.current_id++
	coins:=&coins{
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
func (self *table)EatCoins(Rid string,Id int){
	if coins,ok:=self.coins[Id];ok{
		delete(self.coins,Id)
		self.NotifyEatCoins(coins) //广播给所有其他玩家
	}
}
/**
玩家加入场景
 */
func (self *table)Join(Rid string,session *gate.Session)(map[string]interface{}) {
	randomX:=rand.Float64()*(self.world_width-2*ballRadius)+ballRadius
	randomY:=rand.Float64()*(self.world_height-2*ballRadius)+ballRadius
	player:=&player{
		X:	randomX,
		Y:	randomY,
		Wid:	0,
		Rid:	Rid,
		RotateDirection:1,
		XSpeed: 0,
		YSpeed: 0,
		ballRadius:ballRadius,
		session:session,
	}
	self.player[player.Rid]=player
	self.NotifyJoin(player) //广播给所有其他玩家
	result:=map[string]interface{}{
		"Rid":player.Rid,
		"Player":self.player,
		"Coins":self.coins,
	}
	return result
}

func (self *table)Remove(Rid string)(error) {
	for _,role:=range self.player{
		if role.Rid==Rid{
			delete(self.player,Rid)
		}
	}
	return nil
}
/**
玩家点击屏幕开始奔跑
 */
func (self *table)Fire(Rid string,X float64, Y float64,angle float64,power float64) {
	//发射
	if player,ok:=self.player[Rid];ok{
		player.Fire(X,Y,angle,power)
	}
}
func (self *table)Start(){
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
	go func() {
		tickNotify := time.NewTicker(100 * time.Millisecond)
		defer func() {
			tickNotify.Stop()
		}()
		for !self.stoped{
			select {
			case <-tickNotify.C:
				self.NotifyAxes(nil)
			}
		}
	}()
}
func (self *table)Stop(){
	self.stoped=true
}

/**
定帧计算所有玩家的位置
 */
func (self *table)Update(arge interface{}){
	startTime := time.Now().UnixNano()
	for _,role:=range self.player{
		self.wallBounce(role)
		role.Move()
		role.Rotate()
	}
	s :=(time.Now().UnixNano()-self.proTime)/1000000
	self.proTime = time.Now().UnixNano()
	if s>24{
		log.Debug("exct time %d et %d ns",s ,(time.Now().UnixNano()-startTime)/1000)
	}
}

/**
定期刷新所有玩家的位置
 */
func (self *table)NotifyAxes(arge interface{}){
	b, _ := json.Marshal(self.player)
	for _,role:=range self.player{
		e := role.session.Send("Hitball/OnMove", b)
		if e != "" {
			self.Remove(role.Rid)
		}
	}
	//if !self.stoped{
	//	timer.SetTimer(100,self.NotifyAxes,nil)
	//}

	if len(self.coins)<8{
		self.AddCoins()
	}
}

/**
通知所有玩家有新玩家加入
 */
func (self *table)NotifyJoin(player *player){
	b, _ := json.Marshal(player)
	for _,role:=range self.player{
		e := role.session.Send("Hitball/OnJoin", b)
		if e != "" {
			self.Remove(role.Rid)
		}
	}
}

/**
通知所有玩家新加了金币
 */
func (self *table)NotifyAddCoins(coins *coins){
	b, _ := json.Marshal(coins)
	for _,role:=range self.player{
		e := role.session.Send("Hitball/OnAddCoins", b)
		if e != "" {
			self.Remove(role.Rid)
		}
	}
}
/**
通知所有玩家金币已经被吃掉
 */
func (self *table)NotifyEatCoins(coins *coins){
	b, _ := json.Marshal(coins)
	for _,role:=range self.player{
		e := role.session.Send("Hitball/OnEatCoins", b)
		if e != "" {
			self.Remove(role.Rid)
		}
	}
}

/**
检查小球是否超出世界范围
 */
func (self *table) wallBounce (player *player){
	if(player.X<player.ballRadius){
		player.X=player.ballRadius;
		player.XSpeed*=-1
	}
	if(player.Y<player.ballRadius){
		player.Y=player.ballRadius;
		player.YSpeed*=-1
	}

	if(player.X>self.world_width-player.ballRadius){
		player.X=self.world_width-player.ballRadius;
		player.XSpeed*=-1
	}
	if(player.Y>self.world_height-player.ballRadius){
		player.Y=self.world_height-player.ballRadius;
		player.YSpeed*=-1
	}
}

