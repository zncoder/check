package check

import (
	"log/slog"
	"os"
)

func Vp[T any](v T, err error) T {
	return V(v, err).V()
}

func Tp(ok bool) {
	T(ok).M()
}

func V[T any](v T, err error) valueError[T] {
	return valueError[T]{v: v, err: err}
}

type valueError[T any] struct {
	v       T
	err     error
	ignored bool
}

func (v valueError[T]) P(msg string, args ...any) T {
	if !v.ignored && v.err != nil {
		slog.Error(msg, append([]any{"err", v.err}, args...)...)
		panic(v.err)
	}
	return v.v
}

func (v valueError[T]) F(msg string, args ...any) T {
	if !v.ignored && v.err != nil {
		slog.Error(msg, append([]any{"err", v.err}, args...)...)
		os.Exit(1)
	}
	return v.v
}

func (v valueError[T]) I(ignore bool) valueError[T] {
	v.ignored = ignore
	return v
}

func (v valueError[T]) V() T {
	if !v.ignored && v.err != nil {
		panic(v.err)
	}
	return v.v
}

func K[T any](v T, ok bool) valueOK[T] {
	return valueOK[T]{v: v, ok: ok}
}

type valueOK[T any] struct {
	v       T
	ok      bool
	ignored bool
}

func (v valueOK[T]) P(msg string, args ...any) T {
	if !v.ignored && !v.ok {
		slog.Error(msg, args...)
		panic(false)
	}
	return v.v
}

func (v valueOK[T]) F(msg string, args ...any) T {
	if !v.ignored && !v.ok {
		slog.Error(msg, args...)
		os.Exit(1)
	}
	return v.v
}

func (v valueOK[T]) I(ignore bool) valueOK[T] {
	v.ignored = ignore
	return v
}

func (v valueOK[T]) V() T {
	if !v.ignored && !v.ok {
		panic(false)
	}
	return v.v
}

func E(err error) checkE {
	return checkE{err: err}
}

type checkE struct {
	err     error
	ignored bool
}

func (v checkE) P(msg string, args ...any) {
	if !v.ignored && v.err != nil {
		slog.Error(msg, append([]any{"err", v.err}, args...)...)
		panic(v.err)
	}
}

func (v checkE) F(msg string, args ...any) {
	if !v.ignored && v.err != nil {
		slog.Error(msg, append([]any{"err", v.err}, args...)...)
		os.Exit(1)
	}
}

func (v checkE) I(ignore bool) checkE {
	v.ignored = ignore
	return v
}

func (v checkE) M() {
	if !v.ignored && v.err != nil {
		panic(v.err)
	}
}

func T(ok bool) checkT {
	return checkT{ok: ok}
}

type checkT struct {
	ok      bool
	ignored bool
}

func (v checkT) P(msg string, args ...any) {
	if !v.ignored && !v.ok {
		slog.Error(msg, args...)
		panic(false)
	}
}

func (v checkT) F(msg string, args ...any) {
	if !v.ignored && !v.ok {
		slog.Error(msg, args...)
		os.Exit(1)
	}
}

func (v checkT) I(ignore bool) checkT {
	v.ignored = ignore
	return v
}

func (v checkT) M() {
	if !v.ignored && !v.ok {
		panic(false)
	}
}
