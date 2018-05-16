package util

import (
	"bytes"
	"encoding/binary"
	"io"
	"strings"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil/base58"
	"github.com/juju/errors"
	"golang.org/x/crypto/ripemd160"
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
	if v == nil {
		return nil
	}

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

func (p *TypeEncoder) EncodeArrString(v []string) error {
	if err := p.EncodeUVarint(uint64(len(v))); err != nil {
		return errors.Annotatef(err, "TypeEncoder: failed to write string: %v", v)
	}

	for _, val := range v {
		if err := p.EncodeUVarint(uint64(len(val))); err != nil {
			return errors.Annotatef(err, "TypeEncoder: failed to write string: %v", val)
		}
		if _, err := io.Copy(p.w, strings.NewReader(val)); err != nil {
			return errors.Annotatef(err, "TypeEncoder: failed to write string: %v", val)
		}
	}

	return nil
}

func (p *TypeEncoder) EncodeString(v string) error {
	if err := p.EncodeUVarint(uint64(len(v))); err != nil {
		return errors.Annotatef(err, "TypeEncoder: failed to write string: %v", v)
	}

	return p.writeString(v)
}

func (p *TypeEncoder) EncodePubKey(s string) error {
	pkn1 := strings.Join(strings.Split(s, "")[3:], "")
	b58 := base58.Decode(pkn1)
	chs := b58[len(b58)-4:]
	pkn2 := b58[:len(b58)-4]
	chHash := ripemd160.New()
	chHash.Write(pkn2)
	nchs := chHash.Sum(nil)[:4]

	if bytes.Equal(chs, nchs) {
		pkn3, _ := btcec.ParsePubKey(pkn2, btcec.S256())
		if _, err := p.w.Write(pkn3.SerializeCompressed()); err != nil {
			return errors.Annotatef(err, "TypeEncoder: failed to write bytes: %v", pkn3.SerializeCompressed())
		}

		return nil
	}

	return errors.New("Public key is incorrect")
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
