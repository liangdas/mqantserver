// Copyright 2014 hey Author. All Rights Reserved.
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
package xaba_task

import (
	"github.com/liangdas/armyant/task"
	"io"
	"os"
)

type Manager struct {
	// Writer is where results will be written. If nil, results are written to stdout.
	Writer io.Writer
}

func (this *Manager) writer() io.Writer {
	if this.Writer == nil {
		return os.Stdout
	}
	return this.Writer
}
func (this *Manager) Finish(task task.Task) {
	//total := time.Now().Sub(task.Start)
}
func (this *Manager) CreateWork() task.Work {
	return NewWork(this)
}

// Run makes all the requests, prints the summary. It blocks until
// all work is done.
func NewManager(t task.Task) task.WorkManager {
	// append hey's user agent
	this := new(Manager)
	return this
}
