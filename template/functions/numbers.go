package functions

import (
	"math"
	"reflect"
	"strconv"
	"strings"
)

var _ = Register(Function{
	Name:        "parseInteger",
	Category:    "numbers",
	Description: "Interprets a given string <str> as an integer and returns the result. If this is not a valid number it will fail.",
	Parameters: Parameters{{
		Name: "str",
	}},
	Func: parseInt,
}, Function{
	Name:        "parseLong",
	Category:    "numbers",
	Description: "Interprets a given string <str> as a long and returns the result. If this is not a valid number it will fail.",
	Parameters: Parameters{{
		Name: "str",
	}},
	Func: parseInt64,
}, Function{
	Name:        "parseFloat",
	Category:    "numbers",
	Description: "Interprets a given string <str> as a float and returns the result. If this is not a valid number it will fail.",
	Parameters: Parameters{{
		Name: "str",
	}},
	Func: parseFloat32,
}, Function{
	Name:        "parseDouble",
	Category:    "numbers",
	Description: "Interprets a given string <str> as a double and returns the result. If this is not a valid number it will fail.",
	Parameters: Parameters{{
		Name: "str",
	}},
	Func: parseFloat64,
}, Function{
	Name:        "isInteger",
	Category:    "numbers",
	Description: "Will return <true> if the given string <str> is a valid integer.",
	Parameters: Parameters{{
		Name: "str",
	}},
	Func: isInt,
}, Function{
	Name:        "isLong",
	Category:    "numbers",
	Description: "Will return <true> if the given string <str> is a valid long.",
	Parameters: Parameters{{
		Name: "str",
	}},
	Func: isInt,
}, Function{
	Name:        "isFloat",
	Category:    "numbers",
	Description: "Will return <true> if the given string <str> is a valid float.",
	Parameters: Parameters{{
		Name: "str",
	}},
	Func: isFloat,
}, Function{
	Name:        "isDouble",
	Category:    "numbers",
	Description: "Will return <true> if the given string <str> is a valid double.",
	Parameters: Parameters{{
		Name: "str",
	}},
	Func: isFloat,
}, Function{
	Name:        "toInteger",
	Category:    "numbers",
	Description: "Will try to interpret the given <value> as an integer. If this is not possible <0> is returned.",
	Parameters: Parameters{{
		Name: "value",
	}},
	Func: toInt,
}, Function{
	Name:        "toLong",
	Category:    "numbers",
	Description: "Will try to interpret the given <value> as a long. If this is not possible <0> is returned.",
	Parameters: Parameters{{
		Name: "value",
	}},
	Func: toInt64,
}, Function{
	Name:        "toFloat",
	Category:    "numbers",
	Description: "Will try to interpret the given <value> as a float. If this is not possible <0.0> is returned.",
	Parameters: Parameters{{
		Name: "value",
	}},
	Func: toFloat32,
}, Function{
	Name:        "toDouble",
	Category:    "numbers",
	Description: "Will try to interpret the given <value> as a double. If this is not possible <0.0> is returned.",
	Parameters: Parameters{{
		Name: "value",
	}},
	Func: toFloat64,
}, Function{
	Name:        "toBool",
	Category:    "numbers",
	Description: "Will try to interpret the given <value> as a bool. If this is not possible <false> is returned.",
	Parameters: Parameters{{
		Name: "value",
	}},
	Func: toBool,
}, Function{
	Name:        "toBool",
	Category:    "numbers",
	Description: "Will try to interpret the given <value> as a bool. If this is not possible <false> is returned.",
	Parameters: Parameters{{
		Name: "value",
	}},
	Func: toBool,
}, Function{
	Name:        "minInteger",
	Category:    "numbers",
	Description: "Will pick the smallest integer of the <left> and <right>.",
	Parameters: Parameters{{
		Name: "left",
	}, {
		Name: "right",
	}},
	Func: minInt,
}, Function{
	Name:        "minLong",
	Category:    "numbers",
	Description: "Will pick the smallest long of the <left> and <right>.",
	Parameters: Parameters{{
		Name: "left",
	}, {
		Name: "right",
	}},
	Func: minInt64,
}, Function{
	Name:        "maxInteger",
	Category:    "numbers",
	Description: "Will pick the biggest integer of the <left> and <right>.",
	Parameters: Parameters{{
		Name: "left",
	}, {
		Name: "right",
	}},
	Func: maxInt,
}, Function{
	Name:        "maxLong",
	Category:    "numbers",
	Description: "Will pick the biggest long of the <left> and <right>.",
	Parameters: Parameters{{
		Name: "left",
	}, {
		Name: "right",
	}},
	Func: maxInt64,
}, Function{
	Name:        "minFloat",
	Category:    "numbers",
	Description: "Will pick the smallest float of the <left> and <right>.",
	Parameters: Parameters{{
		Name: "left",
	}, {
		Name: "right",
	}},
	Func: minFloat32,
}, Function{
	Name:        "minDouble",
	Category:    "numbers",
	Description: "Will pick the smallest double of the <left> and <right>.",
	Parameters: Parameters{{
		Name: "left",
	}, {
		Name: "right",
	}},
	Func: minFloat64,
}, Function{
	Name:        "maxFloat",
	Category:    "numbers",
	Description: "Will pick the biggest float of the <left> and <right>.",
	Parameters: Parameters{{
		Name: "left",
	}, {
		Name: "right",
	}},
	Func: maxFloat32,
}, Function{
	Name:        "maxDouble",
	Category:    "numbers",
	Description: "Will pick the biggest double of the <left> and <right>.",
	Parameters: Parameters{{
		Name: "left",
	}, {
		Name: "right",
	}},
	Func: maxFloat64,
}, Function{
	Name:        "floorInteger",
	Category:    "numbers",
	Description: "Returns the greatest integer value less than or equal to <in>.",
	Parameters: Parameters{{
		Name: "in",
	}},
	Func: floorInt,
}, Function{
	Name:        "floorLong",
	Category:    "numbers",
	Description: "Returns the greatest long value less than or equal to <in>.",
	Parameters: Parameters{{
		Name: "in",
	}},
	Func: floorInt64,
}, Function{
	Name:        "ceilInteger",
	Category:    "numbers",
	Description: "Returns the least integer value greater than or equal to <in>.",
	Parameters: Parameters{{
		Name: "in",
	}},
	Func: ceilInt,
}, Function{
	Name:        "ceilLong",
	Category:    "numbers",
	Description: "Returns the least long value greater than or equal to <in>.",
	Parameters: Parameters{{
		Name: "in",
	}},
	Func: ceilInt64,
}, Function{
	Name:        "until",
	Category:    "numbers",
	Description: "Will return an array of numbers of <0> to <count-1>.",
	Parameters: Parameters{{
		Name: "count",
	}},
	Func: until,
}, Function{
	Name:        "untilStep",
	Category:    "numbers",
	Description: "Will return an array of numbers of <start> to <stop> with the given <step> size.",
	Parameters: Parameters{{
		Name: "start",
	}, {
		Name: "stop",
	}, {
		Name: "step",
	}},
	Func: untilStep,
})

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

func maxInt(a interface{}, i ...interface{}) int {
	return int(maxInt64(a, i))
}

func maxInt64(a interface{}, i ...interface{}) int64 {
	aa := toInt64(a)
	for _, b := range i {
		bb := toInt64(b)
		if bb > aa {
			aa = bb
		}
	}
	return aa
}

func maxFloat32(a interface{}, i ...interface{}) float32 {
	return float32(maxFloat64(a, i))
}

func maxFloat64(a interface{}, i ...interface{}) float64 {
	aa := toFloat64(a)
	for _, b := range i {
		bb := toFloat64(b)
		if bb > aa {
			aa = bb
		}
	}
	return aa
}

func minInt(a interface{}, i ...interface{}) int {
	return int(minInt64(a, i))
}

func minInt64(a interface{}, i ...interface{}) int64 {
	aa := toInt64(a)
	for _, b := range i {
		bb := toInt64(b)
		if bb < aa {
			aa = bb
		}
	}
	return aa
}

func minFloat32(a interface{}, i ...interface{}) float32 {
	return float32(minFloat64(a, i))
}

func minFloat64(a interface{}, i ...interface{}) float64 {
	aa := toFloat64(a)
	for _, b := range i {
		bb := toFloat64(b)
		if bb < aa {
			aa = bb
		}
	}
	return aa
}

func until(count int) []int {
	step := 1
	if count < 0 {
		step = -1
	}
	return untilStep(0, count, step)
}

func untilStep(start, stop, step int) []int {
	var v []int

	if stop < start {
		if step >= 0 {
			return v
		}
		for i := start; i > stop; i += step {
			v = append(v, i)
		}
		return v
	}

	if step <= 0 {
		return v
	}
	for i := start; i < stop; i += step {
		v = append(v, i)
	}
	return v
}

func floorInt(a interface{}) int {
	return int(floorInt64(a))
}

func floorInt64(a interface{}) int64 {
	aa := toFloat64(a)
	return int64(math.Floor(aa))
}

func ceilInt(a interface{}) int {
	return int(ceilInt64(a))
}

func ceilInt64(a interface{}) int64 {
	aa := toFloat64(a)
	return int64(math.Ceil(aa))
}

func roundInt(a interface{}) int {
	return int(roundInt64(a))
}

func roundInt64(a interface{}) int64 {
	return int64(math.Round(toFloat64(a)))
}

func roundFloat32(a interface{}, p int) float32 {
	return float32(roundFloat64(a, p))
}

func roundFloat64(a interface{}, p int) float64 {
	roundOn := .5
	val := toFloat64(a)
	places := toFloat64(p)

	var round float64
	pow := math.Pow(10, places)
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	return round / pow
}
