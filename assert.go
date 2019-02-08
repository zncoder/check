package assert

import (
	"fmt"
	"log"
	"os"
)

func Nilf(err error, format string, args ...interface{}) {
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

func OKf(ok bool, format string, args ...interface{}) {
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

func Fatalf(skip int, ok bool, format string, args ...interface{}) {
	if ok {
		return
	}
	s := fmt.Sprintf("ASSERT false: %s", fmt.Sprintf(format, args...))
	log.Output(skip+1, s)
	os.Exit(1)
}
