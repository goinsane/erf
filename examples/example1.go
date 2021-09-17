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
	fmt.Println("#### Foo: just show error text")
	if err := Foo(-1); err != nil {
		fmt.Printf("%v\n", err)
	}

	fmt.Println("#### Foo: show with stack trace")
	if err := Foo(-2); err != nil {
		fmt.Printf("%+v\n", err)
	}

	fmt.Println("#### Foo: show with stack trace, and use space as whitespace instead of tab")
	if err := Foo(-3); err != nil {
		fmt.Printf("%+-v\n", err)
	}

	fmt.Println("#### Foo: show with last stack trace")
	if err := Foo(-4); err != nil {
		fmt.Printf("%+0v\n", err)
	}

	fmt.Println("#### Foo: show with stack trace with only file names except full path")
	if err := Foo(-5); err != nil {
		fmt.Printf("%+#v\n", err)
	}

	fmt.Println("#### Foo: show with stack trace with 2 whitespace chars of padding and 1 whitespace char of indent")
	if err := Foo(-6); err != nil {
		fmt.Printf("%+2.1v\n", err)
	}

	fmt.Println("#### Bar: show with stack trace")
	if err := Bar(-7); err != nil {
		fmt.Printf("%+v\n", err)
	}

	fmt.Println("#### Baz: show with stack trace")
	if err := Baz(-8); err != nil {
		fmt.Printf("%+v\n", err)
	}
}
