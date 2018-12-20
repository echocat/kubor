package functions

import (
	"fmt"
	"github.com/aokoli/goutils"
	"github.com/huandu/xstrings"
	"strings"
)

var FuncAbbrev = Function{
	Description: `Abbreviates a string using ellipses. This will turn  the string "Now is the time for all good men" into "Now is the time for..."`,
	Parameters: Parameters{{
		Name:        "width",
		Description: "Maximum length of result string, must be at least 4",
	}, {
		Name: "in",
	}},
}.MustWithFunc(func(width int, in string) (string, error) {
	if width < 4 {
		return in, nil
	}
	return goutils.Abbreviate(in, width)
})

var FuncAbbrevFull = Function{
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
}.MustWithFunc(func(offset, width int, in string) (string, error) {
	if width < 4 || offset > 0 && width < 7 {
		return in, nil
	}
	return goutils.AbbreviateFull(in, offset, width)
})

var FuncUpper = Function{
	Parameters: Parameters{{
		Name: "in",
	}},
	Returns: Return{
		Description: `A copy of the string s with all Unicode letters mapped to their upper case.`,
	},
}.MustWithFunc(strings.ToUpper)

var FuncLower = Function{
	Parameters: Parameters{{
		Name: "in",
	}},
	Returns: Return{
		Description: `A copy of the string s with all Unicode letters mapped to their lower case.`,
	},
}.MustWithFunc(strings.ToLower)

var FuncTrim = Function{
	Description: "Removes space from either side of a string",
	Parameters: Parameters{{
		Name: "in",
	}},
}.MustWithFunc(strings.TrimSpace)

var FuncTrimAll = Function{
	Description: "Remove given characters from the front or back of a string.",
	Parameters: Parameters{{
		Name: "toRemove",
	}, {
		Name: "in",
	}},
}.MustWithFunc(func(toRemove string, in string) string {
	return strings.Trim(in, toRemove)
})

var FuncTrimSuffix = Function{
	Description: "Remove given characters from the back of a string.",
	Parameters: Parameters{{
		Name: "toRemove",
	}, {
		Name: "in",
	}},
}.MustWithFunc(func(toRemove string, in string) string {
	return strings.TrimSuffix(in, toRemove)
})

var FuncCapitalize = Function{
	Description: `Capitalize capitalizes all the delimiter separated words in a string. Only the first letter of each word is changed.
To convert the rest of each word to lowercase at the same time, use CapitalizeFully(str string, delimiters ...rune).
The delimiters represent a set of characters understood to separate words. The first string character
and the first non-delimiter character after a delimiter will be capitalized. A "" input string returns "".
Capitalization uses the Unicode title case, normally equivalent to upper case.`,
	Parameters: Parameters{{
		Name: "in",
	}},
}.MustWithFunc(func(in string) string {
	return goutils.Capitalize(in)
})

var FuncUncapitalize = Function{
	Description: `Uncapitalize uncapitalizes all the whitespace separated words in a string. Only the first letter of each word is changed.
The delimiters represent a set of characters understood to separate words. The first string character and the first non-delimiter
character after a delimiter will be uncapitalized. Whitespace is defined by unicode.IsSpace(char).`,
	Parameters: Parameters{{
		Name: "in",
	}},
}.MustWithFunc(func(in string) string {
	return goutils.Uncapitalize(in)
})

var FuncReplace = Function{
	Description: "Replaces the given <old> string with the <new> string.",
	Parameters: Parameters{{
		Name: "old",
	}, {
		Name: "new",
	}, {
		Name: "in",
	}},
}.MustWithFunc(func(old string, n string, in string) string {
	return strings.Replace(in, old, n, -1)
})

var FuncRepeat = Function{
	Description: "Repeat a string multiple times.",
	Parameters: Parameters{{
		Name: "count",
	}, {
		Name: "in",
	}},
}.MustWithFunc(func(count int, in string) string {
	return strings.Repeat(in, count)
})

var FuncSubstr = Function{
	Description: "Get a substring from a string.",
	Parameters: Parameters{{
		Name: "start",
	}, {
		Name: "length",
	}, {
		Name: "in",
	}},
}.MustWithFunc(func(start int, length int, in string) string {
	if start < 0 {
		return in[:length]
	}
	if length < 0 {
		return in[start:]
	}
	return in[start:length]
})

var FuncTrunc = Function{
	Description: "Truncate a string (and add no suffix).",
	Parameters: Parameters{{
		Name: "length",
	}, {
		Name: "in",
	}},
}.MustWithFunc(func(length int, in string) string {
	if len(in) <= length {
		return in
	}
	return in[0:length]
})

var FuncInitials = Function{
	Description: "Given multiple words, take the first letter of each word and combine.",
	Parameters: Parameters{{
		Name: "in",
	}},
}.MustWithFunc(func(in string) string {
	return goutils.Initials(in)
})

var FuncRandAlphaNum = Function{
	Description: "These four functions generate random strings, but with different base character sets of [0-9A-Za-z].",
	Parameters: Parameters{{
		Name: "count",
	}},
}.MustWithFunc(func(count int) (string, error) {
	return goutils.RandomAlphaNumeric(count)
})

var FuncRandAlpha = Function{
	Description: "These four functions generate random strings, but with different base character sets of [A-Za-z].",
	Parameters: Parameters{{
		Name: "count",
	}},
}.MustWithFunc(func(count int) (string, error) {
	return goutils.RandomAlphabetic(count)
})

var FuncRandNum = Function{
	Description: "These four functions generate random strings, but with different base character sets of [0-9].",
	Parameters: Parameters{{
		Name: "count",
	}},
}.MustWithFunc(func(count int) (string, error) {
	return goutils.RandomNumeric(count)
})

var FuncWarp = Function{
	Description: `Wrap wraps a single line of text, identifying words by ' '.
New lines will be separated by '\n'. Very int64 words, such as URLs will not be wrapped.
Leading spaces on a new line are stripped. Trailing spaces are not stripped.`,
	Parameters: Parameters{{
		Name: "length",
	}, {
		Name: "in",
	}},
}.MustWithFunc(func(length int, in string) string {
	return goutils.Wrap(in, length)
})

var FuncWarpCustom = Function{
	Description: `WrapCustom wraps a single line of text, identifying words by ' '.
Leading spaces on a new line are stripped. Trailing spaces are not stripped.`,
	Parameters: Parameters{{
		Name: "length",
	}, {
		Name: "newLine",
	}, {
		Name: "wrapInt64Words",
	}, {
		Name: "in",
	}},
}.MustWithFunc(func(length int, newLine string, wrapInt64Words bool, in string) string {
	return goutils.WrapCustom(in, length, newLine, wrapInt64Words)
})

var FuncHasPrefix = Function{
	Description: `Test whether a string has a given prefix.`,
	Parameters: Parameters{{
		Name: "toSearchFor",
	}, {
		Name: "in",
	}},
}.MustWithFunc(func(toSearchFor string, in string) bool {
	return strings.HasPrefix(in, toSearchFor)
})

var FuncHasSuffix = Function{
	Description: `Test whether a string has a given suffix.`,
	Parameters: Parameters{{
		Name: "toSearchFor",
	}, {
		Name: "in",
	}},
}.MustWithFunc(func(toSearchFor string, in string) bool {
	return strings.HasSuffix(in, toSearchFor)
})

var FuncQuote = Function{
	Description: `Wrap a string in float64 quotes.`,
	Parameters: Parameters{{
		Name: "in",
	}},
}.MustWithFunc(func(str ...interface{}) string {
	out := make([]string, len(str))
	for i, s := range str {
		out[i] = fmt.Sprintf("%q", strval(s))
	}
	return strings.Join(out, " ")
})

var FuncSQuote = Function{
	Description: `Wrap a string in single quotes.`,
	Parameters: Parameters{{
		Name: "in",
	}},
}.MustWithFunc(func(str ...interface{}) string {
	out := make([]string, len(str))
	for i, s := range str {
		out[i] = fmt.Sprintf("'%v'", s)
	}
	return strings.Join(out, " ")
})

var FuncCat = Function{
	Description: `Concatenates multiple strings together into one, separating them with spaces.`,
	Parameters: Parameters{{
		Name: "in",
	}},
}.MustWithFunc(func(v ...interface{}) string {
	r := strings.TrimSpace(strings.Repeat("%v ", len(v)))
	return fmt.Sprintf(r, v...)
})

var FuncIndent = Function{
	Description: `Indents every line in a given string to the specified indent width. This is useful when aligning multi-line strings.`,
	Parameters: Parameters{{
		Name: "indent",
	}, {
		Name: "str",
	}},
}.MustWithFunc(func(indent int, str string) string {
	pad := strings.Repeat(" ", indent)
	return pad + strings.Replace(str, "\n", "\n"+pad, -1)
})

var FuncSnakeCase = Function{
	Description: `convert all upper case characters in a string to snake case format.`,
	Parameters: Parameters{{
		Name: "in",
	}},
}.MustWithFunc(xstrings.ToSnakeCase)

var FuncCamelCase = Function{
	Description: `Convert all lower case characters behind underscores to upper case character.`,
	Parameters: Parameters{{
		Name: "in",
	}},
}.MustWithFunc(xstrings.ToCamelCase)

var FuncKebabCase = Function{
	Description: `Convert all upper case characters in a string to kebab case format.`,
	Parameters: Parameters{{
		Name: "in",
	}},
}.MustWithFunc(xstrings.ToKebabCase)

var FuncShuffle = Function{
	Description: `Shuffle randomizes runes in a string and returns the result.`,
	Parameters: Parameters{{
		Name: "in",
	}},
}.MustWithFunc(xstrings.Shuffle)

var FuncsStrings = Functions{
	"abbrev":       FuncAbbrev,
	"abbrevFull":   FuncAbbrevFull,
	"upper":        FuncUpper,
	"lower":        FuncLower,
	"trim":         FuncTrim,
	"trimAll":      FuncTrimAll,
	"trimSuffix":   FuncTrimSuffix,
	"capitalize":   FuncCapitalize,
	"uncapitalize": FuncUncapitalize,
	"replace":      FuncReplace,
	"repeat":       FuncRepeat,
	"substr":       FuncSubstr,
	"trunc":        FuncTrunc,
	"initials":     FuncInitials,
	"randAlphaNum": FuncRandAlphaNum,
	"randAlpha":    FuncRandAlpha,
	"randNum":      FuncRandNum,
	"warp":         FuncWarp,
	"warpCustom":   FuncWarpCustom,
	"hasPrefix":    FuncHasPrefix,
	"hasSuffix":    FuncHasSuffix,
	"quote":        FuncQuote,
	"sQuote":       FuncSQuote,
	"cat":          FuncCat,
	"indent":       FuncIndent,
	"snakeCase":    FuncSnakeCase,
	"camelCase":    FuncCamelCase,
	"kebabCase":    FuncKebabCase,
	"shuffle":      FuncShuffle,
}
var CategoryStrings = Category{
	Functions: FuncsStrings,
}

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
