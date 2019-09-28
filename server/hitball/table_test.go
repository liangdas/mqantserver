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
package hitball

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestTable(t *testing.T) {
	stoped := false
	table := NewTable()
	index := 0
	runtime.GOMAXPROCS(runtime.NumCPU())
	go func() {
		//这里设置为22ms但实际上每次循环大概是23-25ms左右，根据机器定,客户端设置的帧为40
		tick := time.NewTicker(22 * time.Millisecond)
		for !stoped {
			select {
			case <-tick.C:
				table.Update(nil)
			}
		}
	}()
	go func() {
		//这里设置为22ms但实际上每次循环大概是23-25ms左右，根据机器定,客户端设置的帧为40
		tick := time.NewTicker(6 * time.Millisecond)
		for index <= 5000 {
			select {
			case <-tick.C:
				err := table.PutQueue("EatCoins", "127.0.0.1", 1)
				if err != nil {
					fmt.Println("PutQueue", err.Error())
				}
				index++
			}
		}
		fmt.Println("PutQueue end 1")
	}()
	go func() {
		//这里设置为22ms但实际上每次循环大概是23-25ms左右，根据机器定,客户端设置的帧为40
		tick := time.NewTicker(6 * time.Millisecond)
		for index <= 5000 {
			select {
			case <-tick.C:
				err := table.PutQueue("Fire", "127.0.0.1", float64(30), float64(45), float64(23), float64(43))
				if err != nil {
					fmt.Println("PutQueue", err.Error())
				}
				index++
			}
		}
		fmt.Println("PutQueue end 2")
	}()
	for index <= 5000 {

	}
	stoped = true
	time.Sleep(1 * time.Second)

}
