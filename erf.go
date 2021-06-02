package erf

import (
	"bytes"
	"errors"
	"fmt"
	"unsafe"
)

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

// Unwrap returns underlying error.
func (e *Erf) Unwrap() error {
	if err, ok := e.err.(WrappedErrorIfc); ok {
		return err.Unwrap()
	}
	return nil
}

// String is implementation of fmt.Stringer.
func (e *Erf) String() string {
	return fmt.Sprintf("%s", e)
}

// Format is implementation of fmt.Formatter.
func (e *Erf) Format(f fmt.State, verb rune) {
	buf := bytes.NewBuffer(nil)
	switch verb {
	case 's', 'v':
		buf.WriteString(e.err.Error())
		if f.Flag('+') {
			format := "%+"
			for _, r := range []rune{'_', '#', ' '} {
				if f.Flag(int(r)) {
					format += string(r)
				}
			}
			w, _ := f.Width()
			w++
			if f.Flag(' ') {
				w++
			}
			format += fmt.Sprintf("%d", w)
			if p, ok := f.Precision(); ok {
				format += fmt.Sprintf(".%d", p)
			}
			format += "s\n"
			buf.WriteRune('\n')
			buf.WriteString(fmt.Sprintf(format, NewStackTrace(e.pc...)))
			if !f.Flag('0') {
				for err := e.Unwrap(); err != nil; {
					if e2, ok := err.(*Erf); ok {
						buf.WriteRune('\n')
						buf.WriteString(e2.Error())
						buf.WriteRune('\n')
						buf.WriteString(fmt.Sprintf(format, NewStackTrace(e2.pc...)))
					}
					if wErr, ok := err.(WrappedErrorIfc); ok {
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

func (e *Erf) initialize() {
	e.pc = getPC(int(4096/unsafe.Sizeof(uintptr(0))), 4)
}

func New(text string) *Erf {
	e := &Erf{
		err: errors.New(text),
	}
	e.initialize()
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

func Newf(format string, args ...interface{}) *Erf {
	e := newf(format, args...)
	e.initialize()
	return e
}

func Errorf(format string, a ...interface{}) error {
	e := newf(format, a...)
	e.initialize()
	return e
}

func Wrap(err error) error {
	e := newf("%w", err)
	e.initialize()
	return e
}
