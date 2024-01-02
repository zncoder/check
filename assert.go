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
	log.Output(2, s)
	panic(s)
}

func Nil(err error) {
	if err == nil {
		return
	}
	s := fmt.Sprintf("ASSERT err:%v", err)
	log.Output(2, s)
	panic(s)
}

func OKf(ok bool, format string, args ...any) {
	if ok {
		return
	}
	s := fmt.Sprintf("ASSERT: %s", fmt.Sprintf(format, args...))
	log.Output(2, s)
	panic(s)
}

func OK(ok bool) {
	if ok {
		return
	}
	s := "ASSERT false"
	log.Output(2, s)
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
	log.Output(2, s)
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
		slog.Info(msg, append(args, "err", v.err))
		panic(v.err)
	}
	return v.v
}

func (v valueError[T]) Fatal(msg string, args ...any) T {
	if !v.ignored && v.err != nil {
		slog.Info(msg, append(args, "err", v.err))
		os.Exit(1)
	}
	return v.v
}

func (v valueError[T]) Ignore(ignore bool) valueError[T] {
	v.ignored = ignore
	return v
}
