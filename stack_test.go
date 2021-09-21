package erf_test

import (
	"fmt"

	"github.com/goinsane/erf"
)

func ExampleStackCaller() {
	e := erf.New("an example erf error")
	sc := e.StackTrace().Caller(0)

	fmt.Println("just show function and entry without padding and indent.")
	fmt.Printf("%s\n\n", sc)

	fmt.Println("show file path, line and pc. padding char '\\t', default padding 0, default indent 1.")
	fmt.Printf("%+s\n\n", sc)

	fmt.Println("use file name as file path.")
	fmt.Printf("%#s\n\n", sc)

	fmt.Println("padding 2, indent 1 by default.")
	fmt.Printf("%#2s\n\n", sc)

	fmt.Println("padding 2, indent 3.")
	fmt.Printf("%#2.3s\n\n", sc)

	fmt.Println("show file path, line and pc. padding char ' ', default padding 0, default indent 2.")
	fmt.Printf("%-s\n\n", sc)

	fmt.Println("use file name as file path.")
	fmt.Printf("%-#s\n\n", sc)

	fmt.Println("padding 2, indent 2 by default.")
	fmt.Printf("%-#2s\n\n", sc)

	fmt.Println("padding 2, indent 3.")
	fmt.Printf("%-#2.3s\n\n", sc)
}

func ExampleStackTrace() {
	e := erf.New("an example erf error")
	st := e.StackTrace()

	fmt.Println("default")
	fmt.Printf("%s\n\n", st)

	fmt.Println("show file path, line and pc. padding char '\\t', default padding 0, default indent 1.")
	fmt.Printf("%+s\n\n", st)

	fmt.Println("show file path, line and pc. padding char ' ', default padding 0, default indent 2.")
	fmt.Printf("%+s\n\n", st)
}
