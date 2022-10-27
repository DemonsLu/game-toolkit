package entity

import "ecsframework/component"

// you know, for example

type Role struct {
	Entity
}

func NewRole() *Role {
	r := &Role{Entity: Entity{}}
	r.RegisteredComponent(component.NewWingComponent())
	return r
}
