package main

import "fmt"

type Node struct {
	Key   int
	Value int
	Next  *Node
}

type MyHashMap struct {
	buckets []*Node
	size    int
}

func Constructor() MyHashMap {
	size := 1000
	return MyHashMap{
		buckets: make([]*Node, size),
		size:    size,
	}
}

func (m *MyHashMap) hash(key int) int {
	return key % m.size
}

func (m *MyHashMap) Put(key int, value int) {
	index := m.hash(key)

	if m.buckets[index] == nil {
		m.buckets[index] = &Node{
			Key:   key,
			Value: value,
		}
	} else {
		node := m.buckets[index]
		for node.Next != nil {
			if node.Key == key {
				node.Value = value

				return
			}

			node = node.Next
		}

		node.Next = &Node{
			Key:   key,
			Value: value,
		}
	}
}

func (m *MyHashMap) Get(key int) int {
	index := m.hash(key)
	current := m.buckets[index]

	for current != nil {
		if current.Key == key {
			return current.Value
		}
		current = current.Next
	}

	return -1
}

func (m *MyHashMap) Remove(key int) {
	index := m.hash(key)

	if m.buckets[index] != nil {
		node := m.buckets[index]
		if node.Key == key {
			node = node.Next

			return
		}

		prev := node
		curr := node.Next
		for curr != nil {
			if curr.Key == key {
				prev.Next = curr.Next

				return
			}

			prev = curr
			curr = curr.Next
		}
	}
}

func main() {
	obj := Constructor()
	obj.Put(1, 2)
	fmt.Println(obj)
	param_2 := obj.Get(2)
	fmt.Println(param_2)
	obj.Remove(2)
	fmt.Println(obj)
}
