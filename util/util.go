package util

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

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
