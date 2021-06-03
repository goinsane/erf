package erf

type Wrapped interface {
	Unwrap() error
}

type WrappedError interface {
	error
	Wrapped
}
