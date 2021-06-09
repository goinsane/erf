package erf

type WrappedError interface {
	error
	Unwrap() error
}
