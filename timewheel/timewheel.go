package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type TimerCallback func(param any)

type Timer struct {
	Id              uint64
	Callback        TimerCallback
	TriggerTick     uint64 // 下次触发应该在第几帧
	RemainExecCount int    // 剩余执行次数 -1表示无限制
	IntervalTick    uint64 // 触发tick间隔
	param           any
	IsDeleted       bool // 软删除
}

// slot当做循环队列使用，当slot index为len-1时，说明已经完整遍历完一整轮，此时需要让上级轮的current index事件全部pop，往本级轮中存放
type Wheel struct {
	CurIndex int
	Slots    []*NodeManager
}

func NewWheel(slotNum int) (wheel *Wheel) {
	wheel = &Wheel{
		CurIndex: 0,
		Slots:    make([]*NodeManager, 0, slotNum),
	}
	for i := 0; i < slotNum; i++ {
		wheel.Slots = append(wheel.Slots, NewNodeManager())
	}
	return
}

// TimeWheel
/*
	[思路]
	第一层的轮子为执行轮，类似于时钟的秒针，是时间轮一次转动的最小单位。走到对应的index时，需要真正执行对应index中的全部timer callback
	内层的轮子为缓存轮，类似于时钟的分针等，需要等秒针完整转动一周后，内层的轮子才往前推进一格，然后把current index里的全部timer放到上一个轮子里
	Slot直接使用单链表。在Timer里加个bool变量进行软删除即可。
*/

// 将无限次的 timerId 和 有限次的 timerId用不同段表示，防止程序长时间允许导致id溢出归零后，新的timer会和原有无限次执行的timerId重复导致错误
// 小于divisionIdIndex的是无限次timerId段，大于divisionIdIndex的是有限次timerId段
const divisionIdIndex = 10000000

type TimeWheel struct {
	InstantId uint64
	ForeverId uint64
	CurTick   uint64
	Wheels    []*Wheel
	mapMutex  sync.Mutex
	TimerMap  map[uint64]*Timer
}

// tick: 当前帧数
func (tw *TimeWheel) ExecTick(tick uint64, isReached bool) {
	// 把当前执行轮currentIndex的事件全部执行
	tw.CurTick = tick
	execWheel := tw.Wheels[0]
	execTimers := execWheel.Slots[execWheel.CurIndex].PopAll()
	for _, timer := range execTimers {
		// 如果是执行无限次数的timer，在追帧的过程中不要触发，防止过程中压力再次过大
		if !timer.IsDeleted && (timer.RemainExecCount != -1 || isReached) {
			timer.Callback(timer.param)
		}
		timer.TriggerTick = tw.CurTick + timer.IntervalTick
		timer.RemainExecCount--
		if timer.RemainExecCount != 0 && !timer.IsDeleted {
			tw.ReCalculateTimer(timer)
		} else {
			tw.mapMutex.Lock()
			delete(tw.TimerMap, timer.Id)
			tw.mapMutex.Unlock()
		}
	}
	execWheel.CurIndex = (execWheel.CurIndex + 1) % len(execWheel.Slots)

	// 执行轮已经转完了一圈，应该给内部轮中的curIndex中的东西出队进内部轮
	if execWheel.CurIndex == 0 {
		for i := 1; i < len(tw.Wheels); i++ {
			cacheWheel := tw.Wheels[i]
			cacheWheel.CurIndex = (cacheWheel.CurIndex + 1) % len(cacheWheel.Slots)

			timers := cacheWheel.Slots[cacheWheel.CurIndex].PopAll()
			for _, timer := range timers {
				tw.ReCalculateTimer(timer)
			}

			if cacheWheel.CurIndex != 0 {
				break
			}
		}
	}
}

func (tw *TimeWheel) ReCalculateTimer(timer *Timer) {
	isSuccess := false
	offset := timer.TriggerTick - tw.CurTick

	for i := 0; i < len(tw.Wheels); i++ {
		wheel := tw.Wheels[i]
		slotLen := len(wheel.Slots)
		if int(offset) < slotLen {
			targetIndex := (wheel.CurIndex + int(offset)) % slotLen
			wheel.Slots[targetIndex].Add(timer)
			isSuccess = true
			break
		}
		// eg: 当前是第35秒，有个距离当前100秒的timer。那么需要放在 (35 + 100) / 60 = 2，即放在分针的两格后的位置.因为还有25秒分针就要前进1格了
		offset = (offset + uint64(wheel.CurIndex)) / uint64(slotLen)
	}
	if !isSuccess {
		fmt.Println("add timer failed")
	}
}

func (tw *TimeWheel) AddTimer(startOffsetTick uint64, cb TimerCallback, execCount int, interval uint64) (timerId uint64) {
	timer := new(Timer)
	if execCount > 0 {
		timer.Id = atomic.AddUint64(&tw.InstantId, 1)
		if timer.Id < divisionIdIndex {
			if atomic.CompareAndSwapUint64(&tw.InstantId, 0, divisionIdIndex+1) {
				timer.Id = divisionIdIndex + 1
			} else {
				timer.Id = atomic.AddUint64(&tw.InstantId, 1)
			}
		}
	} else {
		timer.Id = atomic.AddUint64(&tw.ForeverId, 1)
		if timer.Id >= divisionIdIndex {
			if atomic.CompareAndSwapUint64(&tw.ForeverId, divisionIdIndex, 1) {
				timer.Id = 1
			} else {
				timer.Id = atomic.AddUint64(&tw.ForeverId, 1)
			}
		}
	}
	timer.TriggerTick = tw.CurTick + startOffsetTick
	timer.IntervalTick = interval
	timer.RemainExecCount = execCount
	timer.Callback = cb

	isSuccess := false
	offset := startOffsetTick
	for i := 0; i < len(tw.Wheels); i++ {
		wheel := tw.Wheels[i]
		slotLen := len(wheel.Slots)
		if int(offset) < slotLen {
			targetIndex := wheel.CurIndex + int(offset)
			wheel.Slots[targetIndex].Add(timer)
			isSuccess = true
			break
		}
		// eg: 当前是第35秒，有个距离当前100秒的timer。那么需要放在 (35 + 100) / 60 = 2，即放在分针的两格后的位置.因为还有25秒分针就要前进1格了
		offset = (offset + uint64(wheel.CurIndex)) / uint64(slotLen)
	}
	if !isSuccess {
		fmt.Println("add timer failed")
		return
	}
	tw.mapMutex.Lock()
	tw.TimerMap[timer.Id] = timer
	tw.mapMutex.Unlock()
	return timer.Id
}

func (tw *TimeWheel) DeleteTimer(timerId uint64) (isSuccess bool) {
	tw.mapMutex.Lock()
	timer, ok := tw.TimerMap[timerId]
	if ok && timer != nil {
		timer.IsDeleted = true
		isSuccess = true
	}
	tw.mapMutex.Unlock()
	return
}
