package erf_test

import (
	"fmt"

	"github.com/goinsane/erf"
)

func ExampleErf() {
	e := erf.New("an example erf error")
	err := erf.Errorf("we have an example error: %w", e)

	fmt.Println("just show the first error message without padding and indent.")
	fmt.Printf("%v\n\n", err)

	fmt.Println("list all error messages by using indent and show StackTrace of errors by using format '%+s'.")
	fmt.Printf("%x\n\n", err)

	fmt.Println("list all error messages by using indent and show StackTrace of errors by using format '% s'.")
	fmt.Printf("% x\n\n", err)

	fmt.Println("list all error messages by using indent and show StackTrace of errors by using format '%#s'.")
	fmt.Printf("%#x\n\n", err)

	fmt.Println("list all error messages by using indent and show StackTrace of errors by using format '% #s'.")
	fmt.Printf("% #x\n\n", err)

	fmt.Println("show the first error message by using indent and show the StackTrace of error by using format '%+s'.")
	fmt.Printf("%X\n\n", err)

	fmt.Println("show the first error message by using indent and show the StackTrace of error by using format '% s'.")
	fmt.Printf("% X\n\n", err)

	fmt.Println("show the first error message by using indent and show the StackTrace of error by using format '%#s'.")
	fmt.Printf("%#X\n\n", err)

	fmt.Println("show the first error message by using indent and show the StackTrace of error by using format '% #s'.")
	fmt.Printf("% #X\n\n", err)

	fmt.Println("don't show any error messages, just show all of StackTrace of errors.")
	fmt.Printf("%-x\n\n", err)

	fmt.Println("don't show the error message, just show the StackTrace of the first error.")
	fmt.Printf("%-X\n\n", err)

	fmt.Println("padding 2, indent 1 by default. padding char '\\t'.")
	fmt.Printf("%2x\n\n", err)

	fmt.Println("padding 2, indent 3. padding char '\\t'.")
	fmt.Printf("%2.3x\n\n", err)

	fmt.Println("padding 0 by default, indent 3. padding char '\\t'.")
	fmt.Printf("%.3x\n\n", err)

	fmt.Println("padding 2, indent 2 by default. padding char ' '.")
	fmt.Printf("% 2x\n\n", err)

	fmt.Println("padding 2, indent 3. padding char ' '.")
	fmt.Printf("% 2.3x\n\n", err)

	fmt.Println("padding 0 by default, indent 3. padding char ' '.")
	fmt.Printf("% .3x\n\n", err)
}
