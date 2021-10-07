package erf

// WrappedError is an interface to simulate GoLang's wrapped errors.
type WrappedError interface {
	error
	Unwrap() error
}
