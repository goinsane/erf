# erf

[![Go Reference](https://pkg.go.dev/badge/github.com/goinsane/erf.svg)](https://pkg.go.dev/github.com/goinsane/erf)

Package erf provides error management with stack trace.
Erf is an error type that wraps the underlying error that stores and formats the stack trace.
Please see [godoc](https://pkg.go.dev/github.com/goinsane/erf).

## Examples

To run any example, please use the command like the following:

    cd examples/
    go run example1.go

## Tests

To run all tests, please use the following command:

    go test -v

To run all examples, please use the following command:

    go test -v -run=^Example

To run all benchmarks, please use the following command:

    go test -v -run=^Benchmark -bench=.
