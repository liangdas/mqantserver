package main

import (
	. "robot/agent"
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	var robotCount = 5
	var taskCount=10000
	fmt.Println("start...count:", robotCount)

	now := time.Now().UnixNano()
	wg.Add(robotCount)
	for i := 0; i < robotCount; i++ {
		go Run(i,taskCount)
	}
	wg.Wait()
	use := time.Now().UnixNano() - now
	ms := use / int64(time.Millisecond)
	qps := int64(robotCount*taskCount) / (ms / 1000)
	fmt.Println("task over=>", " result:", "", robotCount*taskCount, "  usetime:ms", use/int64(time.Millisecond), " qps:", qps)
	fmt.Println("end")
}

func Run(index int,taskCount int) {

	acc := fmt.Sprintf("magicsea_%d", index)
	robot := NewRobot(acc, "111")
	err := robot.Start()
	if err!=nil{
		fmt.Println(err.Error())
	}
	robot.RunTask(taskCount)
	wg.Done()
}
