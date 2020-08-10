package support

import (
	"strings"
)

func NormalizeLabelValue(in string) string {
	if in == "" {
		return ""
	}
	if len(in) > 63 {
		in = in[:63]
	}
	buf := []byte(strings.TrimFunc(in, func(c rune) bool {
		return !((c >= 'A' && c <= 'Z') ||
			(c >= 'a' && c <= 'z') ||
			(c >= '0' && c <= '9'))
	}))
	for i, c := range buf {
		if !((c >= 'A' && c <= 'Z') ||
			(c >= 'a' && c <= 'z') ||
			(c >= '0' && c <= '9') ||
			c == '_' || c == '-' || c == '.') {
			buf[i] = '_'
		}
	}
	return string(buf)
}
