package common

import "time"

func Pstring(in string) *string {
	return &in
}

func Pbool(in bool) *bool {
	return &in
}

func PtimeDuration(in time.Duration) *time.Duration {
	return &in
}
