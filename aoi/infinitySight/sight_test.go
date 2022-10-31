package infinitySight

import (
	"testing"
)

func TestEntity(t *testing.T) {
	e1 := Entity{Id: 1}
	e1.EnterMap()

	e2 := Entity{Id: 2}
	e2.EnterMap()

	e3 := Entity{Id: 3}
	e3.EnterMap()

	e1.ChangePosition(1, 2)

	e3.LeaveMap()
	e2.LeaveMap()
	e1.LeaveMap()
}
