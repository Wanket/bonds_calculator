package marshaling

import (
	"fmt"
	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
	"reflect"
)

var errUnsupportedType = fmt.Errorf("unsupported type")

func MarshalBaseType(t any, writer *jwriter.Writer) { //nolint:cyclop
	switch obj := t.(type) {
	case bool:
		writer.Bool(obj)
	case int:
		writer.Int(obj)
	case int8:
		writer.Int8(obj)
	case int16:
		writer.Int16(obj)
	case int32:
		writer.Int32(obj)
	case int64:
		writer.Int64(obj)
	case uint:
		writer.Uint(obj)
	case uint8:
		writer.Uint8(obj)
	case uint16:
		writer.Uint16(obj)
	case uint32:
		writer.Uint32(obj)
	case uint64:
		writer.Uint64(obj)
	case float32:
		writer.Float32(obj)
	case float64:
		writer.Float64(obj)
	case string:
		writer.String(obj)
	default:
		panic(fmt.Errorf("MarshalBaseType: %w %s", errUnsupportedType, reflect.TypeOf(obj)))
	}
}

func UnmarshalBaseType(lexer *jlexer.Lexer, out any) { //nolint:cyclop
	switch obj := out.(type) {
	case *bool:
		*obj = lexer.Bool()
	case *int:
		*obj = lexer.Int()
	case *int8:
		*obj = lexer.Int8()
	case *int16:
		*obj = lexer.Int16()
	case *int32:
		*obj = lexer.Int32()
	case *int64:
		*obj = lexer.Int64()
	case *uint:
		*obj = lexer.Uint()
	case *uint8:
		*obj = lexer.Uint8()
	case *uint16:
		*obj = lexer.Uint16()
	case *uint32:
		*obj = lexer.Uint32()
	case *uint64:
		*obj = lexer.Uint64()
	case *float32:
		*obj = lexer.Float32()
	case *float64:
		*obj = lexer.Float64()
	case *string:
		*obj = lexer.String()
	default:
		panic(fmt.Errorf("UnmarshalBaseType: %w %s", errUnsupportedType, reflect.TypeOf(obj)))
	}
}
