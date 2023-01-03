package prop

import (
	"fmt"
	"testing"
)

const (
	rootModuleId       = 1
	equipModuleId      = 10
	equipLevelModuleId = 100

	attackAbsolute          = 1
	attackRootPercent       = 2
	attackEquipPercent      = 3
	attackEquipLevelPercent = 4
)

func initConfig() {
	rootModuleConfigId = rootModuleId
	rootModule := ModuleConfig{
		ModuleId:       rootModuleId,
		ParentModuleId: 0,
	}
	equipModule := ModuleConfig{
		ModuleId:       equipModuleId,
		ParentModuleId: rootModuleId,
	}
	equipLevelModule := ModuleConfig{
		ModuleId:       equipLevelModuleId,
		ParentModuleId: equipModuleId,
	}
	moduleConfig[rootModuleId] = rootModule
	moduleConfig[equipModuleId] = equipModule
	moduleConfig[equipLevelModuleId] = equipLevelModule

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
	}
	attackEquipLevelPercentConfig := PropConfig{
		PropId:          attackEquipLevelPercent,
		IsPercentage:    true,
		RelativePropId:  attackAbsolute,
		ConcernModuleId: equipLevelModuleId,
	}
	propConfig[attackAbsolute] = attackAbsoluteConfig
	propConfig[attackRootPercent] = attackRootPercentConfig
	propConfig[attackEquipPercent] = attackEquipPercentConfig
	propConfig[attackEquipLevelPercent] = attackEquipLevelPercentConfig

	modulePropConfig[rootModuleId] = map[int]int{attackAbsolute: attackRootPercent}
	modulePropConfig[equipModuleId] = map[int]int{attackAbsolute: attackEquipPercent}
	modulePropConfig[equipLevelModuleId] = map[int]int{attackAbsolute: attackEquipLevelPercent}
}

func TestProp(t *testing.T) {
	initConfig()
	tree := NewTreeWithProp(map[int]map[int]float64{
		equipLevelModuleId: {
			attackAbsolute:          100,
			attackEquipLevelPercent: 50,
			attackEquipPercent:      20,
			attackRootPercent:       10,
		},
	})
	// 100 * 1.5 * 1.2 * 1.1 = 198
	fmt.Printf("%+v\n", tree.RootNode.PropResult)

	// 120 * 1.5 * 1.2 * 1.1 = 396
	tree.ChangeModuleProp(equipLevelModuleId, map[int]float64{attackAbsolute: 200})
	fmt.Printf("%+v\n", tree.RootNode.PropResult)
}
