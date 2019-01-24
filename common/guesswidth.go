// +build appengine !linux,!freebsd,!darwin,!dragonfly,!netbsd,!openbsd

package common

import "io"

func GuessOutputWidth(w io.Writer) int {
	return 80
}
