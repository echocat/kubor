package functions

import (
	"regexp"
)

var FuncRegexpMatch = Function{
	Description: `Reports whether the string <str> contains any match of the regular expression <pattern>.
See https://golang.org/pkg/regexp/ for more details.`,
	Parameters: Parameters{{
		Name: "pattern",
	}, {
		Name: "str",
	}},
}.MustWithFunc(func(pattern string, str string) (bool, error) {
	return regexp.MatchString(pattern, str)
})

var FuncRegexpFindAll = Function{
	Description: `Returns an array of all matches of the regular expression <pattern> in the input string.
See https://golang.org/pkg/regexp/ for more details.`,
	Parameters: Parameters{{
		Name: "pattern",
	}, {
		Name:        "n",
		Description: "Limits the maximum number of matches that should be returned. If smaller then <0> it will be unlimited.",
	}, {
		Name: "str",
	}},
}.MustWithFunc(func(pattern string, n int, str string) ([]string, error) {
	if r, err := regexp.Compile(pattern); err != nil {
		return []string{}, err
	} else {
		return r.FindAllString(str, n), nil
	}
})

var FuncRegexpFind = Function{
	Description: `Returns the first match of the regular expression <pattern> in the input string.
See https://golang.org/pkg/regexp/ for more details.`,
	Parameters: Parameters{{
		Name: "pattern",
	}, {
		Name: "str",
	}},
}.MustWithFunc(func(pattern string, str string) (string, error) {
	if r, err := regexp.Compile(pattern); err != nil {
		return "", err
	} else {
		return r.FindString(str), nil
	}
})

var FuncRegexpReplaceAll = Function{
	Description: `Returns copy of src, replacing matches of the Regexp with the replacement string repl.
Inside repl, $ signs are interpreted as in Expand, so for instance $1 represents the text of the first submatch.
See https://golang.org/pkg/regexp/ for more details.`,
	Parameters: Parameters{{
		Name: "pattern",
	}, {
		Name: "replacement",
	}, {
		Name: "str",
	}},
}.MustWithFunc(func(pattern string, replacement string, str string) (string, error) {
	if r, err := regexp.Compile(pattern); err != nil {
		return "", err
	} else {
		return r.ReplaceAllString(str, replacement), nil
	}
})

var FuncRegexpSplit = Function{
	Description: `Splits into substrings separated by the expression and returns a slice of the substrings between those expression matches.
See https://golang.org/pkg/regexp/ for more details.`,
	Parameters: Parameters{{
		Name: "pattern",
	}, {
		Name:        "n",
		Description: "Limits the maximum number of parts that should be returned. If smaller then <0> it will be unlimited.",
	}, {
		Name: "str",
	}},
}.MustWithFunc(func(pattern string, n int, str string) ([]string, error) {
	if r, err := regexp.Compile(pattern); err != nil {
		return []string{}, err
	} else {
		return r.Split(str, n), nil
	}
})

var FuncsRegexp = Functions{
	"regexpMatch":      FuncRegexpMatch,
	"regexpFindAll":    FuncRegexpFindAll,
	"regexpFind":       FuncRegexpFind,
	"regexpReplaceAll": FuncRegexpReplaceAll,
	"regexpSplit":      FuncRegexpSplit,
}
var CategoryRegexp = Category{
	Functions: FuncsRegexp,
}
