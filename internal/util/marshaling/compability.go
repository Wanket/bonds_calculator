//nolint:wrapcheck
package marshaling

import (
	"fmt"
	"github.com/mailru/easyjson"
)

var (
	errTypeIsNotUnmarshaler = fmt.Errorf("type is not a easyjson.Unmarshaler")
	errTypeIsNotMarshaler   = fmt.Errorf("type is not a easyjson.Marshaler")
)

func Unmarshal(data []byte, obj interface{}) error {
	if unmarshaler, is := obj.(easyjson.Unmarshaler); is {
		return easyjson.Unmarshal(data, unmarshaler)
	}

	panic(fmt.Errorf("unmarshal: %w, type: %T", errTypeIsNotUnmarshaler, obj))
}

func Marshal(obj interface{}) ([]byte, error) {
	if marshaler, is := obj.(easyjson.Marshaler); is {
		return easyjson.Marshal(marshaler)
	}

	panic(fmt.Errorf("marshal: %w, type: %T", errTypeIsNotMarshaler, obj))
}
