package erf_test

import (
	"fmt"

	"github.com/goinsane/erf"
)

func ExampleErf() {
	e := erf.New("an example erf error")
	e2 := erf.Wrap(e)

	fmt.Println("just show first error message (default: padding char '\\t', padding 0, indent 1)")
	fmt.Printf("%s\n\n", e2)

	fmt.Println("show error message with indent, append stack trace using format '%+s'")
	fmt.Printf("%+s\n\n", e2)

	fmt.Println("show first error message with indent, append stack trace using format '%+s'")
	fmt.Printf("%+ s\n\n", e2)

	fmt.Println("use file name as file path for StackCaller")
	fmt.Printf("%+#s\n\n", e2)

	fmt.Println("padding 2, indent 1 by default")
	fmt.Printf("%+#2s\n\n", e2)

	fmt.Println("padding 2, indent 3")
	fmt.Printf("%+#2.3s\n\n", e2)

	fmt.Println("use ' ' as padding char (padding 0, indent 2)")
	fmt.Printf("%-s\n\n", e2)

	fmt.Println("show error message with indent, append stack trace using format '%+s'")
	fmt.Printf("%-+s\n\n", e2)

	fmt.Println("show first error message with indent, append stack trace using format '%+s'")
	fmt.Printf("%-+ s\n\n", e2)

	fmt.Println("use file name as file path for StackCaller")
	fmt.Printf("%-+#s\n\n", e2)

	fmt.Println("padding 2, indent 2 by default")
	fmt.Printf("%-+#2s\n\n", e2)

	fmt.Println("padding 2, indent 3")
	fmt.Printf("%-+#2.3s\n\n", e2)
}
