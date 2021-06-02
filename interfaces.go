package erf

type UnwrapIfc interface {
	Unwrap() error
}

type WrappedErrorIfc interface {
	error
	UnwrapIfc
}
