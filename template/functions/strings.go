package functions

import (
	"fmt"
	"github.com/aokoli/goutils"
	"github.com/huandu/xstrings"
	"strings"
)

var _ = Register(Function{
	Name:        "abbrev",
	Category:    "strings",
	Description: `Abbreviates a string using ellipses. This will turn  the string "Now is the time for all good men" into "Now is the time for..."`,
	Parameters: Parameters{{
		Name:        "width",
		Description: "Maximum length of result string, must be at least 4",
	}, {
		Name: "in",
	}},
	Func: func(width int, in string) (string, error) {
		if width < 4 {
			return in, nil
		}
		return goutils.Abbreviate(in, width)
	},
}, Function{
	Name:     "abbrevFull",
	Category: "strings",
	Description: `AbbreviateFull abbreviates a string using ellipses. This will turn the string "Now is the time for all good men" into "...is the time for..."
This function works like Abbreviate(string, int), but allows you to specify a "left edge" offset. Note that this left edge is not
necessarily going to be the leftmost character in the result, or the first character following the ellipses, but it will appear
somewhere in the result.
In no case will it return a string of length greater than maxWidth.`,
	Parameters: Parameters{{
		Name:        "offset",
		Description: "Left edge of source string",
	}, {
		Name:        "width",
		Description: "Maximum length of result string, must be at least 4",
	}, {
		Name: "in",
	}},
	Func: func(offset, width int, in string) (string, error) {
		if width < 4 || offset > 0 && width < 7 {
			return in, nil
		}
		return goutils.AbbreviateFull(in, offset, width)
	},
}, Function{
	Name:     "upper",
	Category: "strings",
	Parameters: Parameters{{
		Name: "in",
	}},
	Returns: Return{
		Description: `A copy of the string s with all Unicode letters mapped to their upper case.`,
	},
	Func: strings.ToUpper,
}, Function{
	Name:     "lower",
	Category: "strings",
	Parameters: Parameters{{
		Name: "in",
	}},
	Returns: Return{
		Description: `A copy of the string s with all Unicode letters mapped to their lower case.`,
	},
	Func: strings.ToLower,
}, Function{
	Name:        "trim",
	Category:    "strings",
	Description: "Removes space from either side of a string",
	Parameters: Parameters{{
		Name: "in",
	}},
	Func: strings.TrimSpace,
}, Function{
	Name:        "trimAll",
	Category:    "strings",
	Description: "Remove given characters from the front or back of a string.",
	Parameters: Parameters{{
		Name: "toRemove",
	}, {
		Name: "in",
	}},
	Func: func(toRemove string, in string) string {
		return strings.Trim(in, toRemove)
	},
}, Function{
	Name:        "trimSuffix",
	Category:    "strings",
	Description: "Remove given characters from the back of a string.",
	Parameters: Parameters{{
		Name: "toRemove",
	}, {
		Name: "in",
	}},
	Func: func(toRemove string, in string) string {
		return strings.TrimSuffix(in, toRemove)
	},
}, Function{
	Name:     "capitalize",
	Category: "strings",
	Description: `Capitalize capitalizes all the delimiter separated words in a string. Only the first letter of each word is changed.
To convert the rest of each word to lowercase at the same time, use CapitalizeFully(str string, delimiters ...rune).
The delimiters represent a set of characters understood to separate words. The first string character
and the first non-delimiter character after a delimiter will be capitalized. A "" input string returns "".
Capitalization uses the Unicode title case, normally equivalent to upper case.`,
	Parameters: Parameters{{
		Name: "in",
	}},
	Func: func(in string) string {
		return goutils.Capitalize(in)
	},
}, Function{
	Name:     "uncapitalize",
	Category: "strings",
	Description: `Uncapitalize uncapitalizes all the whitespace separated words in a string. Only the first letter of each word is changed.
The delimiters represent a set of characters understood to separate words. The first string character and the first non-delimiter
character after a delimiter will be uncapitalized. Whitespace is defined by unicode.IsSpace(char).`,
	Parameters: Parameters{{
		Name: "in",
	}},
	Func: func(in string) string {
		return goutils.Uncapitalize(in)
	},
}, Function{
	Name:        "replace",
	Category:    "strings",
	Description: "Replaces the given <old> string with the <new> string.",
	Parameters: Parameters{{
		Name: "old",
	}, {
		Name: "new",
	}, {
		Name: "in",
	}},
	Func: func(old string, n string, in string) string {
		return strings.Replace(in, old, n, -1)
	},
}, Function{
	Name:        "repeat",
	Category:    "strings",
	Description: "Repeat a string multiple times.",
	Parameters: Parameters{{
		Name: "count",
	}, {
		Name: "in",
	}},
	Func: func(count int, in string) string {
		return strings.Repeat(in, count)
	},
}, Function{
	Name:        "substr",
	Category:    "strings",
	Description: "Get a substring from a string.",
	Parameters: Parameters{{
		Name: "start",
	}, {
		Name: "length",
	}, {
		Name: "in",
	}},
	Func: func(start int, length int, in string) string {
		if start < 0 {
			return in[:length]
		}
		if length < 0 {
			return in[start:]
		}
		return in[start:length]
	},
}, Function{
	Name:        "trunc",
	Category:    "strings",
	Description: "Truncate a string (and add no suffix).",
	Parameters: Parameters{{
		Name: "length",
	}, {
		Name: "in",
	}},
	Func: func(length int, in string) string {
		if len(in) <= length {
			return in
		}
		return in[0:length]
	},
}, Function{
	Name:        "initials",
	Category:    "strings",
	Description: "Given multiple words, take the first letter of each word and combine.",
	Parameters: Parameters{{
		Name: "in",
	}},
	Func: func(in string) string {
		return goutils.Initials(in)
	},
}, Function{
	Name:        "randAlphaNum",
	Category:    "strings",
	Description: "These four functions generate random strings, but with different base character sets of [0-9A-Za-z].",
	Parameters: Parameters{{
		Name: "count",
	}},
	Func: func(count int) (string, error) {
		return goutils.RandomAlphaNumeric(count)
	},
}, Function{
	Name:        "randAlpha",
	Category:    "strings",
	Description: "These four functions generate random strings, but with different base character sets of [A-Za-z].",
	Parameters: Parameters{{
		Name: "count",
	}},
	Func: func(count int) (string, error) {
		return goutils.RandomAlphabetic(count)
	},
}, Function{
	Name:        "randNum",
	Category:    "strings",
	Description: "These four functions generate random strings, but with different base character sets of [0-9].",
	Parameters: Parameters{{
		Name: "count",
	}},
	Func: func(count int) (string, error) {
		return goutils.RandomNumeric(count)
	},
}, Function{
	Name:     "warp",
	Category: "strings",
	Description: `Wrap wraps a single line of text, identifying words by ' '.
New lines will be separated by '\n'. Very long words, such as URLs will not be wrapped.
Leading spaces on a new line are stripped. Trailing spaces are not stripped.`,
	Parameters: Parameters{{
		Name: "length",
	}, {
		Name: "in",
	}},
	Func: func(length int, in string) string {
		return goutils.Wrap(in, length)
	},
}, Function{
	Name:     "warpCustom",
	Category: "strings",
	Description: `WrapCustom wraps a single line of text, identifying words by ' '.
Leading spaces on a new line are stripped. Trailing spaces are not stripped.`,
	Parameters: Parameters{{
		Name: "length",
	}, {
		Name: "newLine",
	}, {
		Name: "wrapLongWords",
	}, {
		Name: "in",
	}},
	Func: func(length int, newLine string, wrapLongWords bool, in string) string {
		return goutils.WrapCustom(in, length, newLine, wrapLongWords)
	},
}, Function{
	Name:        "hasPrefix",
	Category:    "strings",
	Description: `Test whether a string has a given prefix.`,
	Parameters: Parameters{{
		Name: "toSearchFor",
	}, {
		Name: "in",
	}},
	Func: func(toSearchFor string, in string) bool {
		return strings.HasPrefix(in, toSearchFor)
	},
}, Function{
	Name:        "hasSuffix",
	Category:    "strings",
	Description: `Test whether a string has a given suffix.`,
	Parameters: Parameters{{
		Name: "toSearchFor",
	}, {
		Name: "in",
	}},
	Func: func(toSearchFor string, in string) bool {
		return strings.HasSuffix(in, toSearchFor)
	},
}, Function{
	Name:        "quote",
	Category:    "strings",
	Description: `Wrap a string in double quotes.`,
	Parameters: Parameters{{
		Name: "in",
	}},
	Func: func(str ...interface{}) string {
		out := make([]string, len(str))
		for i, s := range str {
			out[i] = fmt.Sprintf("%q", strval(s))
		}
		return strings.Join(out, " ")
	},
}, Function{
	Name:        "sQuote",
	Category:    "strings",
	Description: `Wrap a string in single quotes.`,
	Parameters: Parameters{{
		Name: "in",
	}},
	Func: func(str ...interface{}) string {
		out := make([]string, len(str))
		for i, s := range str {
			out[i] = fmt.Sprintf("'%v'", s)
		}
		return strings.Join(out, " ")
	},
}, Function{
	Name:        "cat",
	Category:    "strings",
	Description: `Concatenates multiple strings together into one, separating them with spaces.`,
	Parameters: Parameters{{
		Name: "in",
	}},
	Func: func(v ...interface{}) string {
		r := strings.TrimSpace(strings.Repeat("%v ", len(v)))
		return fmt.Sprintf(r, v...)
	},
}, Function{
	Name:        "indent",
	Category:    "strings",
	Description: `Indents every line in a given string to the specified indent width. This is useful when aligning multi-line strings.`,
	Parameters: Parameters{{
		Name: "indent",
	}, {
		Name: "str",
	}},
	Func: func(indent int, str string) string {
		pad := strings.Repeat(" ", indent)
		return pad + strings.Replace(str, "\n", "\n"+pad, -1)
	},
}, Function{
	Name:        "snakeCase",
	Category:    "strings",
	Description: `convert all upper case characters in a string to snake case format.`,
	Parameters: Parameters{{
		Name: "in",
	}},
	Func: xstrings.ToSnakeCase,
}, Function{
	Name:        "camelCase",
	Category:    "strings",
	Description: `Convert all lower case characters behind underscores to upper case character.`,
	Parameters: Parameters{{
		Name: "in",
	}},
	Func: xstrings.ToCamelCase,
}, Function{
	Name:        "kebabCase",
	Category:    "strings",
	Description: `Convert all upper case characters in a string to kebab case format.`,
	Parameters: Parameters{{
		Name: "in",
	}},
	Func: xstrings.ToKebabCase,
}, Function{
	Name:        "shuffle",
	Category:    "strings",
	Description: `Shuffle randomizes runes in a string and returns the result.`,
	Parameters: Parameters{{
		Name: "in",
	}},
	Func: xstrings.Shuffle,
})

func strval(v interface{}) string {
	switch v := v.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case error:
		return v.Error()
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}
