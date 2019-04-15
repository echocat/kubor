package model

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

var (
	ErrBeginOfStatementExpected = errors.New("begin of statement expected")
	ErrIllegalStatement         = errors.New("illegal statement")
	ErrUnexpectedToken          = errors.New("unexpected token")

	defaultStatementParser = map[string]StatementParser{}
)

func IsValidStatementName(name string) bool {
	if len(name) <= 2 {
		return false
	}
	if !(name[0] >= 'a' && name[0] <= 'z') {
		return false
	}
	for _, ch := range name {
		if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9')) {
			return false
		}
	}
	return true
}

func RegisterStatementParser(name string, parser StatementParser) StatementParser {
	if !IsValidStatementName(name) {
		panic(fmt.Sprintf("%s is not a valid statement parser name", name))
	}
	defaultStatementParser[name] = parser
	return parser
}

func DefaultStatementParser() (result map[string]StatementParser) {
	result = make(map[string]StatementParser)
	for key, value := range defaultStatementParser {
		result[key] = value
	}
	return result
}

func NewRuntimeEnvironmentParsingError(source string, position [2]uint32, err error, details ...interface{}) *RuntimeEnvironmentParsingError {
	var detail string
	if len(details) >= 2 {
		detail = fmt.Sprintf(details[0].(string), details[1:]...)
	} else if len(details) == 1 {
		detail = details[0].(string)
	}
	return &RuntimeEnvironmentParsingError{
		Position: position,
		Source:   source,
		Err:      err,
		Detail:   detail,
	}
}

func WrapRuntimeEnvironmentParsingError(source string, position [2]uint32, err error) *RuntimeEnvironmentParsingError {
	if rErr, ok := err.(*RuntimeEnvironmentParsingError); ok {
		return rErr
	} else if rErr, ok := err.(RuntimeEnvironmentParsingError); ok {
		return &rErr
	} else {
		return NewRuntimeEnvironmentParsingError(source, position, err)
	}
}

type RuntimeEnvironmentParsingError struct {
	Position [2]uint32
	Source   string
	Err      error
	Detail   string
}

func (instance RuntimeEnvironmentParsingError) Error() string {
	return instance.String()
}

func (instance RuntimeEnvironmentParsingError) String() string {
	message := instance.Err.Error()
	if instance.Detail != "" {
		message += ": " + instance.Detail
	}
	return fmt.Sprintf("%s(%d:%d): %s", instance.Source, instance.Position[0], instance.Position[1], message)
}

type RuntimeEnvironmentParser struct {
	Parser map[string]StatementParser
}

func (instance *RuntimeEnvironmentParser) Parse(source string, target *RuntimeEnvironment, reader io.Reader) error {
	task := runtimeEnvironmentParserTask{
		RuntimeEnvironmentParser: instance,
		target:                   target,
		reader: &SourceReader{
			Reader: bufio.NewReader(reader),
			Source: source,
		},
	}
	return task.execute()
}

type runtimeEnvironmentParserTask struct {
	*RuntimeEnvironmentParser
	target *RuntimeEnvironment
	reader *SourceReader

	lineStarted      bool
	commentStarted   bool
	statementStarted bool
	statementName    strings.Builder
}

func (instance *runtimeEnvironmentParserTask) execute() error {
	source := instance.reader.Source
	position := instance.reader.Position()
	r, err := instance.reader.ReadRune()
	for err == nil {
		if r == '\r' {
			// Ignore
		} else if r == ' ' || r == '\t' {
			if instance.statementStarted {
				if parser, ok := instance.Parser[instance.statementName.String()]; !ok {
					return NewRuntimeEnvironmentParsingError(source, position, ErrIllegalStatement, instance.statementName.String())
				} else if err := parser.Parse(instance.target, instance.reader); err != nil {
					return WrapRuntimeEnvironmentParsingError(source, position, err)
				} else {
					instance.lineStarted = false
					instance.commentStarted = false
					instance.statementStarted = false
					instance.statementName.Reset()
				}
			} else {
				instance.lineStarted = true
			}
		} else if r == '\n' {
			instance.lineStarted = false
			instance.commentStarted = false
			instance.statementStarted = false
			instance.statementName.Reset()
		} else if r == '#' {
			instance.commentStarted = true
			instance.lineStarted = true
		} else if instance.commentStarted {
			// Ignore
		} else if instance.statementStarted {
			instance.statementName.WriteRune(r)
		} else if r >= 'a' && r <= 'z' {
			if instance.lineStarted {
				return NewRuntimeEnvironmentParsingError(source, position, ErrBeginOfStatementExpected, "unexpected token: '%c'", r)
			}
			instance.lineStarted = true
			instance.statementStarted = true
			instance.statementName.WriteRune(r)
		} else {
			return NewRuntimeEnvironmentParsingError(source, position, ErrUnexpectedToken, "'%c'", r)
		}
		position = instance.reader.Position()
		r, err = instance.reader.ReadRune()
	}
	if err == io.EOF {
		err = nil
	}
	if err != nil {
		if instance.statementStarted {
			return NewRuntimeEnvironmentParsingError(source, position, io.EOF)
		}
	}
	return err
}

type StatementParser interface {
	Parse(target *RuntimeEnvironment, reader *SourceReader) error
}

type SourceReader struct {
	Reader *bufio.Reader
	Source string
	Line   uint32
	Column uint32
}

func (instance *SourceReader) ReadRune() (r rune, err error) {
	r, _, err = instance.Reader.ReadRune()
	if r == '\n' {
		instance.Line++
		instance.Column = 0
	}
	if r != '\r' {
		instance.Column++
	}
	return
}

func (instance SourceReader) Position() [2]uint32 {
	return [2]uint32{instance.Line, instance.Column}
}
