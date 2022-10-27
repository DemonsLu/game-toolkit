package system

import (
	"ecsframework/component"
	"ecsframework/consts"
	"ecsframework/entity"
)

// you know, for example

type sWing struct{}

var SWing sWing

func (s *sWing) LevelUp(entity entity.IEntity) {
	if entity == nil {
		return
	}
	c := entity.GetComponent(consts.CSTypeWings)
	if c == nil {
		return
	}
	cWing, ok := c.(*component.CWing)
	if cWing == nil || !ok {
		return
	}
	cWing.WingLevel++
	// 每升5级进一阶
	if cWing.WingLevel%5 == 0 {
		cWing.WingStage++
	}
}
