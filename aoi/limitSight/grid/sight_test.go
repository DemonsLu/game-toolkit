package grid

import "testing"

func TestEntity(t *testing.T) {
	e1 := NewEntity()
	e2 := NewEntity()
	e3 := NewEntity()
	e4 := NewEntity()

	e1.EnterMap(1, 1)
	e2.EnterMap(2, 2)
	e3.EnterMap(15, 10)
	e4.EnterMap(25, 10)

	e1.ChangePosition(35, 11)

	e1.LeaveMap()
	e2.LeaveMap()
	e3.LeaveMap()
	e4.LeaveMap()
}
