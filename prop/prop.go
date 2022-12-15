package prop

type Tree struct {
	RootNode    *ModuleNode
	ModuleIdMap map[int]*ModuleNode // key: moduleId value: node
}

type ModuleNode struct {
	ModuleId int // 模块ID

	PropAbsolute map[int]float64 // key: propId value: 绝对值
	PropPercent  map[int]float64 // key: propId value: 百分比

	PropResult map[int]float64 // key: propId value: 结果值

	ParentModuleNode *ModuleNode         // 该节点的父模块集合
	ChildModuleNode  map[int]*ModuleNode // 该节点的子模块集合 key: 子节点ModuleId, value: 子节点
}

// 这三个全局变量应该在程序初始化的时候，通过读取配置的方式将其赋值
var moduleConfig = make(map[int]ModuleConfig)
var propConfig = make(map[int]PropConfig)
var modulePropConfig = make(map[int]map[int]int) // outer key: moduleId, inner key: propId(绝对值属性), value: 当前模块绝对值属性受影响的百分比属性
var rootModuleConfigId int

type ModuleConfig struct {
	ModuleId       int // 当前模块Id
	ParentModuleId int // 父模块Id (如果没有父模块ID，则为0)
}

type PropConfig struct {
	PropId int

	IsPercentage    bool
	RelativePropId  int // 受这个属性影响的属性ID，百分比属性需要
	ConcernModuleId int // 关心这个属性的模块ID，百分比属性需要
}

func NewTreeWithProp(modulePropValue map[int]map[int]float64) *Tree {
	root := &ModuleNode{
		ModuleId:         rootModuleConfigId,
		PropAbsolute:     make(map[int]float64),
		PropPercent:      make(map[int]float64),
		PropResult:       make(map[int]float64),
		ParentModuleNode: nil,
		ChildModuleNode:  make(map[int]*ModuleNode),
	}
	t := &Tree{
		RootNode: root,
		ModuleIdMap: map[int]*ModuleNode{
			rootModuleConfigId: root,
		},
	}
	for moduleId, propValue := range modulePropValue {
		t.ChangeModuleProp(moduleId, propValue)
	}
	return t
}

func (t *Tree) ChangeModuleProp(moduleId int, propValue map[int]float64) {
	t.BuildTreeEnsure(t.moduleChain(moduleId))
	n := t.ModuleIdMap[moduleId]
	// 调用者只能给叶子结点所在的moduleId发送AddModuleProp
	if n.ChildModuleNode != nil && len(n.ChildModuleNode) > 0 {
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
			if c.ConcernModuleId == n.ModuleId {
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
			latestResult = n.PropAbsolute[id] * (100 + n.PropPercent[modulePropConfig[n.ModuleId][id]]) / 100
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

func (t *Tree) BuildTreeEnsure(moduleChain []int) {
	for _, checkModuleId := range moduleChain {
		if _, ok := t.ModuleIdMap[checkModuleId]; ok {
			continue
		}

		p := t.ModuleIdMap[moduleConfig[checkModuleId].ParentModuleId]
		n := ModuleNode{
			ModuleId:         checkModuleId,
			PropAbsolute:     make(map[int]float64),
			PropPercent:      make(map[int]float64),
			PropResult:       make(map[int]float64),
			ParentModuleNode: p,
			ChildModuleNode:  make(map[int]*ModuleNode),
		}
		p.ChildModuleNode[checkModuleId] = &n
		t.ModuleIdMap[checkModuleId] = &n
	}
}

func (t *Tree) moduleChain(moduleId int) []int {
	result := make([]int, 0)
	for moduleId != 0 {
		result = append(result, moduleId)
		moduleId = moduleConfig[moduleId].ParentModuleId
	}
	return Reverse(result)
}

func Reverse(input []int) []int {
	var output []int

	for i := len(input) - 1; i >= 0; i-- {
		output = append(output, input[i])
	}

	return output
}
