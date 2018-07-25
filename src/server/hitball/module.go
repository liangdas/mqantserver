/**
一定要记得在confin.json配置这个模块的参数,否则无法使用
*/
package hitball

import (
	"math/rand"
	"encoding/json"
	"github.com/liangdas/mqant/conf"
	"github.com/liangdas/mqant/gate"
	"github.com/liangdas/mqant/log"
	"github.com/liangdas/mqant/module"
	"github.com/liangdas/mqant/module/base"
	"time"
)

var Module = func() module.Module {
	gate := new(Hitball)
	return gate
}

type Hitball struct {
	basemodule.BaseModule
	room	*Room
	proTime int64
	table *Table
}
//生成随机字符串
func GetRandomString(lenght int) string{
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < lenght; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
func (self *Hitball) GetType() string {
	//很关键,需要与配置文件中的Module配置对应
	return "Hitball"
}
func (self *Hitball) Version() string {
	//可以在监控时了解代码版本
	return "1.0.0"
}
func (self *Hitball) OnInit(app module.App, settings *conf.ModuleSettings) {
	self.BaseModule.OnInit(self, app, settings)
	self.room=NewRoom(self)
	self.table,_=self.room.GetEmptyTable()
	//self.SetListener(new(chat.Listener))
	self.GetServer().RegisterGO("HD_Move", self.move)
	self.GetServer().RegisterGO("HD_Join", self.join)
	self.GetServer().RegisterGO("HD_Fire", self.fire)
	self.GetServer().RegisterGO("HD_EatCoin", self.eatCoin)
}

func (self *Hitball) Run(closeSig chan bool) {
	self.table.Start()
}

func (self *Hitball) OnDestroy() {
	//一定别忘了关闭RPC
	self.GetServer().OnDestroy()
}

func (self *Hitball)join(session gate.Session, msg map[string]interface{})(result string, err string){
	if session.GetUserId()==""{
		session.Bind(GetRandomString(8))
		//return "","no login"
	}
	erro:=self.table.PutQueue("Join",session)
	if erro!=nil{
		return "",erro.Error()
	}
	return "success",""
}

func (self *Hitball)fire(session gate.Session, msg map[string]interface{})(result string, err string){
	if msg["Angle"] == nil ||msg["Power"] == nil||msg["X"] == nil ||msg["Y"] == nil{
		err = "Angle , Power X ,Y cannot be nil"
		return
	}
	Angle := msg["Angle"].(float64)
	Power := msg["Power"].(float64)
	X := msg["X"].(float64)
	Y := msg["Y"].(float64)
	erro:=self.table.PutQueue("Fire",session,float64(X),float64(Y),float64(Angle),float64(Power))
	if erro!=nil{
		return "",erro.Error()
	}
	return "success",""
}

func (self *Hitball)eatCoin(session gate.Session, msg map[string]interface{})(result string, err string){
	if msg["Id"] == nil {
		err = "Id cannot be nil"
		return
	}
	Id := int(msg["Id"].(float64))
	erro:=self.table.PutQueue("EatCoins",session,Id)
	if erro!=nil{
		return "",erro.Error()
	}
	return "success",""
}

func (self *Hitball) move(session gate.Session, msg map[string]interface{}) (result string, err string) {
	if msg["war"] == nil || msg["wid"] == nil || msg["x"] == nil || msg["y"] == nil {
		err = "war , wid ,x ,y cannot be nil"
		return
	}
	//log.Debug("exct time %d", (time.Now().UnixNano()-self.proTime)/1000000)
	//self.proTime = time.Now().UnixNano()
	//war := msg["war"].(string)
	//wid := msg["wid"].(string)
	x := msg["x"].(float64)
	y := msg["y"].(float64)
	//passWord:=msg["passWord"].(string)
	roles := []map[string]float64{
		map[string]float64{
			"x": x,
			"y": y,
		},
	}
	re := map[string]interface{}{}
	re["roles"] = roles
	b, _ := json.Marshal(re)
	e := session.SendNR("Hitball/OnMove", b)
	if e != "" {
		log.Error(e)
	}
	//log.Debug(fmt.Sprintf("move success x:%v,y:%v", x, y))
	return "success", ""
}
