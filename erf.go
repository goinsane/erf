package erf

import (
	"bytes"
	"errors"
	"fmt"
	"unsafe"
)

// Erf is an error type that wraps the underlying error that stores and formats the stack trace.
type Erf struct {
	err    error
	format string
	args   []interface{}
	pc     []uintptr
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
		buf.WriteString(e.err.Error())
		if f.Flag('+') {
			format := "%+"
			for _, r := range []rune{'-', '#'} {
				if f.Flag(int(r)) {
					format += string(r)
				}
			}
			wid, prec := 0, 1
			if f.Flag('-') {
				prec = 2
			}
			if w, ok := f.Width(); ok {
				wid = w
			}
			if p, ok := f.Precision(); ok {
				prec = p
			}
			wid += prec
			format += fmt.Sprintf("%d.%d", wid, prec)
			format += "s\n"
			buf.WriteRune('\n')
			buf.WriteString(fmt.Sprintf(format, e.StackTrace()))
			if !f.Flag('0') {
				for err := e.Unwrap(); err != nil; {
					if e2, ok := err.(*Erf); ok {
						buf.WriteRune('\n')
						buf.WriteString(e2.Error())
						buf.WriteRune('\n')
						buf.WriteString(fmt.Sprintf(format, e2.StackTrace()))
					}
					if wErr, ok := err.(WrappedError); ok {
						err = wErr.Unwrap()
					} else {
						err = nil
					}
				}
			}
		}
	}
	if buf.Len() > 0 {
		_, _ = f.Write(buf.Bytes())
	}
}

// StackTrace returns a StackTrace of Erf.
func (e *Erf) StackTrace() *StackTrace {
	return NewStackTrace(e.pc...)
}

// Fmt returns the format argument of the formatting functions (Newf, Errorf or Wrap) that created Erf.
func (e *Erf) Fmt() string {
	return e.format
}

// Len returns the length of arguments.
func (e *Erf) Len() int {
	return len(e.args)
}

// Arg returns an argument on the given index. It panics if index is out of range.
func (e *Erf) Arg(index int) interface{} {
	if index < 0 || index >= e.Len() {
		panic("out of range")
	}
	return e.args[index]
}

// Args returns all arguments. It returns nil if Erf didn't create with formatting functions.
func (e *Erf) Args() []interface{} {
	if e.args == nil {
		return nil
	}
	result := make([]interface{}, len(e.args))
	copy(result, e.args)
	return result
}

func (e *Erf) initialize(skip int) {
	e.pc = getPC(int(4096/unsafe.Sizeof(uintptr(0))), skip)
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

// Errorf is synonymous with Newf except that it returns the error interface instead of the Erf pointer.
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
