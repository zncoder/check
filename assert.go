package assert

import (
	"fmt"
	"log"
	"log/slog"
	"os"
)

func Nilf(err error, format string, args ...any) {
	if err == nil {
		return
	}
	s := fmt.Sprintf("ASSERT err:%v %s", err, fmt.Sprintf(format, args...))
	panic(s)
}

func Nil(err error) {
	if err == nil {
		return
	}
	s := fmt.Sprintf("ASSERT err:%v", err)
	panic(s)
}

func OKf(ok bool, format string, args ...any) {
	if ok {
		return
	}
	s := fmt.Sprintf("ASSERT: %s", fmt.Sprintf(format, args...))
	panic(s)
}

func OK(ok bool) {
	if ok {
		return
	}
	s := "ASSERT false"
	panic(s)
}

func Fatalf(skip int, ok bool, format string, args ...any) {
	if ok {
		return
	}
	s := fmt.Sprintf("ASSERT false: %s", fmt.Sprintf(format, args...))
	log.Output(skip+1, s)
	os.Exit(1)
}

func Must[T any](v T, err error) T {
	if err == nil {
		return v
	}
	s := fmt.Sprintf("ASSERT err:%v", err)
	panic(s)
}

func V[T any](v T, err error) valueError[T] {
	return valueError[T]{v: v, err: err}
}

type valueError[T any] struct {
	v       T
	err     error
	ignored bool
}

func (v valueError[T]) Panic(msg string, args ...any) T {
	if !v.ignored && v.err != nil {
		slog.Error(msg, append([]any{"err", v.err}, args...)...)
		panic(v.err)
	}
	return v.v
}

func (v valueError[T]) Fatal(msg string, args ...any) T {
	if !v.ignored && v.err != nil {
		slog.Error(msg, append([]any{"err", v.err}, args...)...)
		os.Exit(1)
	}
	return v.v
}

func (v valueError[T]) Ignore(ignore bool) valueError[T] {
	v.ignored = ignore
	return v
}

func (v valueError[T]) Must() T {
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

func (v valueOK[T]) Panic(msg string, args ...any) T {
	if !v.ignored && !v.ok {
		slog.Error(msg, args...)
		panic(false)
	}
	return v.v
}

func (v valueOK[T]) Fatal(msg string, args ...any) T {
	if !v.ignored && !v.ok {
		slog.Error(msg, args...)
		os.Exit(1)
	}
	return v.v
}

func (v valueOK[T]) Ignore(ignore bool) valueOK[T] {
	v.ignored = ignore
	return v
}

func (v valueOK[T]) Must() T {
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

func (v checkE) Panic(msg string, args ...any) {
	if !v.ignored && v.err != nil {
		slog.Error(msg, append([]any{"err", v.err}, args...)...)
		panic(v.err)
	}
}

func (v checkE) Fatal(msg string, args ...any) {
	if !v.ignored && v.err != nil {
		slog.Error(msg, append([]any{"err", v.err}, args...)...)
		os.Exit(1)
	}
}

func (v checkE) Ignore(ignore bool) checkE {
	v.ignored = ignore
	return v
}

func (v checkE) Must() {
	if !v.ignored && v.err != nil {
		panic(v.err)
	}
}
