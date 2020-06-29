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

package scheduler

import (
	"errors"
	"fmt"
	"time"
	"github.com/DemonZSD/timerscheduler/tasks"
)

//TIME_WHEEL_SIZE the length of time wheels
const (
	MOD             = 3
	TIME_WHEEL_SIZE time.Duration = 2 << MOD
)

//TimeWheels is the time wheel witch includs TimeWheelItem
type TimeWheels struct {
	TimeWheel map[int]TimeWheelItem
}

//TimeWheelItem the time wheel's item
type TimeWheelItem struct {
	Size            int
	TimerSchedulers map[int]tasks.TimerScheduler
}

//TimeWheelInterface allows to operate TimerScheduler task, such as AddSchedule
type TimeWheelInterface interface {
	//add new TimerScheduler to TimeWheels
	AddSchedule(timerScheduler *tasks.TimerScheduler) error
	//update scheduler, such as the cycle when
	UpdateSchedule(timerScheduler *tasks.TimerScheduler) error
	//start running task
	Start(timeWheels *TimeWheels) error
}

//New a time wheel
func New() *TimeWheels {
	return &TimeWheels{
		make(map[int]TimeWheelItem),
	}
}

//AddScheduler add new TimerScheduler to TimeWheels
func (timeWheels *TimeWheels) AddScheduler(timerScheduler *tasks.TimerScheduler) error {
	if timerScheduler == nil {
		return errors.New("the timer-scheduler can not be empty")
	}

	if timerScheduler.Interval <= 0 {
		return errors.New("the timer-scheduler exec interval should be bigger then zero")
	}
	// copy a new timer scheduler
	newTimerScheduler := timerScheduler
	//calculate the timerScheduler wheel index and the cycle number by interval
	wheelIndex, cycle := GetTimerIndexAndCycle(timerScheduler)
	if timeWheels.TimeWheel == nil {
		return errors.New("the timer wheel is not init")
	}
	newTimerScheduler.CycNum = cycle
	timerWheelItem := timeWheels.TimeWheel[wheelIndex]
	// if the timer wheel grade's item is nil, init a new one.
	if timerWheelItem.TimerSchedulers == nil {
		timerWheelItem.Size = 0
		timerWheelItem.TimerSchedulers = make(map[int]tasks.TimerScheduler)
	}
	//put the timer Scheduler into the map, the key is the timerWheelItem.Size
	timerWheelItem.Size++
	timerWheelItem.TimerSchedulers[timerWheelItem.Size] = *newTimerScheduler
	timeWheels.TimeWheel[wheelIndex] = timerWheelItem
	return nil
}

//UpdateScheduler update scheduler cycle when time is comming, the cycle = cycle - 1
func (timeWheels *TimeWheels) UpdateScheduler(timerScheduler *tasks.TimerScheduler, updateIndex int) error {
	if timerScheduler == nil {
		return errors.New("the timer-scheduler can not be empty")
	}

	if timerScheduler.Interval <= 0 {
		return errors.New("the timer-scheduler exec interval should be bigger then zero")
	}
	newTimerScheduler := timerScheduler
	wheelIndex, cycle := GetTimerIndexAndCycle(timerScheduler)
	if timerScheduler.CycNum == -1 {
		newTimerScheduler.CycNum = cycle
	}

	if timeWheels.TimeWheel == nil {
		return errors.New("the timer wheel is not init")
	}
	timeWheelItem := timeWheels.TimeWheel[wheelIndex]

	if timeWheelItem.TimerSchedulers == nil {
		return errors.New("timer schedulers can not be nil")
	}
	//update the new timerScheduler to timeWheelItem.TimerSchedulers[updateIndex]
	timeWheelItem.TimerSchedulers[updateIndex] = *newTimerScheduler
	return nil
}

//GetTimerIndexAndCycle get the timer scheduler interval and cycle number
func GetTimerIndexAndCycle(timerScheduler *tasks.TimerScheduler) (int, int) {
	return int(timerScheduler.Interval&TIME_WHEEL_SIZE - 1), int(timerScheduler.Interval >> MOD)
}

//Start the timer wheel
func Start(timeWheels *TimeWheels) error {
	fmt.Println("Start the timer scheduler")
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		var wheelIndex = 0
		//each ticker to exec this statement
		for _ = range ticker.C {
			if timeWheels.TimeWheel[wheelIndex].Size > 0 {
				if timeWheels.TimeWheel[wheelIndex].Size == len(timeWheels.TimeWheel[wheelIndex].TimerSchedulers) {
					fmt.Println("The size is not eq with length of TimerScheduler array")
				}
				for scheduleIndex, timerScheduler := range timeWheels.TimeWheel[wheelIndex].TimerSchedulers {
					if timerScheduler.CycNum == 0 {
						timerScheduler.DefaultExecuteSchedule()
					} else {
						timerScheduler.CycNum--
						timeWheels.UpdateScheduler(&timerScheduler, scheduleIndex)
					}
				}
			}
			if wheelIndex == int(TIME_WHEEL_SIZE-1) {
				wheelIndex = 0
			} else {
				wheelIndex++
			}

		}
	}()
	return nil
}
