package base

import (
	"bytes"
	"fmt"
	"io"
	"runtime"
	"strings"
)

type Error struct {
	error error
	stack stack
	text  string
}

const (
	gFILTER_KEY = "/errors/gerror/gerror"
)

var (
	goRootForFilter = runtime.GOROOT()
)

func init() {
	if goRootForFilter != "" {
		goRootForFilter = strings.Replace(goRootForFilter, "\\", "/", -1)
	}
}

func (err *Error) Error() string {
	if err.text != "" {
		if err.error != nil {
			return err.text + ": " + err.error.Error()
		}
		return err.text
	}
	return err.error.Error()
}

func (err *Error) Cause() error {
	loop := err
	for loop != nil {
		if loop.error != nil {
			if e, ok := loop.error.(*Error); ok {
				loop = e
			} else {
				return loop.error
			}
		} else {
			return loop
		}
	}
	return nil
}

func (err *Error) Format(s fmt.State, verb rune) {
	switch verb {
	case 's', 'v':
		switch {
		case s.Flag('-'):
			if err.text != "" {
				io.WriteString(s, err.text)
			} else {
				io.WriteString(s, err.Error())
			}
		case s.Flag('+'):
			if verb == 's' {
				io.WriteString(s, err.Stack())
			} else {
				io.WriteString(s, err.Error()+"\n"+err.Stack())
			}
		default:
			io.WriteString(s, err.Error())
		}
	}
}

func (err *Error) Stack() string {
	if err == nil {
		return ""
	}
	loop := err
	index := 1
	buffer := bytes.NewBuffer(nil)
	for loop != nil {
		buffer.WriteString(fmt.Sprintf("%d. %-v\n", index, loop))
		index++
		formatSubStack(loop.stack, buffer)
		if loop.error != nil {
			if e, ok := loop.error.(*Error); ok {
				loop = e
			} else {
				buffer.WriteString(fmt.Sprintf("%d. %s\n", index, loop.error.Error()))
				index++
				break
			}
		} else {
			break
		}
	}
	return buffer.String()
}

func formatSubStack(st stack, buffer *bytes.Buffer) {
	index := 1
	space := "  "
	for _, p := range st {
		if fn := runtime.FuncForPC(p - 1); fn != nil {
			file, line := fn.FileLine(p - 1)
			if strings.Contains(file, gFILTER_KEY) {
				continue
			}

			if strings.Contains(file, "<") {
				continue
			}
			if goRootForFilter != "" && len(file) >= len(goRootForFilter) && file[0:len(goRootForFilter)] == goRootForFilter {
				continue
			}
			if index > 9 {
				space = " "
			}
			buffer.WriteString(fmt.Sprintf("   %d).%s%s\n    \t%s:%d\n", index, space, fn.Name(), file, line))
			index++
		}
	}
}
