package grid

import (
	"fmt"
	"sync"
	"sync/atomic"
)

/*
	[思路]
	九宫格AOI算法以自身所在格子为中心，其周围的八个格子为包围圈，总共九个格子。
	当玩家从自身所在的格子移动到其他格子时，重新计算新的九宫格。
	向离开的格子内的Entity发送leave sight命令，向进入的格子内的Entity发送enter sight

	[缺点]
	无法把想要的视野范围设置的很大，如果设置的视野范围超过了九宫格的范围，那将毫无意义。
	所以如果采用这种设计，客户端表现上的视野范围要小于等于九宫格的范围
*/

const gridSize = 10 // 给定格子边长

// 这里先假定地图无限大

var EntityIdGen int64
var GridManager *gridManager

func init() {
	GridManager = &gridManager{
		GridCache: make(map[uint64]*Grid),
	}
}

type gridManager struct {
	sync.Mutex
	GridCache map[uint64]*Grid
}

const (
	TypeEnterMap = iota
	TypeEnterSight
	TypeChangePosition
	TypeLeaveSight
	TypeLeaveMap
	TypeSync // 自身的探测，我移动之后能看见谁
)

func (g *gridManager) SetEntityInGrid(entity *Entity) {
	if entity == nil {
		return
	}
	gridId := entity.GridId
	g.Lock()
	grid, ok := g.GridCache[gridId]
	if !ok {
		grid = &Grid{EntityCache: make(map[int64]*Entity)}
		g.GridCache[gridId] = grid
	}
	g.Unlock()

	grid.Lock()
	grid.EntityCache[entity.Id] = entity
	grid.Unlock()
}

func (g *gridManager) RemoveEntityFromGrid(entity *Entity) {
	if entity == nil {
		return
	}
	gridId := entity.GridId
	g.Lock()
	grid, ok := g.GridCache[gridId]
	g.Unlock()
	if !ok {
		return
	}

	grid.RemoveEntity(entity)
}

func (g *gridManager) GetGridsByIds(gridIds []uint64) (result []*Grid) {
	g.Lock()
	defer g.Unlock()
	for _, gridId := range gridIds {
		grid, ok := g.GridCache[gridId]
		if !ok {
			continue
		}
		result = append(result, grid)
	}
	return
}

func (g *gridManager) GetGridsById(gridId uint64) (result *Grid) {
	g.Lock()
	defer g.Unlock()
	result = g.GridCache[gridId]
	return
}

type Grid struct {
	sync.Mutex
	EntityCache map[int64]*Entity
}

func (g *Grid) DoAction(entity *Entity, category int) {
	g.Lock()
	for id, e := range g.EntityCache {
		if id == entity.Id {
			continue
		}
		var data any
		switch category {
		case TypeEnterMap:
			data = "enter map"
		case TypeEnterSight:
			data = "enter sight"
		case TypeChangePosition:
			data = fmt.Sprintf("change position to x: %d, y: %d", entity.PositionX, entity.PositionY)
		case TypeLeaveSight:
			data = "leave sight"
		case TypeLeaveMap:
			data = "leave map"
		case TypeSync:
			data = fmt.Sprintf("sync, I can see %d on x: %d, y: %d", e.Id, e.PositionX, e.PositionY)
		}
		if category == TypeSync {
			entity.ReceiveNews(e.Id, data)
		} else {
			e.ReceiveNews(entity.Id, data)
		}
	}
	g.Unlock()
}

func (g *Grid) RemoveEntity(entity *Entity) {
	g.Lock()
	delete(g.EntityCache, entity.Id)
	g.Unlock()
}

type Entity struct {
	Id        int64
	GridId    uint64
	PositionX int
	PositionY int
}

func NewEntity() (e *Entity) {
	return &Entity{
		Id: atomic.AddInt64(&EntityIdGen, 1),
	}
}

func (e *Entity) ReceiveNews(entityId int64, data any) {
	fmt.Printf("%d received %d's message, data is %+v\n", e.Id, entityId, data)
}

func (e *Entity) EnterMap(x, y int) {
	// 1. 根据x & y算出gridId
	e.PositionX, e.PositionY = x, y
	e.setGridId()
	GridManager.SetEntityInGrid(e)

	grids := GridManager.GetGridsByIds(getSightGridIds(e.GridId))
	for _, grid := range grids {
		grid.DoAction(e, TypeEnterMap)
		grid.DoAction(e, TypeSync)
	}
}

func (e *Entity) LeaveMap() {
	grids := GridManager.GetGridsByIds(getSightGridIds(e.GridId))
	for _, grid := range grids {
		grid.DoAction(e, TypeLeaveMap)
	}
	GridManager.RemoveEntityFromGrid(e)
}

func (e *Entity) ChangePosition(x, y int) {
	originGridId := e.GridId
	e.PositionX, e.PositionY = x, y
	e.setGridId()
	currentGridId := e.GridId
	currentGridIds := getSightGridIds(e.GridId)
	currentGrids := GridManager.GetGridsByIds(currentGridIds)

	// 这里表示伴随着格子的切换，需要发送进/出的消息给对应的玩家
	// origin从有到无发leave sight，current从无到有发enter sight
	if currentGridId != originGridId {
		oldGrid := GridManager.GetGridsById(originGridId)
		oldGrid.RemoveEntity(e)

		originGridIds := getSightGridIds(originGridId)
		currentGridMap := changeSlice2Map(currentGridIds)
		originGridMap := changeSlice2Map(originGridIds)

		result := make(map[uint64]bool) // true表示enter, false表示leave
		for id := range currentGridMap {
			if _, ok := originGridMap[id]; !ok {
				result[id] = true
			}
		}
		for id := range originGridMap {
			if _, ok := currentGridMap[id]; !ok {
				result[id] = false
			}
		}
		for id, v := range result {
			grid := GridManager.GetGridsById(id)
			if grid == nil {
				continue
			}
			if v {
				grid.DoAction(e, TypeEnterSight)
				grid.DoAction(e, TypeSync)
			} else {
				grid.DoAction(e, TypeLeaveSight)
			}
		}
	}

	// 发送坐标改变信息
	for _, grid := range currentGrids {
		grid.DoAction(e, TypeChangePosition)
	}
}

// 根据位置计算格子ID
func (e *Entity) setGridId() {
	e.GridId = calculateGridIdByXY(uint64(e.PositionX/gridSize), uint64(e.PositionY/gridSize))
}

func changeSlice2Map(s []uint64) (m map[uint64]struct{}) {
	m = make(map[uint64]struct{}, len(s))
	for _, v := range s {
		m[v] = struct{}{}
	}
	return
}

// 九宫格
func getSightGridIds(gridId uint64) (gridIds []uint64) {
	xIndex, yIndex := calculateBYByGridId(gridId)
	gridIds = append(gridIds, calculateGridIdByXY(xIndex-1, yIndex-1))
	gridIds = append(gridIds, calculateGridIdByXY(xIndex, yIndex-1))
	gridIds = append(gridIds, calculateGridIdByXY(xIndex+1, yIndex-1))
	gridIds = append(gridIds, calculateGridIdByXY(xIndex-1, yIndex))
	gridIds = append(gridIds, calculateGridIdByXY(xIndex, yIndex))
	gridIds = append(gridIds, calculateGridIdByXY(xIndex+1, yIndex))
	gridIds = append(gridIds, calculateGridIdByXY(xIndex-1, yIndex+1))
	gridIds = append(gridIds, calculateGridIdByXY(xIndex, yIndex+1))
	gridIds = append(gridIds, calculateGridIdByXY(xIndex+1, yIndex+1))
	return gridIds
}

func calculateGridIdByXY(xIndex, yIndex uint64) (gridId uint64) {
	return xIndex<<32 + yIndex
}

func calculateBYByGridId(gridId uint64) (xIndex, yIndex uint64) {
	xIndex, yIndex = gridId>>32, gridId&0x00000000FFFFFFFF
	return
}
