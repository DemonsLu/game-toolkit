package crosschain

import (
	"fmt"
	"math"
	"sync/atomic"
)

const (
	_                = iota
	NodeSentinelDown // 下界哨兵
	NodeSentinelUp   // 上界哨兵
	NodeEntity

	AxisX = 1
	AxisY = 2

	InvalidSight = -1
)

/*
	[思路]
	在X轴和Y轴方向上各建立一个链表，故称之为十字链表
	以自身(比如小A)的视角来看，自己的左右哨兵覆盖范围内的所有Entity，都能被自己看到
	以其他Entity的视角来看，其他Entity的左右哨兵覆盖范围内如果有小A，那么其他Entity可以看到小A
	优势: 每个Entity都可以有自己的视野范围，视野设置更加灵活
	缺点: 对CPU的消耗较大
*/

var EntityIdGen int64

var NodeManager *nodeManager

type nodeManager struct {
	XNodeHead *Node
	YNodeHead *Node

	entityMap map[int64]*Entity
}

func init() {
	NodeManager = &nodeManager{
		XNodeHead: nil,
		YNodeHead: nil,
		entityMap: make(map[int64]*Entity),
	}
}

func (m *nodeManager) AddEntity(e *Entity) {
	if e == nil {
		return
	}

	m.addNode(e.XNode[0], AxisX)
	m.addNode(e.XNode[1], AxisX)
	m.addNode(e.XNode[2], AxisX)
	m.addNode(e.YNode[0], AxisY)
	m.addNode(e.YNode[1], AxisY)
	m.addNode(e.YNode[2], AxisY)
	m.entityMap[e.Id] = e

	// 我能看到谁
	for _, entity := range m.GetSelfSightEntities(e.Id) {
		e.SendMessage(fmt.Sprintf("I Enter Map, I can see %d", entity.Id))
	}
	// 谁能看到我
	for _, entity := range m.GetOtherSightEntities(e.Id) {
		entity.SendMessage(fmt.Sprintf("I can see %d enter map right now", e.Id))
	}
}

func (m *nodeManager) RemoveEntity(e *Entity) {
	if e == nil {
		return
	}

	// 我能看到谁
	for _, entity := range m.GetSelfSightEntities(e.Id) {
		e.SendMessage(fmt.Sprintf("I Leave Map, so I can see %d any more", entity.Id))
	}
	// 谁能看到我
	for _, entity := range m.GetOtherSightEntities(e.Id) {
		entity.SendMessage(fmt.Sprintf("I can see %d leave map right now", e.Id))
	}

	m.RemoveNode(e.XNode[0], AxisX)
	m.RemoveNode(e.XNode[1], AxisX)
	m.RemoveNode(e.XNode[2], AxisX)
	m.RemoveNode(e.YNode[0], AxisY)
	m.RemoveNode(e.YNode[1], AxisY)
	m.RemoveNode(e.YNode[2], AxisY)
	delete(m.entityMap, e.Id)
}

func (m *nodeManager) ChangePosition(e *Entity, x, y int) {
	// 我能看到谁
	selfSights := make(map[int64]struct{})
	otherSights := make(map[int64]struct{})
	for _, entity := range m.GetSelfSightEntities(e.Id) {
		selfSights[entity.Id] = struct{}{}
	}
	// 谁能看到我
	for _, entity := range m.GetOtherSightEntities(e.Id) {
		otherSights[entity.Id] = struct{}{}
	}

	//todo: 先暂且这么写，把整体思路先实现，代码后续考虑优化
	xDiff := x - e.XNode[1].Value
	yDiff := y - e.YNode[1].Value
	m.RemoveNode(e.XNode[0], AxisX)
	m.RemoveNode(e.XNode[1], AxisX)
	m.RemoveNode(e.XNode[2], AxisX)
	m.RemoveNode(e.YNode[0], AxisY)
	m.RemoveNode(e.YNode[1], AxisY)
	m.RemoveNode(e.YNode[2], AxisY)
	e.XNode[0].Value += xDiff
	e.XNode[1].Value += xDiff
	e.XNode[2].Value += xDiff
	e.YNode[0].Value += yDiff
	e.YNode[1].Value += yDiff
	e.YNode[2].Value += yDiff
	m.addNode(e.XNode[0], AxisX)
	m.addNode(e.XNode[1], AxisX)
	m.addNode(e.XNode[2], AxisX)
	m.addNode(e.YNode[0], AxisY)
	m.addNode(e.YNode[1], AxisY)
	m.addNode(e.YNode[2], AxisY)

	for _, entity := range m.GetSelfSightEntities(e.Id) {
		if _, ok := selfSights[entity.Id]; !ok {
			// 我新看到了谁
			e.SendMessage(fmt.Sprintf("I change position, I can see new friend %d", entity.Id))
		} else {
			delete(selfSights, entity.Id)
		}
	}
	// 谁从我的视野里消失了
	for id := range selfSights {
		sightOut := m.entityMap[id]
		e.SendMessage(fmt.Sprintf("I change position, I cannot see %d", sightOut.Id))
	}

	for _, entity := range m.GetOtherSightEntities(e.Id) {
		if _, ok := otherSights[entity.Id]; !ok {
			// 谁从现在开始看到了我
			entity.SendMessage(fmt.Sprintf("I can see a new friend %d", e.Id))
		} else {
			delete(otherSights, entity.Id)
		}
	}
	// 我从谁的视野里消失了
	for id := range otherSights {
		sightOut := m.entityMap[id]
		sightOut.SendMessage(fmt.Sprintf("I cannot see %d", e.Id))
	}
}

// GetSelfSightEntities 获取自己能看到的Entity列表
func (m *nodeManager) GetSelfSightEntities(id int64) (result []*Entity) {
	self, ok := m.entityMap[id]
	if !ok {
		return
	}

	xIds := make(map[int64]struct{})
	for n := self.XNode[0]; n != self.XNode[2]; n = n.Next {
		if n.Category != NodeEntity || n.ID == id {
			continue
		}
		xIds[n.ID] = struct{}{}
	}

	for n := self.YNode[0]; n != self.YNode[2]; n = n.Next {
		if n.Category != NodeEntity || n.ID == id {
			continue
		}
		// X轴和Y轴均有交集，那么说明是要找的entity
		if _, ok := xIds[n.ID]; ok {
			result = append(result, m.entityMap[n.ID])
		}
	}
	return
}

// 获取能看到自己的Entity列表
func (m *nodeManager) GetOtherSightEntities(id int64) (result []*Entity) {
	self, ok := m.entityMap[id]
	if !ok {
		return
	}

	xIds := make(map[int64]struct{})
	for n := m.XNodeHead; n != nil; n = n.Next {
		if n.Category != NodeEntity || n.ID == id {
			continue
		}
		other := m.entityMap[n.ID]
		if int(math.Abs(float64(n.Value-self.XNode[1].Value))) < other.Sight {
			xIds[n.ID] = struct{}{}
		}
	}

	for n := m.YNodeHead; n != nil; n = n.Next {
		if n.Category != NodeEntity || n.ID == id {
			continue
		}
		other := m.entityMap[n.ID]
		if int(math.Abs(float64(n.Value-self.YNode[1].Value))) >= other.Sight {
			continue
		}
		// X轴和Y轴均有交集，那么说明是要找的entity
		if _, ok := xIds[n.ID]; ok {
			result = append(result, m.entityMap[n.ID])
		}
	}
	return
}

func (m *nodeManager) addNode(target *Node, category int) {
	if category == AxisX && m.XNodeHead == nil {
		m.XNodeHead = target
		return
	}
	if category == AxisY && m.YNodeHead == nil {
		m.YNodeHead = target
		return
	}

	var header *Node
	if category == AxisX {
		header = m.XNodeHead
	} else {
		header = m.YNodeHead
	}
	for n := header; n != nil; n = n.Next {
		// 找到了要插入的node，target应该放到n的前面
		if target.Value < n.Value {
			target.Next = n
			if n.Front != nil {
				target.Front = n.Front
				n.Front.Next = target
				n.Front = target
			}
			return
		}

		// 说明value是最大的，此时替换即可
		if n.Next == nil {
			n.Next = target
			target.Front = n
			return
		}
	}
}

func (m *nodeManager) RemoveNode(target *Node, category int) {
	// 如果节点是头节点，则需要特殊处理
	if category == AxisX {
		if m.XNodeHead == target {
			if m.XNodeHead.Next == nil {
				m.XNodeHead = nil
			} else {
				m.XNodeHead = target.Next
			}
		}
	} else {
		if m.YNodeHead == target {
			if m.YNodeHead.Next == nil {
				m.YNodeHead = nil
			} else {
				m.YNodeHead = target.Next
			}
		}
	}

	if target.Front != nil {
		target.Front.Next = target.Next
	}
	if target.Next != nil {
		target.Next.Front = target.Front
	}
	target.Front = nil
	target.Next = nil
}

func (m *nodeManager) changeNode(target *Node, category int, value int) {
	// todo: finish this
}

type Node struct {
	ID       int64 // entityId
	Category int
	Value    int
	Front    *Node
	Next     *Node
}

type Entity struct {
	Id int64

	Sight int      // 视野距离
	XNode [3]*Node // 0:左哨兵 1:自身 2:右哨兵
	YNode [3]*Node // 0:下哨兵 1:自身 2:上哨兵
}

func NewEntity() (e *Entity) {
	return &Entity{Id: atomic.AddInt64(&EntityIdGen, 1), Sight: InvalidSight}
}

func (e *Entity) EnterMap(sight, x, y int) {
	id := e.Id
	x1sentinel := &Node{ID: id, Category: NodeSentinelDown, Value: x - sight}
	xNode := &Node{ID: id, Category: NodeEntity, Value: x}
	x2sentinel := &Node{ID: id, Category: NodeSentinelUp, Value: x + sight}

	y1sentinel := &Node{ID: id, Category: NodeSentinelDown, Value: y - sight}
	yNode := &Node{ID: id, Category: NodeEntity, Value: y}
	y2sentinel := &Node{ID: id, Category: NodeSentinelUp, Value: y + sight}

	e.Sight = sight
	e.XNode = [3]*Node{x1sentinel, xNode, x2sentinel}
	e.YNode = [3]*Node{y1sentinel, yNode, y2sentinel}
	NodeManager.AddEntity(e)
}

func (e *Entity) LeaveMap() {
	// 说明没有进入Map
	if e.Sight == InvalidSight {
		return
	}
	NodeManager.RemoveEntity(e)
}

func (e *Entity) ChangePosition(x, y int) {
	NodeManager.ChangePosition(e, x, y)
}

func (e *Entity) SendMessage(data any) {
	fmt.Printf("%d receive message: %+v\n", e.Id, data)
}
