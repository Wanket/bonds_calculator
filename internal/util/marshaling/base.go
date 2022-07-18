package marshaling

import (
	"fmt"
	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
	"reflect"
)

func MarshalBaseType(t any, writer *jwriter.Writer) {
	switch t := t.(type) {
	case bool:
		writer.Bool(t)
	case int:
		writer.Int(t)
	case int8:
		writer.Int8(t)
	case int16:
		writer.Int16(t)
	case int32:
		writer.Int32(t)
	case int64:
		writer.Int64(t)
	case uint:
		writer.Uint(t)
	case uint8:
		writer.Uint8(t)
	case uint16:
		writer.Uint16(t)
	case uint32:
		writer.Uint32(t)
	case uint64:
		writer.Uint64(t)
	case float32:
		writer.Float32(t)
	case float64:
		writer.Float64(t)
	case string:
		writer.String(t)
	default:
		panic(fmt.Errorf("unsupported type: %s", reflect.TypeOf(t)))
	}
}

func UnmarshalBaseType(lexer *jlexer.Lexer, out any) {
	switch t := out.(type) {
	case *bool:
		*t = lexer.Bool()
	case *int:
		*t = lexer.Int()
	case *int8:
		*t = lexer.Int8()
	case *int16:
		*t = lexer.Int16()
	case *int32:
		*t = lexer.Int32()
	case *int64:
		*t = lexer.Int64()
	case *uint:
		*t = lexer.Uint()
	case *uint8:
		*t = lexer.Uint8()
	case *uint16:
		*t = lexer.Uint16()
	case *uint32:
		*t = lexer.Uint32()
	case *uint64:
		*t = lexer.Uint64()
	case *float32:
		*t = lexer.Float32()
	case *float64:
		*t = lexer.Float64()
	case *string:
		*t = lexer.String()
	default:
		panic(fmt.Errorf("unsupported type: %s", reflect.TypeOf(t)))
	}
}
