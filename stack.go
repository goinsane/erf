package erf

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
)

type StackCaller struct {
	runtime.Frame
}

func (c StackCaller) String() string {
	return fmt.Sprintf("%s", c)
}

func (c StackCaller) Format(f fmt.State, verb rune) {
	buf := bytes.NewBuffer(nil)
	switch verb {
	case 's', 'v':
		pad, wid, prec := byte('\t'), 0, 1
		if f.Flag(' ') {
			pad = ' '
			prec = 2
		}
		if w, ok := f.Width(); ok {
			wid = w
		}
		if p, ok := f.Precision(); ok {
			prec = p
		}
		padding := bytes.Repeat([]byte{pad}, wid)
		var str string
		buf.Write(padding)
		str = "???"
		if c.Function != "" {
			str = trimSrcPath(c.Function)
		}
		buf.WriteString(fmt.Sprintf("%s(%#x)", str, c.Entry))
		if f.Flag('+') {
			buf.WriteRune('\n')
			buf.Write(padding)
			str = trimSrcPath(c.File)
			if f.Flag('#') {
				str = trimDirs(str)
			}
			buf.WriteString(fmt.Sprintf("%s%s:%d +%#x", strings.Repeat(string(pad), prec), str, c.Line, c.PC-c.Entry))
		}
	}
	if buf.Len() > 0 {
		_, _ = f.Write(buf.Bytes())
	}
}

type StackTrace struct {
	pc      []uintptr
	callers []StackCaller
}

func NewStackTrace(pc ...uintptr) *StackTrace {
	t := &StackTrace{
		pc:      make([]uintptr, len(pc)),
		callers: make([]StackCaller, 0, len(pc)),
	}
	copy(t.pc, pc)
	frames := runtime.CallersFrames(t.pc)
	for {
		frame, more := frames.Next()
		caller := StackCaller{
			Frame: frame,
		}
		t.callers = append(t.callers, caller)
		if !more {
			break
		}
	}
	return t
}

func (t *StackTrace) Duplicate() *StackTrace {
	if t == nil {
		return nil
	}
	t2 := &StackTrace{
		pc:      make([]uintptr, len(t.pc), cap(t.pc)),
		callers: make([]StackCaller, len(t.callers), cap(t.callers)),
	}
	copy(t2.pc, t.pc)
	copy(t2.callers, t.callers)
	return t2
}

func (t *StackTrace) Len() int {
	return len(t.callers)
}

func (t *StackTrace) Caller(index int) StackCaller {
	if index < 0 || index >= t.Len() {
		panic("out of range")
	}
	return t.callers[index]
}

func (t *StackTrace) String() string {
	return fmt.Sprintf("%s", t)
}

func (t *StackTrace) Format(f fmt.State, verb rune) {
	buf := bytes.NewBuffer(nil)
	switch verb {
	case 's', 'v':
		format := "%"
		for _, r := range []rune{'+', '_', '#', ' ', '0'} {
			if f.Flag(int(r)) {
				format += string(r)
			}
		}
		if w, ok := f.Width(); ok {
			format += fmt.Sprintf("%d", w)
		}
		if p, ok := f.Precision(); ok {
			format += fmt.Sprintf(".%d", p)
		}
		format += "s"
		for i, c := range t.callers {
			if i > 0 {
				buf.WriteRune('\n')
			}
			buf.WriteString(fmt.Sprintf(format, c))
		}
	}
	if buf.Len() > 0 {
		_, _ = f.Write(buf.Bytes())
	}
}
