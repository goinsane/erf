// +build examples

package main

import (
	"fmt"

	"github.com/goinsane/erf"
)

func main() {
	e := erf.New("an example erf error")
	e2 := erf.Wrap(e)

	fmt.Println(e2)
}
