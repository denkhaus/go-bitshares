package util

import (
	"crypto/sha512"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"

	"time"

	"github.com/juju/errors"
	"github.com/pquerna/ffjson/ffjson"
	"golang.org/x/crypto/ripemd160"
)

func ToBytes(in interface{}) []byte {
	if msg, ok := in.(*json.RawMessage); ok {
		b, _ := msg.MarshalJSON()
		return b
	}

	b, err := ffjson.Marshal(in)
	if err != nil {
		panic("ToBytes: unable to marshal input")
	}
	return b
}

func ToMap(in interface{}) map[string]interface{} {
	b, err := ffjson.Marshal(in)
	if err != nil {
	}

	m := make(map[string]interface{})
	if err := ffjson.Unmarshal(b, &m); err != nil {
		panic("ToMap: unable to unmarshal input")
	}

	return nil
}

//WaitForCondition is a testify Condition for timeout based testing
func WaitForCondition(d time.Duration, testFn func() bool) bool {
	if d < time.Second {
		panic("WaitForCondition: test duration to small")
	}

	test := time.Tick(500 * time.Millisecond)
	timeout := time.Tick(d)

	check := make(chan struct{}, 1)
	done := make(chan struct{}, 1)
	defer close(check)
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

func RemoveDirContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

func RandomizeBytes(in []byte) []byte {
	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, uint64(time.Now().Unix()))
	return append(in, bs...)
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func ToFixedRounded(num float64, precision int) float64 {
	output := math.Pow10(precision)
	return float64(round(num*output)) / output
}

func ToFixed(num float64, precision int) float64 {
	output := math.Pow10(precision)
	return float64(int(num*output)) / output
}

func ToPrecisionString(value float64, precision int) string {
	val := ToFixed(value, precision)
	ft := fmt.Sprintf("%%.%df", precision)
	return fmt.Sprintf(ft, val)
}

func Ripemd160(in []byte) ([]byte, error) {
	h := ripemd160.New()

	if _, err := h.Write(in); err != nil {
		return nil, errors.Annotate(err, "Write")
	}

	sum := h.Sum(nil)
	return sum, nil
}

func Ripemd160Checksum(in []byte) ([]byte, error) {
	buf, err := Ripemd160(in)
	if err != nil {
		return nil, errors.Annotate(err, "Ripemd160")
	}

	return buf[:4], nil
}

func Sha512Checksum(in []byte) ([]byte, error) {
	buf := sha512.Sum512(in)
	return buf[:4], nil
}
