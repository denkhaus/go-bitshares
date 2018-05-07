package util

import (
	"encoding/binary"
	"io"
	"strings"

	"github.com/juju/errors"
)

type TypeMarshaller interface {
	Marshal(enc *TypeEncoder) error
}
type TypeEncoder struct {
	w io.Writer
}

func NewTypeEncoder(w io.Writer) *TypeEncoder {
	return &TypeEncoder{w}
}

func (p *TypeEncoder) EncodeVarint(i int64) error {
	if i >= 0 {
		return p.EncodeUVarint(uint64(i))
	}

	b := make([]byte, binary.MaxVarintLen64)
	n := binary.PutVarint(b, i)
	return p.writeBytes(b[:n])
}

func (p *TypeEncoder) EncodeUVarint(i uint64) error {
	b := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(b, i)
	return p.writeBytes(b[:n])
}

func (p *TypeEncoder) EncodeNumber(v interface{}) error {
	if err := binary.Write(p.w, binary.LittleEndian, v); err != nil {
		return errors.Annotatef(err, "TypeEncoder: failed to write number: %v", v)
	}
	return nil
}

func (p *TypeEncoder) Encode(v interface{}) error {
	if m, ok := v.(TypeMarshaller); ok {
		return m.Marshal(p)
	}

	switch v := v.(type) {
	case int:
		return p.EncodeNumber(v)
	case int8:
		return p.EncodeNumber(v)
	case int16:
		return p.EncodeNumber(v)
	case int32:
		return p.EncodeNumber(v)
	case int64:
		return p.EncodeNumber(v)
	case uint:
		return p.EncodeNumber(v)
	case uint8:
		return p.EncodeNumber(v)
	case uint16:
		return p.EncodeNumber(v)
	case uint32:
		return p.EncodeNumber(v)
	case uint64:
		return p.EncodeNumber(v)
	case float32:
		return p.EncodeNumber(v)
	case float64:
		return p.EncodeNumber(v)
	case string:
		return p.EncodeString(v)
	case []byte:
		return p.writeBytes(v)
	case bool:
		if v {
			return p.EncodeNumber(uint8(1))
		} else {
			return p.EncodeNumber(uint8(0))
		}

	default:
		return errors.Errorf("TypeEncoder: unsupported type encountered")
	}
}

func (p *TypeEncoder) EncodeString(v string) error {
	if err := p.EncodeUVarint(uint64(len(v))); err != nil {
		return errors.Annotatef(err, "TypeEncoder: failed to write string: %v", v)
	}

	return p.writeString(v)
}

func (p *TypeEncoder) writeBytes(bs []byte) error {
	if _, err := p.w.Write(bs); err != nil {
		return errors.Annotatef(err, "TypeEncoder: failed to write bytes: %v", bs)
	}
	return nil
}

func (p *TypeEncoder) writeString(s string) error {
	if _, err := io.Copy(p.w, strings.NewReader(s)); err != nil {
		return errors.Annotatef(err, "TypeEncoder: failed to write string: %v", s)
	}
	return nil
}
