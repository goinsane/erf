// +build examples

package main

import (
	"fmt"

	"github.com/goinsane/erf"
)

var (
	ErrValueBelowZero = erf.New("value below zero")
)

type InvalidArgumentError struct{ *erf.Erf }

func NewInvalidArgumentError(name string, err error) error {
	return (&InvalidArgumentError{erf.Newf("invalid argument %q: %w", name, err)}).Attach("name")
}

func Foo(x int) error {
	if x < 0 {
		return NewInvalidArgumentError("x", ErrValueBelowZero)
	}
	return nil
}

func Bar(y int) error {
	if y < 0 {
		return erf.Errorf("%w: %d", ErrValueBelowZero, y)
	}
	return nil
}

func Baz(z int) error {
	if z < 0 {
		return erf.Wrap(ErrValueBelowZero)
	}
	return nil
}

func main() {
	fmt.Println("#### Foo: just the show first error message without padding and indent")
	if err := Foo(-1); err != nil {
		fmt.Printf("%v\n", err)
	}

	fmt.Println("#### Foo: show with stack trace")
	if err := Foo(-2); err != nil {
		fmt.Printf("%x\n", err)
	}

	fmt.Println("#### Foo: show with stack trace, and use space as whitespace instead of tab")
	if err := Foo(-3); err != nil {
		fmt.Printf("% x\n", err)
	}

	fmt.Println("#### Foo: show with first stack trace")
	if err := Foo(-4); err != nil {
		fmt.Printf("%X\n", err)
	}

	fmt.Println("#### Foo: show with stack trace with only file names except full path")
	if err := Foo(-5); err != nil {
		fmt.Printf("%#x\n", err)
	}

	fmt.Println("#### Foo: show with stack trace (padding 2, indent 1)")
	if err := Foo(-6); err != nil {
		fmt.Printf("%2.1x\n", err)
	}

	fmt.Println("#### Bar: show with stack trace")
	if err := Bar(-7); err != nil {
		fmt.Printf("%x\n", err)
	}

	fmt.Println("#### Bar: show wrapped error with stack trace (padding 2, indent 1)")
	if err := Foo(-8); err != nil {
		err = erf.Wrap(err)
		fmt.Printf("%2.1x\n", err)
	}

	fmt.Println("#### Baz: show with stack trace")
	if err := Baz(-9); err != nil {
		fmt.Printf("%x\n", err)
	}

	fmt.Println("#### Baz: just show all of stack traces of errors")
	if err := Baz(-10); err != nil {
		fmt.Printf("%-x\n", err)
	}

	fmt.Println("#### Baz: just show the stack trace of the first error")
	if err := Baz(-11); err != nil {
		fmt.Printf("%-X\n", err)
	}
}
