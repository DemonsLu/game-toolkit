package crosschain

import "testing"

func TestSight(t *testing.T) {
	e1 := NewEntity()
	e2 := NewEntity()
	e3 := NewEntity()

	e1.EnterMap(100, 0, 0)
	t.Log("-------------------")
	e2.EnterMap(100, 1, 1)
	t.Log("-------------------")
	e3.EnterMap(500, 300, 200)
	t.Log("-------------------")
	e2.ChangePosition(350, 250)
	t.Log("-------------------")
	e3.LeaveMap()
}
