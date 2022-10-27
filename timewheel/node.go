package main

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type NodeManager struct {
	header *Node
	tail   *Node
}

func NewNodeManager() *NodeManager {
	return &NodeManager{
		header: nil,
		tail:   nil,
	}
}

type Node struct {
	Id   uint64 // 用来比较大小
	Data *Timer
	Next *Node
}

func (m *NodeManager) Add(data *Timer) {
	// 判断SkipList初始化状态
	if m == nil || data == nil {
		return
	}
	id := data.Id
	dstNode := &Node{
		Id:   id,
		Data: data,
		Next: nil,
	}
	// 一个元素都没有
	if m.header == nil && m.tail == nil {
		m.header = dstNode
	} else {
		node := m.tail
		node.Next = dstNode
	}
	m.tail = dstNode
}

func (m *NodeManager) GetAll() (results []*Timer) {
	results = make([]*Timer, 0)
	if m == nil || m.header == nil {
		return
	}
	for node := m.header; node != nil; node = node.Next {
		results = append(results, node.Data)
	}
	return
}

func (m *NodeManager) PopAll() (results []*Timer) {
	results = m.GetAll()
	if m == nil || m.header == nil {
		return
	}
	m.header = nil
	m.tail = nil

	return
}
