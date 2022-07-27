package box

import "time"

type ConcurrentCacheBox[T any] struct {
	ConcurrentBox[T]

	changedTime time.Time
}

func NewConcurrentCacheBox[T any](value T) ConcurrentCacheBox[T] {
	return ConcurrentCacheBox[T]{
		ConcurrentBox: NewConcurrentBox(value),

		changedTime: time.Now(),
	}
}

func (box *ConcurrentCacheBox[T]) Set(value T) {
	box.lock.Lock()

	box.value = value
	box.changedTime = time.Now()

	box.lock.Unlock()
}

func (box *ConcurrentCacheBox[T]) SafeReadWithTime() (T, time.Time) { //nolint:ireturn
	result := box.ConcurrentBox.SafeRead()

	return result, box.changedTime
}
