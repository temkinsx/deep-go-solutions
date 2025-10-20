package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Task struct {
	Identifier int
	Priority   int
}

type Scheduler struct {
	heap []*Task
}

func NewScheduler() Scheduler {
	return Scheduler{heap: []*Task{}}
}

func (s *Scheduler) AddTask(task Task) {
	s.heap = append(s.heap, &task)
	i := len(s.heap) - 1
	for i > 0 {
		p := (i - 1) / 2
		if s.heap[p].Priority >= s.heap[i].Priority {
			break
		}
		s.heap[p], s.heap[i] = s.heap[i], s.heap[p]
		i = p
	}
}

func (s *Scheduler) GetTask() Task {
	n := len(s.heap)
	if n == 0 {
		return Task{}
	}
	if n == 1 {
		res := s.heap[0]
		s.heap = s.heap[:0]
		return *res
	}

	res := s.heap[0]
	last := s.heap[n-1]
	s.heap = s.heap[:n-1]
	s.heap[0] = last
	s.siftDown(0)
	return *res
}

func (s *Scheduler) siftDown(i int) {
	n := len(s.heap)
	for {
		l := 2*i + 1
		r := l + 1
		largest := i
		if l < n && s.heap[l].Priority > s.heap[largest].Priority {
			largest = l
		}
		if r < n && s.heap[r].Priority > s.heap[largest].Priority {
			largest = r
		}
		if largest == i {
			break
		}
		s.heap[i], s.heap[largest] = s.heap[largest], s.heap[i]
		i = largest
	}
}

func (s *Scheduler) siftUp(i int) {
	for i > 0 {
		p := (i - 1) / 2
		if s.heap[p].Priority >= s.heap[i].Priority {
			break
		}
		s.heap[p], s.heap[i] = s.heap[i], s.heap[p]
		i = p
	}
}

func (s *Scheduler) ChangeTaskPriority(taskID int, newPriority int) {
	idx := -1
	for i := 0; i < len(s.heap); i++ {
		if s.heap[i].Identifier == taskID {
			idx = i
			break
		}
	}
	if idx == -1 {
		return
	}

	old := s.heap[idx].Priority
	s.heap[idx].Priority = newPriority

	if newPriority > old {
		s.siftUp(idx)
	} else if newPriority < old {
		s.siftDown(idx)
	}
}

func TestTrace(t *testing.T) {
	task1 := Task{Identifier: 1, Priority: 10}
	task2 := Task{Identifier: 2, Priority: 20}
	task3 := Task{Identifier: 3, Priority: 30}
	task4 := Task{Identifier: 4, Priority: 40}
	task5 := Task{Identifier: 5, Priority: 50}

	scheduler := NewScheduler()
	scheduler.AddTask(task1)
	scheduler.AddTask(task2)
	scheduler.AddTask(task3)
	scheduler.AddTask(task4)
	scheduler.AddTask(task5)

	task := scheduler.GetTask()
	assert.Equal(t, task5, task)

	task = scheduler.GetTask()
	assert.Equal(t, task4, task)

	scheduler.ChangeTaskPriority(1, 100)

	task = scheduler.GetTask()
	assert.Equal(t, Task{Identifier: 1, Priority: 100}, task)

	task = scheduler.GetTask()
	assert.Equal(t, task3, task)
}
