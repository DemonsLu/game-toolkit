package component

type IComponent interface {
	GetType() (t int)
}

type Component struct {
	CType int
}

func (c *Component) GetType() int {
	return c.CType
}
