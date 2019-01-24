package functions

import (
	"math"
)

var FuncMinInt = Function{
	Description: "Will pick the smallest int of the <left> and <right>.",
	Parameters: Parameters{{
		Name: "left",
	}, {
		Name: "right",
	}},
}.MustWithFunc(minInt)

var FuncMinInt64 = Function{
	Description: "Will pick the smallest int64 of the <left> and <right>.",
	Parameters: Parameters{{
		Name: "left",
	}, {
		Name: "right",
	}},
}.MustWithFunc(minInt64)

var FuncMaxInt = Function{
	Description: "Will pick the biggest int of the <left> and <right>.",
	Parameters: Parameters{{
		Name: "left",
	}, {
		Name: "right",
	}},
}.MustWithFunc(maxInt)

var FuncMaxInt64 = Function{
	Description: "Will pick the biggest int64 of the <left> and <right>.",
	Parameters: Parameters{{
		Name: "left",
	}, {
		Name: "right",
	}},
}.MustWithFunc(maxInt64)

var FuncMinFloat = Function{
	Description: "Will pick the smallest float of the <left> and <right>.",
	Parameters: Parameters{{
		Name: "left",
	}, {
		Name: "right",
	}},
}.MustWithFunc(minFloat32)

var FuncMinDouble = Function{
	Description: "Will pick the smallest float64 of the <left> and <right>.",
	Parameters: Parameters{{
		Name: "left",
	}, {
		Name: "right",
	}},
}.MustWithFunc(minFloat64)

var FuncMaxFloat = Function{
	Description: "Will pick the biggest float of the <left> and <right>.",
	Parameters: Parameters{{
		Name: "left",
	}, {
		Name: "right",
	}},
}.MustWithFunc(maxFloat32)

var FuncMaxDouble = Function{
	Description: "Will pick the biggest float64 of the <left> and <right>.",
	Parameters: Parameters{{
		Name: "left",
	}, {
		Name: "right",
	}},
}.MustWithFunc(maxFloat64)

var FuncFloorInt = Function{
	Description: "Returns the greatest int value less than or equal to <in>.",
	Parameters: Parameters{{
		Name: "in",
	}},
}.MustWithFunc(floorInt)

var FuncFloorInt64 = Function{
	Description: "Returns the greatest int64 value less than or equal to <in>.",
	Parameters: Parameters{{
		Name: "in",
	}},
}.MustWithFunc(floorInt64)

var FuncCeilInt = Function{
	Description: "Returns the least int value greater than or equal to <in>.",
	Parameters: Parameters{{
		Name: "in",
	}},
}.MustWithFunc(ceilInt)

var FuncCeilInt64 = Function{
	Description: "Returns the least int64 value greater than or equal to <in>.",
	Parameters: Parameters{{
		Name: "in",
	}},
}.MustWithFunc(ceilInt64)

var FuncUntil = Function{
	Description: "Will return an array of numbers of <0> to <count-1>.",
	Parameters: Parameters{{
		Name: "count",
	}},
}.MustWithFunc(until)

var FuncUntilStep = Function{
	Description: "Will return an array of numbers of <start> to <stop> with the given <step> size.",
	Parameters: Parameters{{
		Name: "start",
	}, {
		Name: "stop",
	}, {
		Name: "step",
	}},
}.MustWithFunc(untilStep)

var FuncRoundInt = Function{
	Description: "Will round the given input to an int.",
	Parameters: Parameters{{
		Name: "in",
	}},
}.MustWithFunc(roundInt)

var FuncRoundInt64 = Function{
	Description: "Will round the given input to a int64.",
	Parameters: Parameters{{
		Name: "in",
	}},
}.MustWithFunc(roundInt64)

var FuncRoundFloat = Function{
	Description: "Will round the given input to a float.",
	Parameters: Parameters{{
		Name:        "precision",
		Description: "Defines how many decimal positions should remain.",
	}, {
		Name: "in",
	}},
}.MustWithFunc(func(precision int, in interface{}) float32 {
	return roundFloat32(in, precision)
})

var FuncRoundDouble = Function{
	Description: "Will round the given input to a float64.",
	Parameters: Parameters{{
		Name:        "precision",
		Description: "Defines how many decimal positions should remain.",
	}, {
		Name: "in",
	}},
}.MustWithFunc(func(precision int, in interface{}) float64 {
	return roundFloat64(in, precision)
})

var FuncsMath = Functions{
	"roundFloat64": FuncRoundDouble,
	"roundFloat32": FuncRoundFloat,
	"roundInt64":   FuncRoundInt64,
	"roundInt":     FuncRoundInt,
	"untilStep":    FuncUntilStep,
	"ceilInt64":    FuncCeilInt64,
	"ceilInt":      FuncCeilInt,
	"floorInt64":   FuncFloorInt64,
	"minInt":       FuncMinInt,
	"minInt64":     FuncMinInt64,
	"maxInt":       FuncMaxInt,
	"maxInt64":     FuncMaxInt64,
	"minFloat32":   FuncMinFloat,
	"minFloat64":   FuncMinDouble,
	"maxFloat32":   FuncMaxFloat,
	"maxFloat64":   FuncMaxDouble,
	"floorInt":     FuncFloorInt,
	"until":        FuncUntil,
}
var CategoryMath = Category{
	Functions: FuncsMath,
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
