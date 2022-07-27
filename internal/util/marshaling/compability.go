//nolint:wrapcheck
package marshaling

import "github.com/mailru/easyjson"

func Unmarshal(data []byte, obj interface{}) error {
	if unmarshaler, is := obj.(easyjson.Unmarshaler); is {
		return easyjson.Unmarshal(data, unmarshaler)
	}

	panic("object is not easyjson.Unmarshaler")
}

func Marshal(obj interface{}) ([]byte, error) {
	if marshaler, is := obj.(easyjson.Marshaler); is {
		return easyjson.Marshal(marshaler)
	}

	panic("object is not easyjson.Marshaler")
}
