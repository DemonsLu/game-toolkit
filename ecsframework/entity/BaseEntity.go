package entity

import (
	"ecsframework/component"
)

type Entity struct {
	Components map[int]component.IComponent
}

type IEntity interface {
	RegisteredComponent(iComponent component.IComponent)
	GetComponent(typeId int) (c component.IComponent)
}

func (e *Entity) RegisteredComponent(iComponent component.IComponent) {
	if e.Components == nil {
		e.Components = make(map[int]component.IComponent)
	}
	if iComponent == nil {
		return
	}
	e.Components[iComponent.GetType()] = iComponent
}

func (e *Entity) GetComponent(typeId int) (c component.IComponent) {
	return e.Components[typeId]
}
