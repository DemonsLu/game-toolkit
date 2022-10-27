package main

import (
	"time"
)

const tickInterval = 10 // 每次tick 10ms
var startTs = time.Now().UnixMilli()

var tw TimeWheel
var ch = make(chan Timer, 1)

func InitTimerManager() {
	tw = TimeWheel{
		InstantId: divisionIdIndex,
		CurTick:   0, // 当前执行到了第几帧 (用于当服务器卡顿时的补帧操作)
		TimerMap:  make(map[uint64]*Timer),
		Wheels: []*Wheel{
			NewWheel(256),
			NewWheel(128),
			NewWheel(128),
			NewWheel(128),
			NewWheel(128),
		},
	}
}

func StartTimerManager() {
	go func() {
		for {
			nowTs := time.Now().UnixMilli()
			// 追帧
			for i := tw.CurTick + 1; i <= uint64(nowTs-startTs)/tickInterval; i++ {
				// 当前是否追到了最大帧
				isReached := i == uint64(nowTs-startTs)/tickInterval
				t1 := time.Now().UnixMilli()
				tw.ExecTick(i, i == uint64(nowTs-startTs)/tickInterval)
				t2 := time.Now().UnixMilli()
				// 追帧的时候别睡，直接一把追上，等追到最大帧的时候再睡
				if isReached && t2-t1 < 10 {
					time.Sleep(time.Duration(10-(t2-t1)) * time.Millisecond)
				}
			}
		}
	}()
}

func AddTimer(startOffsetMs int64, cb TimerCallback) (timerId uint64) {
	return tw.AddTimer(uint64(startOffsetMs/tickInterval), cb, 1, 0)
}

func AddMultiExecTimer(startOffsetMs int64, execCount int, intervalMs int64, cb TimerCallback) (timerId uint64) {
	return tw.AddTimer(uint64(startOffsetMs/tickInterval), cb, execCount, uint64(intervalMs/tickInterval))
}

func DeleteTimer(timerId uint64) (isSuccess bool) {
	return tw.DeleteTimer(timerId)
}
