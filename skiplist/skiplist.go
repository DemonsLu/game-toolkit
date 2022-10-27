package main

import (
	"math"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type SkipList struct {
	maxDepth int
	header   *Node
}

func NewSkipList(maxDepth int) *SkipList {
	if maxDepth < 5 {
		maxDepth = 5
	}

	header := &Node{
		Id:       -1,
		Data:     nil,
		Forward:  make([]*Node, maxDepth, maxDepth),
		Previous: make([]*Node, maxDepth, maxDepth),
	}

	return &SkipList{maxDepth, header}
}

type Node struct {
	Id       int64 // 用来比较大小
	Data     any
	Forward  []*Node // 后继
	Previous []*Node // 前驱
}

func (s *SkipList) findNode(id int64) *Node {
	if s == nil || s.header == nil || len(s.header.Forward) <= 0 {
		return nil
	}
	node := s.header
	// 从最上层往下层找
	// 看上去是O(n ^ 2)，实际上由于跳表的概率性质，这里是 O(log N)
	for i := len(s.header.Forward) - 1; i >= 0; i-- {
		n := node.Forward[i]
		for {
			// 说明本层已经找完了，开始在下层找
			if n == nil || n.Id > id {
				break
			}
			if n.Id == id {
				return n
			}
			// 本层还没找完，继续找本层的后一个元素
			node = n
			n = n.Forward[i]
		}
	}
	return nil
}

func (s *SkipList) Find(id int64) (data any) {
	node := s.findNode(id)
	if node == nil {
		return nil
	}
	return node.Data
}

func (s *SkipList) Add(id int64, data any) {
	// 判断SkipList初始化状态
	if s == nil || s.header == nil || len(s.header.Forward) <= 0 || data == nil {
		return
	}
	node := s.header

	recordNodes := make([]*Node, s.maxDepth)
	// 从最上层往下层找
	// 看上去是O(n ^ 2)，实际上由于跳表的概率性质，这里是 O(log N)
	for i := s.maxDepth - 1; i >= 0; i-- {
		n := node.Forward[i]
		for {
			// 说明本层已经找完了，开始在下层找
			if n == nil || n.Id > id {
				// node此时就是待插节点的前一个节点
				recordNodes[i] = node
				break
			}
			// 本层还没找完，继续找本层的后一个元素
			node = n
			n = n.Forward[i]
		}
	}

	// 开始插入，先判断这个新元素需要几层depth
	// 算法:先求幂，将结果和随机系数相乘后再求底，这样得到结果大的底的概率高。最后用maxDepth - 底，得到最终depth
	// depth ∈ [1, maxDepth]
	var depth = s.maxDepth - int(math.Log2(1+(rand.Float64()*(math.Pow(2, float64(s.maxDepth))))))
	if depth <= 0 {
		depth = 1
	}
	dstNode := &Node{
		Id:       id,
		Data:     data,
		Forward:  make([]*Node, depth, depth),
		Previous: make([]*Node, depth, depth),
	}
	for i := 0; i <= depth-1; i++ {
		dstNode.Forward[i] = recordNodes[i].Forward[i]
		dstNode.Previous[i] = recordNodes[i]
		recordNodes[i].Forward[i] = dstNode
		// 前驱肯定有 (因为有Header)，但后继可能为nil
		if dstNode.Forward[i] != nil {
			dstNode.Forward[i].Previous[i] = dstNode
		}
	}
}

func (s *SkipList) Pop(id int64) (data any) {
	node := s.findNode(id)
	if node == nil {
		return
	}
	depth := len(node.Forward)
	//fmt.Println("node depth", depth)
	for i := 0; i < depth; i++ {
		previous := node.Previous[i]
		previous.Forward[i] = node.Forward[i]

		forward := node.Forward[i]
		// 前驱肯定有 (因为有Header)，但后继可能为nil
		if forward != nil {
			forward.Previous[i] = node.Previous[i]
		}

		node.Previous[i], node.Forward[i] = nil, nil
	}
	return node.Data
}

func (s *SkipList) GetAll() (results []any) {
	results = make([]any, 0)
	if s == nil || s.header == nil || len(s.header.Forward) <= 0 {
		return
	}
	for node := s.header.Forward[0]; node != nil; node = node.Forward[0] {
		results = append(results, node.Data)
	}
	return
}

func (s *SkipList) PopAll() (results []any) {
	results = s.GetAll()
	if s == nil || s.header == nil || len(s.header.Forward) <= 0 {
		return
	}
	s.header = &Node{
		Id:       -1,
		Data:     nil,
		Forward:  make([]*Node, s.maxDepth, s.maxDepth),
		Previous: make([]*Node, s.maxDepth, s.maxDepth),
	}

	return
}
