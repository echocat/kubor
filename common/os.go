package common

import (
	"os"
	"strings"
)

func Environ() map[string]string {
	result := map[string]string{}
	for _, keyAndValue := range os.Environ() {
		part := strings.SplitN(keyAndValue, "=", 2)
		if len(part) > 1 {
			result[part[0]] = part[1]
		}
	}
	return result
}
