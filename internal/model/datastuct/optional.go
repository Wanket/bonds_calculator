package datastuct

import (
	"bonds_calculator/internal/util/marshaling"
	"github.com/mailru/easyjson"
	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
)

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

func (optional *Optional[T]) MarshalEasyJSON(writer *jwriter.Writer) {
	if !optional.exist {
		writer.RawString("null")

		return
	}

	if marshaler, is := any(optional.val).(easyjson.Marshaler); is {
		marshaler.MarshalEasyJSON(writer)
	} else {
		marshaling.MarshalBaseType(optional.val, writer)
	}
}

func (optional *Optional[T]) UnmarshalEasyJSON(lexer *jlexer.Lexer) {
	if lexer.IsNull() {
		optional.Empty()

		lexer.Skip()

		return
	}

	if unmarshaler, is := any(optional.val).(easyjson.Unmarshaler); is {
		unmarshaler.UnmarshalEasyJSON(lexer)
	} else {
		marshaling.UnmarshalBaseType(lexer, optional.val)
	}
}
