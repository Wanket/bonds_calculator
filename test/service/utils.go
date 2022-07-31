package service

import (
	"fmt"
	"github.com/mailru/easyjson"
	testdataloader "github.com/peteole/testdata-loader"
)

var errTypeIsNotMarshalerUnmarshaler = fmt.Errorf("type is not a easyjson.MarshalerUnmarshaler")

func LoadExpectedJSON[T any](path string) (string, error) {
	file := testdataloader.GetTestFile(path)

	var result T

	marshalerUnmarshaler, is := any(&result).(easyjson.MarshalerUnmarshaler)
	if !is {
		panic(fmt.Errorf("LoadExpectedJSON: %w, type: %T", errTypeIsNotMarshalerUnmarshaler, result))
	}

	if err := easyjson.Unmarshal(file, marshalerUnmarshaler); err != nil {
		panic(fmt.Errorf("cannot unmarshal type %T: %w", result, err))
	}

	stringResult, err := easyjson.Marshal(marshalerUnmarshaler)
	if err != nil {
		panic(fmt.Errorf("cannot marshal type %T: %w", result, err))
	}

	return string(stringResult), nil
}
