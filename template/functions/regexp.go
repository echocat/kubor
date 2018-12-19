package functions

import (
	"regexp"
)

var _ = Register(Function{
	Name:     "regexMatch",
	Category: "regexp",
	Description: `Reports whether the string <str> contains any match of the regular expression <pattern>.
See https://golang.org/pkg/regexp/ for more details.`,
	Parameters: Parameters{{
		Name: "pattern",
	}, {
		Name: "str",
	}},
	Func: func(pattern string, str string) (bool, error) {
		return regexp.MatchString(pattern, str)
	},
}, Function{
	Name:     "regexFindAll",
	Category: "regexp",
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
	Func: func(pattern string, n int, str string) ([]string, error) {
		if r, err := regexp.Compile(pattern); err != nil {
			return []string{}, err
		} else {
			return r.FindAllString(str, n), nil
		}
	},
}, Function{
	Name:     "regexFind",
	Category: "regexp",
	Description: `Returns the first match of the regular expression <pattern> in the input string.
See https://golang.org/pkg/regexp/ for more details.`,
	Parameters: Parameters{{
		Name: "pattern",
	}, {
		Name: "str",
	}},
	Func: func(pattern string, str string) (string, error) {
		if r, err := regexp.Compile(pattern); err != nil {
			return "", err
		} else {
			return r.FindString(str), nil
		}
	},
}, Function{
	Name:     "regexReplaceAll",
	Category: "regexp",
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
	Func: func(pattern string, replacement string, str string) (string, error) {
		if r, err := regexp.Compile(pattern); err != nil {
			return "", err
		} else {
			return r.ReplaceAllString(str, replacement), nil
		}
	},
}, Function{
	Name:     "regexSplit",
	Category: "regexp",
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
	Func: func(pattern string, n int, str string) ([]string, error) {
		if r, err := regexp.Compile(pattern); err != nil {
			return []string{}, err
		} else {
			return r.Split(str, n), nil
		}
	},
})
