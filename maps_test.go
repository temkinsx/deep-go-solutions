package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v maps_test.go
type Node struct {
	Key   int
	Value int
	Left  *Node
	Right *Node
}
type OrderedMap struct {
	Root *Node
	size int
}

func NewOrderedMap() *OrderedMap {
	return &OrderedMap{
		Root: nil,
		size: 0,
	}
}

func (m *OrderedMap) Insert(key, value int) {
	x := m.Root
	var y *Node
	for x != nil {
		if key == x.Key {
			x.Value = value
			return
		} else {
			y = x
			if key < x.Key {
				x = x.Left
			} else {
				x = x.Right
			}
		}
	}

	newNode := &Node{Key: key, Value: value}
	if y == nil {
		m.Root = newNode
	} else {
		if key < y.Key {
			y.Left = newNode
		} else {
			y.Right = newNode
		}
	}
	m.size++
}

func (m *OrderedMap) Erase(key int) {
	if m.Root == nil {
		return
	}

	// 1) Найти узел x с ключом key и его родителя parent
	var parent *Node
	x := m.Root
	for x != nil && x.Key != key {
		parent = x
		if key < x.Key {
			x = x.Left
		} else {
			x = x.Right
		}
	}
	if x == nil {
		return // не нашли
	}

	// 2) Если у x два ребёнка – заменить значениями с преемником и удалить преемника
	if x.Left != nil && x.Right != nil {
		succParent := x
		succ := x.Right
		for succ.Left != nil {
			succParent = succ
			succ = succ.Left
		}
		// Копируем ключ/значение преемника в x
		x.Key, x.Value = succ.Key, succ.Value
		// Теперь будем удалять succ как узел с <=1 ребёнком
		x, parent = succ, succParent
	}

	// 3) У x 0 или 1 ребёнок
	var child *Node
	if x.Left != nil {
		child = x.Left
	} else {
		child = x.Right
	}

	// 4) Переподвесить
	if parent == nil {
		// удаляем корень
		m.Root = child
	} else if parent.Left == x {
		parent.Left = child
	} else {
		parent.Right = child
	}

	m.size--
}

func (m *OrderedMap) Contains(key int) bool {
	x := m.Root
	for x != nil {
		if key == x.Key {
			return true
		} else {
			if key < x.Key {
				x = x.Left
			} else {
				x = x.Right
			}
		}
	}
	return false
}

func (m *OrderedMap) Size() int {
	return m.size
}

func (m *OrderedMap) ForEach(action func(int, int)) {
	var traverse func(n *Node)
	traverse = func(n *Node) {
		if n == nil {
			return
		}
		traverse(n.Left)
		action(n.Key, n.Value)
		traverse(n.Right)
	}
	traverse(m.Root)
}

func TestMap1(t *testing.T) {
	data := NewOrderedMap()
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}
