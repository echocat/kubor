package model

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	StreamTargetStdout = StreamTarget("stdout")
	StreamTargetStderr = StreamTarget("stderr")
)

type StreamTarget string

func (instance *StreamTarget) Set(plain string) error {
	return instance.UnmarshalText([]byte(plain))
}

func (instance StreamTarget) String() string {
	return string(instance)
}

func (instance StreamTarget) MarshalText() (text []byte, err error) {
	return []byte(instance), nil
}

func (instance *StreamTarget) UnmarshalText(text []byte) error {
	*instance = StreamTarget(text)
	return nil
}

func (instance StreamTarget) OpenForWrite() (io.WriteCloser, error) {
	switch instance {
	case StreamTargetStdout:
		return &streamTargetWriterWrapper{os.Stdout, false}, nil
	case StreamTargetStderr:
		return &streamTargetWriterWrapper{os.Stderr, false}, nil
	default:
		dir := filepath.Dir(string(instance))
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("cannot ensure parent directory of %v: %w", instance, err)
		} else if f, err := os.OpenFile(string(instance), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644); err != nil {
			return nil, fmt.Errorf("cannot open %v: %w", instance, err)
		} else {
			return &streamTargetWriterWrapper{f, true}, nil
		}
	}
}

type streamTargetWriterWrapper struct {
	io.WriteCloser
	canBeClosed bool
}

func (instance *streamTargetWriterWrapper) Close() error {
	if !instance.canBeClosed {
		return nil
	}
	return instance.WriteCloser.Close()
}
