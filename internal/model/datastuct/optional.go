package datastuct

type Optional[T any] struct {
	val   T
	exist bool
}

func NewOptional[T any](val T) Optional[T] {
	return Optional[T]{
		val:   val,
		exist: true,
	}
}

func (optional *Optional[T]) Get() (T, bool) {
	return optional.val, optional.exist
}

func (optional *Optional[T]) Set(val T) {
	optional.val = val
	optional.exist = true
}

func (optional *Optional[T]) Empty() {
	var zero T
	optional.val = zero
	optional.exist = false
}
