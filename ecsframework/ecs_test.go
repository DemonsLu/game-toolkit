package main

import (
	"ecsframework/component"
	"ecsframework/consts"
	"ecsframework/entity"
	"ecsframework/system"
	"testing"
)

/*
	[思路]
	Component 和 System是 1 * 1 的关系，Component负责定义模块里的字段信息，System负责处理具体逻辑
	Entity 和 Component 是 m * n 的关系，每个Entity都可以注册不同的Component，以组合的方式完成对某一类Entity的模块绑定
	System 属于逻辑模块，在设计上最好和Entity & Component分开，这样这个package可以通过编译为.so文件的方式，实现逻辑的热更新
*/

func TestEcsFramework(t *testing.T) {
	r := entity.NewRole()
	c := r.GetComponent(consts.CSTypeWings).(*component.CWing)

	system.SWing.LevelUp(r)
	t.Logf("%+v", c)
	system.SWing.LevelUp(r)
	t.Logf("%+v", c)
	system.SWing.LevelUp(r)
	t.Logf("%+v", c)
	system.SWing.LevelUp(r)
	t.Logf("%+v", c)
	system.SWing.LevelUp(r)
	t.Logf("%+v", c)
}
