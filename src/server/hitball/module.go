/**
一定要记得在confin.json配置这个模块的参数,否则无法使用
*/
package hitball

import (
	"encoding/json"
	"github.com/liangdas/mqant/conf"
	"github.com/liangdas/mqant/gate"
	"github.com/liangdas/mqant/log"
	"github.com/liangdas/mqant/module"
)

var Module = func() module.Module {
	gate := new(Hitball)
	return gate
}

type Hitball struct {
	module.BaseModule
	proTime int64
	table *table
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
	self.table=NewTable()
	//self.SetListener(new(chat.Listener))
	self.GetServer().RegisterGO("HD_Move", self.move)
	self.GetServer().RegisterGO("HD_Join", self.join)
	self.GetServer().RegisterGO("HD_Fire", self.fire)
	self.GetServer().RegisterGO("HD_EatCoin", self.EatCoin)
}

func (self *Hitball) Run(closeSig chan bool) {
	self.table.Start()
}

func (self *Hitball) OnDestroy() {
	//一定别忘了关闭RPC
	self.GetServer().OnDestroy()
}

func (self *Hitball)join(s map[string]interface{}, msg map[string]interface{})(result map[string]interface{}, err string){
	//if msg["Rid"] == nil {
	//	err = "Rid cannot be nil"
	//	return
	//}
	//Rid := msg["Rid"].(string)
	session := gate.NewSession(self.App, s)
	result=self.table.Join(session.IP,session)
	return result,""
}

func (self *Hitball)fire(s map[string]interface{}, msg map[string]interface{})(result string, err string){
	if msg["Angle"] == nil ||msg["Power"] == nil||msg["X"] == nil ||msg["Y"] == nil{
		err = "Angle , Power X ,Y cannot be nil"
		return
	}
	//Rid := msg["Rid"].(string)
	session := gate.NewSession(self.App, s)
	Angle := msg["Angle"].(float64)
	Power := msg["Power"].(float64)
	X := msg["X"].(float64)
	Y := msg["Y"].(float64)
	self.table.Fire(session.IP,X,Y,Angle,Power)
	return "fire",""
}

func (self *Hitball)EatCoin(s map[string]interface{}, msg map[string]interface{})(result string, err string){
	if msg["Id"] == nil {
		err = "Id cannot be nil"
		return
	}
	//Rid := msg["Rid"].(string)
	session := gate.NewSession(self.App, s)
	Id := int(msg["Id"].(float64))
	self.table.EatCoins(session.IP,Id)
	return "EatCoin",""
}

func (self *Hitball) move(s map[string]interface{}, msg map[string]interface{}) (result string, err string) {
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
	session := gate.NewSession(self.App, s)
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
