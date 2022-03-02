// Package erf provides error management with stack trace.
package erf

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
	"unsafe"
)

var (
	// DefaultPCSize defines max length of Erf program counters in PC() method.
	DefaultPCSize = int(uintptr(os.Getpagesize()) / unsafe.Sizeof(uintptr(0)))
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

// UnwrapAll returns all errors using Unwrap method. The first element in the returned value is e.
func (e *Erf) UnwrapAll() []error {
	result := make([]error, 0, 4096)
	for err := error(e); err != nil; {
		result = append(result, err)
		if wErr, ok := err.(WrappedError); ok {
			err = wErr.Unwrap()
		} else {
			err = nil
		}
	}
	return result
}

// Format is implementation of fmt.Formatter.
// Format lists error messages and appends StackTrace's for underlying Erf and all of wrapped Erf's,
// line by line with given format.
//
// For '%v' (also '%s'):
// 	%v       just show the first error message without padding and indent.
//
// For '%x' and '%X':
// 	%x       list all error messages by using indent and show StackTrace of errors by using format '%+s'.
// 	%+x      similar with '%x', also shows tags.
// 	% x      list all error messages by using indent and show StackTrace of errors by using format '% s'.
// 	%#x      list all error messages by using indent and show StackTrace of errors by using format '%#s'.
// 	% #x     list all error messages by using indent and show StackTrace of errors by using format '% #s'.
// 	%X       show the first error message by using indent and show the StackTrace of error by using format '%+s'.
// 	%+X      similar with '%X', also shows tags.
// 	% X      show the first error message by using indent and show the StackTrace of error by using format '% s'.
// 	%#X      show the first error message by using indent and show the StackTrace of error by using format '%#s'.
// 	% #X     show the first error message by using indent and show the StackTrace of error by using format '% #s'.
// 	%-x      don't show any error messages, just show all of StackTrace of errors.
// 	%-X      don't show the error message, just show the StackTrace of the first error.
// 	%4x      same with '%x', padding 4, indent 1 by default.
// 	%.3x     same with '%x', padding 0 by default, indent 3.
// 	%4.3x    same with '%x', padding 4, indent 3.
// 	%4.x     same with '%x', padding 4, indent 0.
// 	% 4x     same with '% x', padding 4, indent 2 by default.
// 	% .3x    same with '% x', padding 0 by default, indent 3.
// 	% 4.3x   same with '% x', padding 4, indent 3.
// 	% 4.x    same with '% x', padding 4, indent 0.
// 	%#4.3x   same with '%#x', padding 4, indent 3.
// 	% #4.3x  same with '% #x', padding 4, indent 3.
func (e *Erf) Format(f fmt.State, verb rune) {
	buf := bytes.NewBuffer(make([]byte, 0, 4096))
	switch verb {
	case 's', 'v':
		buf.WriteString(e.err.Error())
	case 'x', 'X':
		format := "%+"
		for _, r := range []rune{' ', '#'} {
			if f.Flag(int(r)) {
				format += string(r)
			}
		}
		pad, wid, prec := getPadWidPrec(f)
		format += fmt.Sprintf("%d.%ds", wid, prec)
		padding, indent := bytes.Repeat([]byte{pad}, wid), bytes.Repeat([]byte{pad}, prec)
		for idx, err := range e.UnwrapAll() {
			if idx > 0 {
				buf.WriteRune('\n')
			}
			if e2, ok := err.(*Erf); ok {
				if !f.Flag('-') {
					for _, line := range strings.Split(e2.Error(), "\n") {
						buf.Write(padding)
						buf.Write(indent)
						buf.WriteString(line)
						buf.WriteRune('\n')
					}
				}
				str := fmt.Sprintf(format, e2.StackTrace())
				if str != "" {
					buf.WriteString(str)
				} else {
					buf.WriteString("* ")
				}
				buf.WriteRune('\n')
				if f.Flag('+') {
					tags := e2.Tags()
					if len(tags) > 0 {
						buf.Write(padding)
						buf.WriteString("+ ")
						for idx, tag := range tags {
							if idx > 0 {
								buf.WriteRune(' ')
							}
							buf.WriteString(fmt.Sprintf("%q=%q", tag, fmt.Sprintf("%v", e2.Tag(tag))))
						}
						buf.WriteRune('\n')
					}
				}
			} else {
				if !f.Flag('-') {
					for _, line := range strings.Split(err.Error(), "\n") {
						buf.Write(padding)
						buf.Write(indent)
						buf.WriteString(line)
						buf.WriteRune('\n')
					}
				} else {
					buf.Write(padding)
					buf.WriteString("- ")
					buf.WriteRune('\n')
				}
			}
			buf.Write(padding)
			if verb == 'X' {
				break
			}
		}
	default:
		return
	}
	_, _ = f.Write(buf.Bytes())
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
	if index < 0 || index >= len(e.args) {
		panic("index out of range")
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

// Attach attaches tags to arguments, if arguments are given.
// If tag is "", it passes attaching tag to corresponding argument.
// It panics for given errors:
// 	args are not using
// 	tags are already attached
// 	number of tags is more than args
// 	tag already defined
func (e *Erf) Attach(tags ...string) *Erf {
	if e.args == nil {
		panic("args are not using")
	}
	if e.tagIndexes != nil {
		panic("tags are already attached")
	}
	if len(tags) > len(e.args) {
		panic("number of tags is more than args")
	}
	tagIndexes := make(map[string]int, len(tags))
	for index, tag := range tags {
		if tag == "" {
			continue
		}
		if _, ok := tagIndexes[tag]; ok {
			panic("tag already defined")
		}
		tagIndexes[tag] = index
	}
	e.tagIndexes = tagIndexes
	return e
}

// Attach2 is similar with Attach except that it returns the error interface instead of the Erf pointer.
func (e *Erf) Attach2(tags ...string) error {
	return e.Attach(tags...)
}

// Tag returns an argument value on the given tag. It returns nil if tag is not found.
func (e *Erf) Tag(tag string) interface{} {
	index := -1
	if idx, ok := e.tagIndexes[tag]; ok {
		index = idx
	}
	if index < 0 || index >= len(e.args) {
		return nil
	}
	return e.args[index]
}

// Tags returns all tags sequentially. It returns nil if tags are not attached.
func (e *Erf) Tags() []string {
	if e.tagIndexes == nil {
		return nil
	}
	m := make(map[int]string, len(e.tagIndexes))
	for tag, index := range e.tagIndexes {
		m[index] = tag
	}
	result := make([]string, 0, len(m))
	for i, j := 0, len(m); i < j; i++ {
		if tag, ok := m[i]; ok {
			result = append(result, tag)
		}
	}
	return result
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
	e.pc = PC(DefaultPCSize, skip)
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
		args:   make([]interface{}, 0, len(args)),
	}
	for _, arg := range args {
		if arg == nil {
			panic("arg is nil")
		}
		e.args = append(e.args, arg)
	}
	return e
}

// Newf creates a new Erf object with the given format and args.
// It panics if an any arg is nil.
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
// Wrap is similar with Newf("%w", err) except that it returns nil if err is nil.
func Wrap(err error) error {
	if err == nil {
		return nil
	}
	e := newf("%w", err)
	e.initialize(4)
	return e
}

// Wrapp wraps the error in the given pointer and returns a new Erf object onto the given pointer.
// Wrapp is similar with Newf("%w", err) except that it returns to perr and doesn't affect if perr or *perr is nil.
func Wrapp(perr *error) {
	if perr == nil {
		return
	}
	err := *perr
	if err == nil {
		return
	}
	e := newf("%w", err)
	e.initialize(4)
	*perr = e
}
