package propv2

import (
	"fmt"
	"testing"
)

const (
	rootModuleId  = 0
	equipModuleId = 1

	attackAbsolute           = 1
	attackRootPercent        = 2
	attackEquipPercent       = 3
	attackEquipLayer1Percent = 4
)

var tree *Tree

var equipModule EquipModule
var equipLevelModule EquipLevelModule
var equipRecastModule EquipRecastModule

type EquipLevelModule struct {
	n              *ModuleNode
	EquipLevelData interface{} // custom module data, blabla...
}
type EquipRecastModule struct {
	n               *ModuleNode
	EquipRecastData interface{} // custom module data, blabla...
}

type EquipModule struct {
	n         *ModuleNode
	EquipData interface{} // custom module data, blabla...
}

func initConfig() {
	attackAbsoluteConfig := PropConfig{
		PropId:       attackAbsolute,
		IsPercentage: false,
	}
	attackRootPercentConfig := PropConfig{
		PropId:          attackRootPercent,
		IsPercentage:    true,
		RelativePropId:  attackAbsolute,
		ConcernModuleId: rootModuleId,
	}
	attackEquipPercentConfig := PropConfig{
		PropId:          attackEquipPercent,
		IsPercentage:    true,
		RelativePropId:  attackAbsolute,
		ConcernModuleId: equipModuleId,
		ConcernLayer:    0,
	}
	attackEquipLevelPercentConfig := PropConfig{
		PropId:          attackEquipLayer1Percent,
		IsPercentage:    true,
		RelativePropId:  attackAbsolute,
		ConcernModuleId: equipModuleId,
		ConcernLayer:    1,
	}
	propConfig[attackAbsolute] = attackAbsoluteConfig
	propConfig[attackRootPercent] = attackRootPercentConfig
	propConfig[attackEquipPercent] = attackEquipPercentConfig
	propConfig[attackEquipLayer1Percent] = attackEquipLevelPercentConfig

	modulePropConfig[rootModuleId] = map[int]map[int]int{
		0: {attackAbsolute: attackRootPercent},
	}
	modulePropConfig[equipModuleId] = map[int]map[int]int{
		0: {attackAbsolute: attackEquipPercent},
		1: {attackAbsolute: attackEquipLayer1Percent},
	}

	tree = NewTree()
	equipModule = EquipModule{
		n:         tree.BuildByModule(equipModuleId),
		EquipData: nil,
	}
	equipLevelModule = EquipLevelModule{
		n:              tree.BuildBySubmodule(equipModule.n),
		EquipLevelData: nil,
	}
	equipRecastModule = EquipRecastModule{
		n:               tree.BuildBySubmodule(equipModule.n),
		EquipRecastData: nil,
	}
}

func TestProp(t *testing.T) {
	initConfig()

	tree.ChangeModuleProp(equipLevelModule.n, map[int]float64{
		attackAbsolute:           100,
		attackRootPercent:        10,
		attackEquipPercent:       20,
		attackEquipLayer1Percent: 50,
	})
	tree.ChangeModuleProp(equipRecastModule.n, map[int]float64{
		attackAbsolute:           200,
		attackRootPercent:        1,
		attackEquipPercent:       2,
		attackEquipLayer1Percent: 5,
	})
	// ((100 * 1.5) + (200 * 1.05)) * (1 + 0.2 + 0.02) * (1 + 0.1 + 0.01) = 360 * 1.22 * 1.11 = 487.512
	fmt.Printf("%+v\n", tree.RootNode.PropResult)
}
