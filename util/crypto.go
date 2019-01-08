package util

import (
	"crypto/sha512"

	"github.com/juju/errors"
	"golang.org/x/crypto/ripemd160"
)

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
