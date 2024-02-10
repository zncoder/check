package check

import (
	"log/slog"
	"os"
)

type logAction int

const (
	toPanic logAction = 1 + iota
	toFatal
	toIgnore
	toLog
)

func (la logAction) String() string {
	switch la {
	case toPanic:
		return "panic"
	case toFatal:
		return "fatal"
	case toIgnore:
		return "ignore"
	case toLog:
		return "info"
	default:
		panic(la)
	}
}

func logErr(la logAction, err error, args []any) {
	if len(args) > 0 {
		slog.Error(args[0].(string), append([]any{"err", err}, args[1:]...)...)
	} else if la != toPanic {
		slog.Error(la.String(), "err", err)
	}
	switch la {
	case toPanic:
		panic(err)
	case toFatal:
		os.Exit(1)
	}
}

func logFalse(la logAction, args []any) {
	if len(args) > 0 {
		slog.Error(args[0].(string), args[1:]...)
	} else if la != toPanic {
		slog.Error(la.String() + " FALSE")
	}
	switch la {
	case toPanic:
		panic(false)
	case toFatal:
		os.Exit(1)
	}
}

func L(msg string, args ...any) {
	slog.Info(msg, args...)
}

func P(msg string, args ...any) {
	slog.Error(msg, args...)
	panic(msg)
}

func F(msg string, args ...any) {
	slog.Error(msg, args...)
	os.Exit(1)
}

// V captures a value and an error. Usage:
//
//	V(os.CreateFile(filename, 0600)).P("open file", "filename", filename), or
//	V(os.CreateFile(filename, 0600)).P(), or
//	V(os.CreateFile(filename, 0600)).F("open file", "filename", filename)
//	V(os.CreateFile(filename, 0600)).F()
func V[T any](v T, err error) valueError[T] {
	return valueError[T]{v: v, err: err}
}

type valueError[T any] struct {
	v      T
	err    error
	silent bool
}

// P panics on error. If args is set, args[0] must be a string,
// and the rest args are args to slog.Error.
func (v valueError[T]) P(args ...any) T {
	if !v.silent && v.err != nil {
		logErr(toPanic, v.err, args)
	}
	return v.v
}

// F exits on error.
func (v valueError[T]) F(args ...any) T {
	if !v.silent && v.err != nil {
		logErr(toFatal, v.err, args)
	}
	return v.v
}

// L logs on error.
func (v valueError[T]) L(args ...any) T {
	if !v.silent && v.err != nil {
		logErr(toLog, v.err, args)
	}
	return v.v
}

// S ignores error.
func (v valueError[T]) S(silent bool) valueError[T] {
	v.silent = silent
	return v
}

// K captures a value and a bool.
func K[T any](v T, ok bool) valueOK[T] {
	return valueOK[T]{v: v, ok: ok}
}

type valueOK[T any] struct {
	v      T
	ok     bool
	silent bool
}

func (v valueOK[T]) P(args ...any) T {
	if !v.silent && !v.ok {
		logFalse(toPanic, args)
	}
	return v.v
}

func (v valueOK[T]) F(args ...any) T {
	if !v.silent && !v.ok {
		logFalse(toFatal, args)
	}
	return v.v
}

func (v valueOK[T]) L(args ...any) T {
	if !v.silent && !v.ok {
		logFalse(toLog, args)
	}
	return v.v
}

func (v valueOK[T]) S(silent bool) valueOK[T] {
	v.silent = silent
	return v
}

// E captures an error.
func E(err error) checkE {
	return checkE{err: err}
}

type checkE struct {
	err    error
	silent bool
}

func (v checkE) P(args ...any) {
	if !v.silent && v.err != nil {
		logErr(toPanic, v.err, args)
	}
}

func (v checkE) F(args ...any) {
	if !v.silent && v.err != nil {
		logErr(toFatal, v.err, args)
	}
}

func (v checkE) L(args ...any) {
	if !v.silent && v.err != nil {
		logErr(toLog, v.err, args)
	}
}

func (v checkE) S(silent bool) checkE {
	v.silent = silent
	return v
}

// T captures a bool.
func T(ok bool) checkT {
	return checkT{ok: ok}
}

type checkT struct {
	ok     bool
	silent bool
}

func (v checkT) P(args ...any) {
	if !v.silent && !v.ok {
		logFalse(toPanic, args)
	}
}

func (v checkT) F(args ...any) {
	if !v.silent && !v.ok {
		logFalse(toFatal, args)
	}
}

func (v checkT) L(args ...any) {
	if !v.silent && !v.ok {
		logFalse(toLog, args)
	}
}

func (v checkT) S(silent bool) checkT {
	v.silent = silent
	return v
}
