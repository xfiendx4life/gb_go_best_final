package sqlparser

import (
	"fmt"
)

type Stack struct {
	data []string
}

func (s *Stack) IsEmpty() bool {
	return len(s.data) == 0
}

func (s *Stack) Pop() (string, error) {
	if s.IsEmpty() {
		return "", fmt.Errorf("can't pop from empty stack")
	}
	res := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return res, nil
}

func NewStack(length int) *Stack {
	return &Stack{
		data: make([]string, length),
	}
}

func (s *Stack) Push(d string) {
	s.data = append(s.data, d)
}

func (s *Stack) Len() int {
	return len(s.data)
}
