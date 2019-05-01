package util

import (
	"encoding/binary"
	"io"
	"reflect"

	"github.com/juju/errors"
)

var (
	ErrCannotDecodeNilValue = errors.New("cannot decode nil value")
)

type TypeUnmarshaler interface {
	Unmarshal(enc *TypeDecoder) error
}

type TypeDecoder struct {
	r io.Reader
}

func NewTypeDecoder(r io.Reader) *TypeDecoder {
	return &TypeDecoder{r}
}

func (p *TypeDecoder) DecodeUVarint(v interface{}) error {
	br := ByteReader{p.r}
	val, err := binary.ReadUvarint(br)
	if err != nil {
		return errors.Annotate(err, "ReadUvarint")
	}

	reflect.ValueOf(v).Elem().SetUint(val)
	return nil
}

func (p *TypeDecoder) DecodeNumber(v interface{}) error {
	if err := binary.Read(p.r, binary.LittleEndian, v); err != nil {
		return errors.Annotate(err, "Read")
	}
	return nil
}

func (p *TypeDecoder) Decode(v interface{}) error {
	if v == nil {
		return ErrCannotDecodeNilValue
	}

	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr {
		return errors.New("v must be a pointer to a valid decode target")
	}

	if m, ok := v.(TypeUnmarshaler); ok {
		return m.Unmarshal(p)
	}

	trg := val.Elem().Interface()
	switch trg.(type) {
	case int8:
		return p.DecodeNumber(v)
	case int16:
		return p.DecodeNumber(v)
	case int32:
		return p.DecodeNumber(v)
	case int64:
		return p.DecodeNumber(v)
	case uint:
		return p.DecodeNumber(v)
	case uint8:
		return p.DecodeNumber(v)
	case uint16:
		return p.DecodeNumber(v)
	case uint32:
		return p.DecodeNumber(v)
	case uint64:
		return p.DecodeNumber(v)
	case float32:
		return p.DecodeNumber(v)
	case float64:
		return p.DecodeNumber(v)
	case string:
		return p.DecodeString(v)
	case bool:
		var val uint8
		if err := p.DecodeNumber(&val); err != nil {
			return errors.Annotate(err, "DecodeNumber")
		}
		reflect.ValueOf(v).Elem().SetBool(val == 1)
		return nil

	default:
		return errors.Errorf("TypeDecoder: unsupported type encountered")
	}
}

func (p *TypeDecoder) DecodeString(v interface{}) error {
	var length uint64
	if err := p.DecodeUVarint(&length); err != nil {
		return errors.Annotate(err, "DecodeUVarint")
	}

	buf := make([]byte, length)
	n, err := p.r.Read(buf)
	if err != nil {
		return errors.Annotate(err, "Read")
	}

	buf = buf[:n]

	reflect.ValueOf(v).Elem().SetString(string(buf))
	return nil
}

func (p *TypeDecoder) ReadBytes(v interface{}, len uint64) error {
	buf := make([]byte, len)
	if _, err := p.r.Read(buf); err != nil {
		return errors.Annotate(err, "Read")
	}

	reflect.ValueOf(v).Elem().SetBytes(buf)
	return nil
}

type ByteReader struct {
	io.Reader
}

func (br ByteReader) ReadByte() (byte, error) {
	buf := make([]byte, 1)
	if _, err := br.Read(buf); err != nil {
		return 0, err
	}

	return buf[0], nil
}
