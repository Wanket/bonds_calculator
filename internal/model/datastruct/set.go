package datastruct

type Set[T comparable] struct {
	items map[T]struct{}
}

func NewSet[T comparable](capacity int) Set[T] {
	return Set[T]{
		items: make(map[T]struct{}, capacity),
	}
}

func (set *Set[T]) Add(item T) {
	set.items[item] = struct{}{}
}

func (set *Set[T]) Remove(item T) {
	delete(set.items, item)
}

func (set *Set[T]) Contains(item T) bool {
	_, ok := set.items[item]

	return ok
}

func (set *Set[T]) Size() int {
	return len(set.items)
}

func (set *Set[T]) Range() map[T]struct{} {
	return set.items
}
