// Package erf provides error management with stack trace.
package erf

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"unsafe"
)

// Erf is an error type that wraps the underlying error that stores and formats the stack trace.
type Erf struct {
	err        error
	format     string
	args       []interface{}
	tagIndexes map[string]int
	pc         []uintptr
}

// Error is implementation of error.
func (e *Erf) Error() string {
	return e.err.Error()
}

// Unwrap returns the underlying error.
func (e *Erf) Unwrap() error {
	if err, ok := e.err.(WrappedError); ok {
		return err.Unwrap()
	}
	return nil
}

// Format is implementation of fmt.Formatter.
func (e *Erf) Format(f fmt.State, verb rune) {
	buf := bytes.NewBuffer(nil)
	switch verb {
	case 's', 'v':
		if !f.Flag('+') {
			buf.WriteString(e.err.Error())
			break
		}
		format := "%+"
		for _, r := range []rune{'-', '#'} {
			if f.Flag(int(r)) {
				format += string(r)
			}
		}
		pad, wid, prec := byte('\t'), 0, 1
		if f.Flag('-') {
			pad = ' '
			prec = 2
		}
		if w, ok := f.Width(); ok {
			wid = w
		}
		if p, ok := f.Precision(); ok {
			prec = p
		}
		format += fmt.Sprintf("%d.%d", wid, prec)
		format += "s"
		for _, line := range strings.Split(e.err.Error(), "\n") {
			buf.WriteString(fmt.Sprintf("%s%s", strings.Repeat(string(pad), wid+prec), line))
			buf.WriteRune('\n')
		}
		buf.WriteString(fmt.Sprintf(format, e.StackTrace()))
		buf.WriteRune('\n')
		if !f.Flag('0') {
			for err := e.Unwrap(); err != nil; {
				if e2, ok := err.(*Erf); ok {
					buf.WriteRune('\n')
					for _, line := range strings.Split(e2.err.Error(), "\n") {
						buf.WriteString(fmt.Sprintf("%s%s", strings.Repeat(string(pad), wid+prec), line))
						buf.WriteRune('\n')
					}
					buf.WriteString(fmt.Sprintf(format, e2.StackTrace()))
					buf.WriteRune('\n')
				}
				if wErr, ok := err.(WrappedError); ok {
					err = wErr.Unwrap()
				} else {
					err = nil
				}
			}
		}
	}
	if buf.Len() > 0 {
		_, _ = f.Write(buf.Bytes())
	}
}

// Fmt returns the format argument of the formatting functions (Newf, Errorf or Wrap) that created Erf.
func (e *Erf) Fmt() string {
	return e.format
}

// Len returns the length of the arguments slice.
func (e *Erf) Len() int {
	return len(e.args)
}

// Arg returns an argument value on the given index. It panics if index is out of range.
func (e *Erf) Arg(index int) interface{} {
	if index < 0 || index >= e.Len() {
		panic("index is out of range")
	}
	return e.args[index]
}

// Args returns all argument values. It returns nil if Erf didn't create with formatting functions.
func (e *Erf) Args() []interface{} {
	if e.args == nil {
		return nil
	}
	result := make([]interface{}, len(e.args))
	copy(result, e.args)
	return result
}

// Attach attaches tags to arguments, if arguments are given. It panics if an error occurs.
func (e *Erf) Attach(tags ...string) *Erf {
	if e.args == nil {
		panic("args are not using")
	}
	if e.tagIndexes != nil {
		panic("tags are already attached")
	}
	if len(tags) > len(e.args) {
		panic("tags are more than args")
	}
	tagIndexes := make(map[string]int, len(tags))
	for index, tag := range tags {
		if tag == "" {
			continue
		}
		if _, ok := tagIndexes[tag]; ok {
			panic("tag is already defined")
		}
		tagIndexes[tag] = index
	}
	e.tagIndexes = tagIndexes
	return e
}

// Tag returns an argument value on the given tag. It panics if tag is not found.
func (e *Erf) Tag(tag string) interface{} {
	index := -1
	if idx, ok := e.tagIndexes[tag]; ok {
		index = idx
	}
	if index < 0 || index >= e.Len() {
		panic("tag is not found")
	}
	return e.args[index]
}

// PC returns program counters.
func (e *Erf) PC() []uintptr {
	result := make([]uintptr, len(e.pc))
	copy(result, e.pc)
	return result
}

// StackTrace returns a StackTrace of Erf.
func (e *Erf) StackTrace() *StackTrace {
	return NewStackTrace(e.pc...)
}

func (e *Erf) initialize(skip int) {
	e.pc = PC(int(4096/unsafe.Sizeof(uintptr(0))), skip)
}

// New creates a new Erf object with the given text.
func New(text string) *Erf {
	e := &Erf{
		err: errors.New(text),
	}
	e.initialize(4)
	return e
}

func newf(format string, args ...interface{}) *Erf {
	e := &Erf{
		err:    fmt.Errorf(format, args...),
		format: format,
		args:   make([]interface{}, len(args)),
	}
	copy(e.args, args)
	return e
}

// Newf creates a new Erf object with the given format and args.
func Newf(format string, args ...interface{}) *Erf {
	e := newf(format, args...)
	e.initialize(4)
	return e
}

// Errorf is similar with Newf except that it returns the error interface instead of the Erf pointer.
func Errorf(format string, a ...interface{}) error {
	e := newf(format, a...)
	e.initialize(4)
	return e
}

// Wrap wraps the given error as the underlying error and returns a new Erf object as the error interface.
func Wrap(err error) error {
	e := newf("%w", err)
	e.initialize(4)
	return e
}
