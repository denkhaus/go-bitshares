package util

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
)

func ToBytes(in interface{}) []byte {
	if i, ok := in.(string); ok {
		return []byte(i)
	}

	b, err := json.Marshal(in)
	if err != nil {
		panic("toBytes is unable to marshal input")
	}
	return b
}

func Dump(descr string, in interface{}) {
	fmt.Printf("%s ------------------------- dump start ---------------------------------------\n", descr)
	spew.Dump(in)
	fmt.Printf("%s -------------------------  dump end  ---------------------------------------\n\n", descr)
}

func SafeUnquote(in string) (string, error) {
	if strings.HasPrefix(in, "\"") && strings.HasSuffix(in, "\"") {
		q, err := strconv.Unquote(in)
		if err != nil {
			return "", err
		}

		return q, nil
	}

	return in, nil
}

//WaitForCondition is a testify Condition for timeout based testing
func WaitForCondition(d time.Duration, testFn func() bool) bool {
	if d < time.Second {
		panic("WaitForCondition: test duration to small")
	}

	timeout := time.Tick(d)
	test := time.Tick(500 * time.Millisecond)
	check := make(chan struct{}, 1)
	defer close(check)
	done := make(chan struct{}, 1)
	defer close(done)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-test:
				if testFn() {
					check <- struct{}{}
					return
				}
			}
		}
	}()

	for {
		select {
		case <-check:
			return true
		case <-timeout:
			done <- struct{}{}
			return false
		}
	}
}
