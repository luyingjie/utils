package verror

import (
	"fmt"
)

type ApiStack interface {
	Error() string
	Stack() string
}

type ApiCause interface {
	Error() string
	Cause() error
}

func New(text string) error {
	if text == "" {
		return nil
	}
	return &Error{
		stack: callers(),
		text:  text,
	}
}

func NewSkip(skip int, text string) error {
	if text == "" {
		return nil
	}
	return &Error{
		stack: callers(skip),
		text:  text,
	}
}

func Newf(format string, args ...interface{}) error {
	if format == "" {
		return nil
	}
	return &Error{
		stack: callers(),
		text:  fmt.Sprintf(format, args...),
	}
}

func NewfSkip(skip int, format string, args ...interface{}) error {
	if format == "" {
		return nil
	}
	return &Error{
		stack: callers(skip),
		text:  fmt.Sprintf(format, args...),
	}
}

func Wrap(err error, text string) error {
	if err == nil {
		return nil
	}
	return &Error{
		error: err,
		stack: callers(),
		text:  text,
	}
}

func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return &Error{
		error: err,
		stack: callers(),
		text:  fmt.Sprintf(format, args...),
	}
}

func Cause(err error) error {
	if err != nil {
		if e, ok := err.(ApiCause); ok {
			return e.Cause()
		}
	}
	return err
}

func Stack(err error) string {
	if err == nil {
		return ""
	}
	if e, ok := err.(ApiStack); ok {
		return e.Stack()
	}
	return ""
}
