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

func (box *ConcurrentBox[T]) LockAndRead() T {
	box.lock.RLock()

	return box.value
}

func (box *ConcurrentBox[T]) UnlockRead() {
	box.lock.RUnlock()
}

func (box *ConcurrentBox[T]) Set(value T) {
	box.lock.Lock()

	box.value = value

	box.lock.Unlock()
}
