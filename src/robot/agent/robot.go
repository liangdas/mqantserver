package agent

import (
	"fmt"
	_ "time"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type Robot struct {
	account string
	pwd     string

	gateAddr string
	uid      uint64
	key      string
	result	 chan MQTT.Message
	finish	 chan int
	resultNum int
	taskCount int
	client   MQTT.Client
}

func NewRobot(account, pwd string) *Robot {
	return &Robot{account: account, pwd: pwd}
}

func (robot *Robot) Start()error  {
	return robot.ConnectGate()
}


func (robot *Robot) ConnectGate() error{
	fmt.Println("ConnectGate...")
	robot.result=make(chan MQTT.Message)
	robot.finish=make(chan int)
	robot.resultNum=0
	opts := MQTT.NewClientOptions()
	opts.AddBroker("tcp://127.0.0.1:3563")
	opts.SetClientID("123")
	opts.SetUsername(robot.account)
	opts.SetPassword(robot.pwd)
	opts.SetCleanSession(false)
	opts.SetProtocolVersion(3)
	opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
		//robot.result <- msg
		//fmt.Println("publish",robot.resultNum,msg.Topic(),string(msg.Payload()))
		robot.resultNum++
		if robot.resultNum>=robot.taskCount{
			robot.finish<-0
		}
	})

	robot.client =  MQTT.NewClient(opts)
	if token := robot.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (robot *Robot)RunTask(taskCount int)  {
	robot.taskCount=taskCount
	for i := 0; i < robot.taskCount; i++ {
		//b,_:=json.Marshal(map[string]string{"userName": "liangdas", "passWord": "Hello,anyone!"})
		//fmt.Println("send",i)
		//time.Sleep(time.Millisecond*20)
		s:=fmt.Sprintf("'userName': 'liangdas', 'passWord': 'Hello,anyone! %d'}",i)
		robot.client.Publish("Login/HD_Robot/1",0,false,[]byte(s))
		//if token.Wait() && token.Error() != nil {
		//	panic(token.Error())
		//}
		//msg:=<-robot.result
		//fmt.Println("publish",msg.Topic(),string(msg.Payload()))
	}
	<-robot.finish
}

func (robot *Robot) Finish() {
	robot.client.Disconnect(250)
}
