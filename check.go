package check

import (
	"errors"
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
func V[T any](v T, err error) valueE[T] {
	return valueE[T]{v: v, err: err}
}

type valueE[T any] struct {
	v      T
	err    error
	silent bool
}

// P panics on error. If args is set, args[0] must be a string,
// and the rest args are args to slog.Error.
func (v valueE[T]) P(args ...any) T {
	if !v.silent && v.err != nil {
		logErr(toPanic, v.err, args)
	}
	return v.v
}

// F exits on error.
func (v valueE[T]) F(args ...any) T {
	if !v.silent && v.err != nil {
		logErr(toFatal, v.err, args)
	}
	return v.v
}

// L logs on error.
func (v valueE[T]) L(args ...any) bool {
	_, ok := v.K(args...)
	return ok
}

// K returns value and ok
func (v valueE[T]) K(args ...any) (T, bool) {
	if !v.silent && v.err != nil {
		logErr(toLog, v.err, args)
	}
	return v.v, v.err == nil
}

// S ignores error.
func (v valueE[T]) S() valueE[T] {
	v.silent = true
	return v
}

// I ignores error if it is one of errs.
func (v valueE[T]) I(errs ...error) valueE[T] {
	if v.err == nil {
		return v
	}
	for _, err := range errs {
		if errors.Is(v.err, err) {
			v.err = nil
			return v
		}
	}
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

func (v valueOK[T]) L(args ...any) bool {
	_, ok := v.K(args...)
	return ok
}

func (v valueOK[T]) K(args ...any) (T, bool) {
	if !v.silent && !v.ok {
		logFalse(toLog, args)
	}
	return v.v, v.ok
}

func (v valueOK[T]) S() valueOK[T] {
	v.silent = true
	return v
}

// E captures an error.
func E(err error) checkE {
	return checkE{err: err}
}

type checkE struct {
	silent bool
	err    error
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

func (v checkE) L(args ...any) bool {
	if !v.silent && v.err != nil {
		logErr(toLog, v.err, args)
	}
	return v.err == nil
}

func (v checkE) S() checkE {
	v.silent = true
	return v
}

func (v checkE) I(errs ...error) checkE {
	if v.err == nil {
		return v
	}
	for _, err := range errs {
		if errors.Is(v.err, err) {
			v.err = nil
			return v
		}
	}
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

func (v checkT) L(args ...any) bool {
	if !v.silent && !v.ok {
		logFalse(toLog, args)
	}
	return v.ok
}

func (v checkT) S() checkT {
	v.silent = true
	return v
}
