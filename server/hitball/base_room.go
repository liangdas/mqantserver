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
	"github.com/liangdas/mqant/module"
	"sync"
)

type Room struct {
	module module.Module
	lock   *sync.RWMutex
	tables map[int]*Table
	index  int
	max    int
}

func NewRoom(module module.Module) *Room {
	room := &Room{
		module: module,
		lock:   new(sync.RWMutex),
		tables: map[int]*Table{},
		index:  0,
		max:    0,
	}
	return room
}

func (self *Room) create(module module.Module) *Table {
	self.lock.Lock()
	self.index++
	table := NewTable(module, self.index)
	self.tables[self.index] = table
	self.lock.Unlock()
	return table
}

func (self *Room) GetTable(tableId int) *Table {
	if table, ok := self.tables[tableId]; ok {
		return table
	}
	return nil
}

func (self *Room) GetEmptyTable() (*Table, error) {
	for _, table := range self.tables {
		if table.Empty() {
			return table, nil
		}
	}
	//没有找到已创建的空房间,新创建一个
	table := self.create(self.module)
	if table == nil {
		return nil, fmt.Errorf("fail create table")
	}
	return table, nil
}
