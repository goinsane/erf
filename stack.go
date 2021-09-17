package erf

import (
	"bytes"
	"fmt"
	"runtime"
)

// StackCaller stores the information of stack caller.
// StackCaller can format given information as string by using Format or String methods.
type StackCaller struct {
	runtime.Frame
}

// String is implementation of fmt.Stringer.
// It is synonym with fmt.Sprintf("%s", c).
func (c StackCaller) String() string {
	return fmt.Sprintf("%s", c)
}

// Format is implementation of fmt.Formatter.
//
// For '%s' (also '%v'):
// 	%s       just show function and entry (default: padding char '\t', padding 0, indent 1)
// 	%+s      show file path, line and pc
// 	%#s      use file name as file path
// 	%-s      use ' ' as padding char (padding 0, indent 2)
// 	%4s      padding 4, default indent
// 	%.3s     default padding, indent 3
// 	%4.3s    padding 4, indent 3
// 	%4.s     padding 4, indent 0
func (c StackCaller) Format(f fmt.State, verb rune) {
	buf := bytes.NewBuffer(make([]byte, 0, 4096))
	switch verb {
	case 's', 'v':
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
		padding := bytes.Repeat([]byte{pad}, wid)
		indent := bytes.Repeat([]byte{pad}, prec)
		buf.Write(padding)
		fn := "???"
		if c.Function != "" {
			fn = trimSrcPath(c.Function)
		}
		buf.WriteString(fmt.Sprintf("%s(%#x)", fn, c.Entry))
		if f.Flag('+') {
			buf.WriteRune('\n')
			buf.Write(padding)
			buf.Write(indent)
			file, line := "???", 0
			if c.File != "" {
				file = trimSrcPath(c.File)
				if f.Flag('#') {
					file = trimDirs(file)
				}
			}
			if c.Line > 0 {
				line = c.Line
			}
			buf.WriteString(fmt.Sprintf("%s:%d +%#x", file, line, c.PC-c.Entry))
		}
	default:
		return
	}
	_, _ = f.Write(buf.Bytes())
}

// StackTrace stores the information of stack trace.
type StackTrace struct {
	pc      []uintptr
	callers []StackCaller
}

// NewStackTrace creates a new StackTrace object.
func NewStackTrace(pc ...uintptr) *StackTrace {
	t := &StackTrace{
		pc:      make([]uintptr, len(pc)),
		callers: make([]StackCaller, 0, len(pc)),
	}
	copy(t.pc, pc)
	if len(t.pc) > 0 {
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
	}
	return t
}

// Duplicate duplicates the StackTrace object.
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

// String is implementation of fmt.Stringer.
// It is synonym with fmt.Sprintf("%s", t).
func (t *StackTrace) String() string {
	return fmt.Sprintf("%s", t)
}

// Format is implementation of fmt.Formatter.
// Format lists StackCaller's line by line with given format.
func (t *StackTrace) Format(f fmt.State, verb rune) {
	buf := bytes.NewBuffer(make([]byte, 0, 4096))
	switch verb {
	case 's', 'v':
		format := "%"
		for _, r := range []rune{'+', '-', '#'} {
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
	default:
		return
	}
	_, _ = f.Write(buf.Bytes())
}

// PC returns program counters.
func (t *StackTrace) PC() []uintptr {
	result := make([]uintptr, len(t.pc))
	copy(result, t.pc)
	return result
}

// Len returns the length of the StackCaller slice.
func (t *StackTrace) Len() int {
	return len(t.callers)
}

// Caller returns a StackCaller on the given index. It panics if index is out of range.
func (t *StackTrace) Caller(index int) StackCaller {
	if index < 0 || index >= t.Len() {
		panic("index out of range")
	}
	return t.callers[index]
}
