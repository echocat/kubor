package common

import (
	"io"
)

func ToReadCloser(reader io.Reader) io.ReadCloser {
	return &readCloserDelegate{
		reader,
	}
}

type readCloserDelegate struct {
	io.Reader
}

func (instance *readCloserDelegate) Close() error {
	return nil
}
