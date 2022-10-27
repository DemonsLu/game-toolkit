package main

import "testing"

func TestNewSkipList(t *testing.T) {
	skipList := NewSkipList(5)
	skipList.Add(1, 1)
	skipList.Add(2, 2)
	skipList.Add(3, 3)
	t.Log(skipList.GetAll())

	skipList.Pop(2)
	t.Log(skipList.GetAll())

	skipList.Add(10, 10)
	t.Log(skipList.GetAll())

	skipList.PopAll()
	t.Log(skipList.GetAll())
}
