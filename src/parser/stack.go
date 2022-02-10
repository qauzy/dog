package parser

import (
	"dog/ast"
	"sync"
)

type Stack struct {
	data     []ast.Container
	length   int
	capacity int
	sync.Mutex
}

// 构建一个空栈
func InitStack() *Stack {
	return &Stack{data: make([]ast.Container, 200), length: 0, capacity: 200}
}

// 压栈操作
func (s *Stack) Push(data ast.Container) {
	s.Lock()
	defer s.Unlock()

	if s.length+1 >= s.capacity {
		s.capacity <<= 1
		t := s.data
		s.data = make([]ast.Container, s.capacity)
		copy(s.data, t)
	}

	s.data[s.length] = data
	s.length++
}

// 出栈操作
func (s *Stack) Pop() ast.Container {
	s.Lock()
	defer s.Unlock()

	if s.length <= 0 {
		panic("int stack pop: index out of range")
	}

	t := s.data[s.length-1]
	s.length--

	return t
}

// 返回栈顶元素
func (s *Stack) Peek() ast.Container {
	s.Lock()
	defer s.Unlock()

	if s.length <= 0 {
		panic("empty stack")
	}

	return s.data[s.length-1]
}
func (s *Stack) List() []ast.Container {
	s.Lock()
	defer s.Unlock()

	if s.length <= 0 {
		return nil
	}

	return s.data[:s.length]
}

// 返回当前栈元素个数
func (s *Stack) Count() int {
	s.Lock()
	defer s.Unlock()

	t := s.length

	return t
}

// 清空栈
func (s *Stack) Clear() {
	s.Lock()
	defer s.Unlock()

	s.data = make([]ast.Container, 8)
	s.length = 0
	s.capacity = 8
}

// 栈是否为空
func (s *Stack) IsEmpty() bool {
	s.Lock()
	defer s.Unlock()
	b := s.length == 0
	return b
}
