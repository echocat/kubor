package functions

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

var FuncToString = Function{
	Description: `Converts the given argument to a string.`,
	Parameters: Parameters{{
		Name: "in",
	}},
}.MustWithFunc(func(in ...interface{}) string {
	return fmt.Sprint(in)
})

var FuncParseInt = Function{
	Description: "Interprets a given string <str> as an int and returns the result. If this is not a valid number it will fail.",
	Parameters: Parameters{{
		Name: "str",
	}},
}.MustWithFunc(parseInt)

var FuncParseInt64 = Function{
	Description: "Interprets a given string <str> as a int64 and returns the result. If this is not a valid number it will fail.",
	Parameters: Parameters{{
		Name: "str",
	}},
}.MustWithFunc(parseInt64)

var FuncParseFloat = Function{
	Description: "Interprets a given string <str> as a float and returns the result. If this is not a valid number it will fail.",
	Parameters: Parameters{{
		Name: "str",
	}},
}.MustWithFunc(parseFloat32)

var FuncParseDouble = Function{
	Description: "Interprets a given string <str> as a float64 and returns the result. If this is not a valid number it will fail.",
	Parameters: Parameters{{
		Name: "str",
	}},
}.MustWithFunc(parseFloat64)

var FuncIsInt = Function{
	Description: "Will return <true> if the given string <str> is a valid int.",
	Parameters: Parameters{{
		Name: "str",
	}},
}.MustWithFunc(isInt)

var FuncIsInt64 = Function{
	Description: "Will return <true> if the given string <str> is a valid int64.",
	Parameters: Parameters{{
		Name: "str",
	}},
}.MustWithFunc(isInt)

var FuncIsFloat = Function{
	Description: "Will return <true> if the given string <str> is a valid float.",
	Parameters: Parameters{{
		Name: "str",
	}},
}.MustWithFunc(isFloat)

var FuncIsDouble = Function{
	Description: "Will return <true> if the given string <str> is a valid float64.",
	Parameters: Parameters{{
		Name: "str",
	}},
}.MustWithFunc(isFloat)

var FuncToInt = Function{
	Description: "Will try to interpret the given <value> as an int. If this is not possible <0> is returned.",
	Parameters: Parameters{{
		Name: "value",
	}},
}.MustWithFunc(toInt)

var FuncToInt64 = Function{
	Description: "Will try to interpret the given <value> as a int64. If this is not possible <0> is returned.",
	Parameters: Parameters{{
		Name: "value",
	}},
}.MustWithFunc(toInt64)

var FuncToFloat = Function{
	Description: "Will try to interpret the given <value> as a float. If this is not possible <0.0> is returned.",
	Parameters: Parameters{{
		Name: "value",
	}},
}.MustWithFunc(toFloat32)

var FuncToDouble = Function{
	Description: "Will try to interpret the given <value> as a float64. If this is not possible <0.0> is returned.",
	Parameters: Parameters{{
		Name: "value",
	}},
}.MustWithFunc(toFloat64)

var FuncToBool = Function{
	Description: "Will try to interpret the given <value> as a bool. If this is not possible <false> is returned.",
	Parameters: Parameters{{
		Name: "value",
	}},
}.MustWithFunc(toBool)

var FuncsConversations = Functions{
	"toString":     FuncToString,
	"parseInt":     FuncParseInt,
	"parseInt64":   FuncParseInt64,
	"parseFloat32": FuncParseFloat,
	"parseFloat64": FuncParseDouble,
	"isInt":        FuncIsInt,
	"isInt64":      FuncIsInt64,
	"isFloat32":    FuncIsFloat,
	"isFloat64":    FuncIsDouble,
	"toInt":        FuncToInt,
	"toInt64":      FuncToInt64,
	"toFloat32":    FuncToFloat,
	"toFloat64":    FuncToDouble,
	"toBool":       FuncToBool,
}
var CategoryConversations = Category{
	Functions: FuncsConversations,
}

func isFloat(v string) bool {
	_, err := strconv.ParseFloat(v, 0)
	return err == nil
}

func parseFloat32(v string) (float32, error) {
	i, err := parseFloat64(v)
	return float32(i), err
}

func parseFloat64(v string) (float64, error) {
	return strconv.ParseFloat(v, 64)
}

func toFloat32(v interface{}) float32 {
	return float32(toFloat64(v))
}

func toFloat64(v interface{}) float64 {
	if str, ok := v.(string); ok {
		iv, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return 0
		}
		return iv
	}

	val := reflect.Indirect(reflect.ValueOf(v))
	switch val.Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return float64(val.Int())
	case reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return float64(val.Uint())
	case reflect.Uint, reflect.Uint64:
		return float64(val.Uint())
	case reflect.Float32, reflect.Float64:
		return val.Float()
	case reflect.Bool:
		if val.Bool() == true {
			return 1
		}
		return 0
	default:
		return 0
	}
}

func isInt(v string) bool {
	_, err := strconv.ParseInt(v, 10, 64)
	return err == nil
}

func parseInt(v string) (int, error) {
	i, err := parseFloat64(v)
	return int(i), err
}

func parseInt64(v string) (int64, error) {
	return strconv.ParseInt(v, 10, 64)
}

func toInt(v interface{}) int {
	return int(toInt64(v))
}

func toInt64(v interface{}) int64 {
	if str, ok := v.(string); ok {
		iv, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return 0
		}
		return iv
	}

	val := reflect.Indirect(reflect.ValueOf(v))
	switch val.Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return val.Int()
	case reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return int64(val.Uint())
	case reflect.Uint, reflect.Uint64:
		tv := val.Uint()
		if tv <= math.MaxInt64 {
			return int64(tv)
		}
		return math.MaxInt64
	case reflect.Float32, reflect.Float64:
		return int64(val.Float())
	case reflect.Bool:
		if val.Bool() == true {
			return 1
		}
		return 0
	default:
		return 0
	}
}

func toBool(v interface{}) bool {
	if b, ok := v.(bool); ok {
		return b
	}

	if str, ok := v.(string); ok {
		switch strings.TrimSpace(strings.ToLower(str)) {
		case "true":
			return true
		case "on":
			return true
		case "yes":
			return true
		default:
			return false
		}
	}

	val := reflect.Indirect(reflect.ValueOf(v))
	switch val.Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return val.Int() > 1
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint, reflect.Uint64:
		return val.Uint() > 1
	case reflect.Float32, reflect.Float64:
		return val.Float() > 1
	default:
		return false
	}
}
