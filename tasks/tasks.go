// Copyright 2020 Weshzhu
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

// author: Weshzhu

package tasks

import (
	"time"
)

// TimerSchedulerInterface The timer-scheduler interface allows tasks witch extend it can do something what you want.
type TimerSchedulerInterface interface {
	DefaultExecuteSchedule(exec func (Args *interface{}) interface{}) error // All tasks should be extend this method. then the tasks can be add to scheduler timer-wheels
	CustomExecuteSchedule() error // All tasks should be extend this method. then the tasks can be add to scheduler timer-wheels
}
type  TimerSchedulerExec struct {
	Args *interface{}
	Exec func (Args *interface{}) interface{}
}
type TimerScheduler struct {
	CycNum   int // this time scheduler
	Interval time.Duration // Time interval of scheduled task
	TimerSchedulerExec
}

// DefaultExecuteSchedule timerScheduler executes the customize task,
func (timerScheduler *TimerScheduler) CustomExecuteSchedule(exec func (Args *interface{}) interface{}) interface{} {
	return exec(timerScheduler.Args)
}

// ExecuteSchedule timerScheduler executes task, if use the specify executor,
// you should extends the TimerScheduler
// and implements this func
func (timerScheduler *TimerScheduler) DefaultExecuteSchedule() interface{} {
	return timerScheduler.Exec(timerScheduler.Args)
}


