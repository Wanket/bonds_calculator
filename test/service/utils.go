package service

import (
	"fmt"
	"github.com/mailru/easyjson"
	testdataloader "github.com/peteole/testdata-loader"
)

func LoadExpectedJSON[T any](path string) (string, error) {
	file := testdataloader.GetTestFile(path)

	var result T

	marshalerUnmarshaler, is := any(&result).(easyjson.MarshalerUnmarshaler)
	if !is {
		panic("T is not a easyjson.MarshalerUnmarshaler")
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
