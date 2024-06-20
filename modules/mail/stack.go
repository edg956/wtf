package mail

import (
	"errors"
	"sync"
)

type stack struct {
	lock sync.Mutex // you don't have to do this if you don't want thread safety
	s    []interface{}
}

func newStack() *stack {
	return &stack{sync.Mutex{}, make([]interface{}, 0)}
}

func (s *stack) push(v interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.s = append(s.s, v)
}

func (s *stack) pop() (interface{}, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	l := len(s.s)
	if l == 0 {
		return 0, errors.New("Empty Stack")
	}

	res := s.s[l-1]
	s.s = s.s[:l-1]
	return res, nil
}

func (s *stack) peek() (interface{}, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	l := len(s.s)
	if l == 0 {
		return 0, errors.New("Empty Stack")
	}

	return s.s[l-1], nil
}

func (s *stack) size() int {
	s.lock.Lock()
	defer s.lock.Unlock()

	return len(s.s)
}
