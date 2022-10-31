package infinitySight

import (
	"fmt"
	"sync"
)

/*
	[思路]
	无限视野的方案实现最简单，适用于场景内Entity比较少的情况
	所有Entity的可观测行为均需要广播给在场景内的其他Entity
*/

var Manager manager

type manager struct {
	sync.Mutex
	EntityMap map[int64]*Entity
}

type Entity struct {
	Id        int64
	PositionX int
	PositionY int
}

func init() {
	Manager = manager{
		EntityMap: make(map[int64]*Entity),
	}
}

func (e *Entity) EnterMap() {
	Manager.Lock()
	defer Manager.Unlock()
	for _, entity := range Manager.EntityMap {
		entity.ReceiveNews(e.Id, "enterMap")
	}
	Manager.EntityMap[e.Id] = e
}

func (e *Entity) LeaveMap() {
	Manager.Lock()
	defer Manager.Unlock()
	delete(Manager.EntityMap, e.Id)
	for _, entity := range Manager.EntityMap {
		entity.ReceiveNews(e.Id, "leaveMap")
	}
}

func (e *Entity) ChangePosition(x, y int) {
	Manager.Lock()
	defer Manager.Unlock()

	e.PositionX, e.PositionY = x, y
	for id, entity := range Manager.EntityMap {
		if id == e.Id {
			continue
		}
		entity.ReceiveNews(e.Id, fmt.Sprintf("change position to: x %d, y %d", x, y))
	}
}

func (e *Entity) ReceiveNews(entityId int64, data any) {
	fmt.Printf("%d received %d's message, data is %+v\n", e.Id, entityId, data)
}
