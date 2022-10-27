package component

import "ecsframework/consts"

// you know, for example

type CWing struct {
	Component
	WingLevel int
	WingStage int
}

func NewWingComponent() *CWing {
	return &CWing{
		Component: Component{CType: consts.CSTypeWings},
	}
}
