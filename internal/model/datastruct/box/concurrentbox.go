package box

import "sync"

type ConcurrentBox[T any] struct {
	value T

	lock sync.RWMutex
}

func NewConcurrentBox[T any](value T) ConcurrentBox[T] {
	return ConcurrentBox[T]{
		value: value,
	}
}

func (box *ConcurrentBox[T]) SafeRead() T { //nolint:nolintlint,ireturn
	box.lock.RLock()

	value := box.value

	box.lock.RUnlock()

	return value
}

func (box *ConcurrentBox[T]) Set(value T) {
	box.lock.Lock()

	box.value = value

	box.lock.Unlock()
}
