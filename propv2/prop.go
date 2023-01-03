package propv2

/*
	RootNode的ModuleId和Layer均为零
	RootNode下方是各个模块，它们的ModuleId可以自行配置，Layer为零
	各模块下方是子模块slice，它们按照深度来确定Layer，同一ModuleId，同一Layer的各个模块可以视为是相同的
*/

type Tree struct {
	RootNode    *ModuleNode
	ModuleIdMap map[int]*ModuleNode // key: moduleId value: node
}

type ModuleNode struct {
	ModuleId int // 模块ID
	Layer    int // 深度

	PropAbsolute map[int]float64 // key: propId value: 绝对值
	PropPercent  map[int]float64 // key: propId value: 百分比

	PropResult map[int]float64 // key: propId value: 结果值

	ParentModuleNode *ModuleNode   // 该节点的父模块集合
	ChildModuleNode  []*ModuleNode // 该节点的子模块集合 key: 子节点ModuleId, value: 子节点
}

var propConfig = make(map[int]PropConfig)
var modulePropConfig = make(map[int]map[int]map[int]int) // key0: moduleId, key1: layer, key2: propId(绝对值属性), value: 当前模块绝对值属性受影响的百分比属性

type PropConfig struct {
	PropId int

	IsPercentage    bool
	RelativePropId  int // 受这个属性影响的属性ID，百分比属性需要
	ConcernModuleId int // 关心这个属性的模块ID，百分比属性需要
	ConcernLayer    int // 关心这个属性的Layer，百分比属性需要
}

func NewTree() *Tree {
	root := &ModuleNode{
		ModuleId:         0,
		PropAbsolute:     make(map[int]float64),
		PropPercent:      make(map[int]float64),
		PropResult:       make(map[int]float64),
		ParentModuleNode: nil,
		ChildModuleNode:  make([]*ModuleNode, 0),
	}
	t := &Tree{
		RootNode:    root,
		ModuleIdMap: make(map[int]*ModuleNode),
	}
	return t
}

func (t *Tree) BuildByModule(moduleId int) (result *ModuleNode) {
	result = &ModuleNode{
		ModuleId:         moduleId,
		Layer:            0,
		PropAbsolute:     make(map[int]float64),
		PropPercent:      make(map[int]float64),
		PropResult:       make(map[int]float64),
		ParentModuleNode: t.RootNode,
		ChildModuleNode:  make([]*ModuleNode, 0),
	}
	t.ModuleIdMap[moduleId] = result
	t.RootNode.ChildModuleNode = append(t.RootNode.ChildModuleNode, result)
	return
}

func (t *Tree) BuildBySubmodule(n *ModuleNode) (result *ModuleNode) {
	result = &ModuleNode{
		ModuleId:         n.ModuleId,
		Layer:            n.Layer + 1,
		PropAbsolute:     make(map[int]float64),
		PropPercent:      make(map[int]float64),
		PropResult:       make(map[int]float64),
		ParentModuleNode: n,
		ChildModuleNode:  make([]*ModuleNode, 0),
	}
	n.ChildModuleNode = append(n.ChildModuleNode, result)
	return
}

func (t *Tree) ChangeModuleProp(n *ModuleNode, propValue map[int]float64) {
	if n == nil {
		return
	}
	t.calcProp(n, propValue)
}

func (t *Tree) calcProp(n *ModuleNode, propValue map[int]float64) {
	changedPropIds := make([]int, 0)
	// 若是叶子结点，则propValue表示属性值；若是非叶子结点，则propValue表示属性diff值
	for id, v := range propValue {
		if n.ChildModuleNode == nil || len(n.ChildModuleNode) == 0 {
			if propConfig[id].IsPercentage {
				n.PropPercent[id] = v
			} else {
				n.PropAbsolute[id] = v
			}
		} else {
			if propConfig[id].IsPercentage {
				n.PropPercent[id] += v
			} else {
				n.PropAbsolute[id] += v
			}
		}
		changedPropIds = append(changedPropIds, id)
	}

	propDiff := make(map[int]float64)
	for _, id := range changedPropIds {
		c := propConfig[id]
		var latestResult float64
		if c.IsPercentage {
			if c.ConcernModuleId == n.ModuleId && c.ConcernLayer == n.Layer {
				// 是本模块关心的百分比属性，那么就在这一级计算好最终result
				latestResult = n.PropAbsolute[c.RelativePropId] * (100 + n.PropPercent[id]) / 100
				if _, ok := propDiff[c.RelativePropId]; !ok {
					propDiff[c.RelativePropId] = latestResult - n.PropResult[c.RelativePropId]
				}
				n.PropResult[c.RelativePropId] = latestResult
			} else {
				// 不是本模块关心的百分比属性，说明是上级结点需要关心的，我们只需要将其直接上浮即可
				latestResult = n.PropPercent[id]
				if _, ok := propDiff[id]; !ok {
					propDiff[id] = latestResult - n.PropResult[id]
				}
				n.PropResult[id] = latestResult
			}
		} else {
			// 绝对值属性应该去寻找当前模块关心的对应百分比属性
			latestResult = n.PropAbsolute[id] * (100 + n.PropPercent[modulePropConfig[n.ModuleId][n.Layer][id]]) / 100
			if _, ok := propDiff[id]; !ok {
				propDiff[id] = latestResult - n.PropResult[id]
			}
			n.PropResult[id] = latestResult
		}
	}

	if n.ParentModuleNode != nil {
		t.calcProp(n.ParentModuleNode, propDiff)
	}
}
