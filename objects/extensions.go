package objects

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/util"
)

type Extensions []Extension

//implements TypeMarshaller interface
func (p Extensions) Marshal(enc *util.TypeEncoder) error {

	if err := enc.Encode([]byte{0, 0, 0}); err != nil {
		return err
	}

	/*
		if err := enc.EncodeUVarint(uint64(len(p))); err != nil {
			return errors.Annotate(err, "encode Extensions length")
		}
					if err := enc.Encode([]byte{0}); err != nil {
						return err
					}
			for _, ex := range p {
					if err := enc.Encode(ex); err != nil {
						return errors.Annotate(err, "encode Extension")
					}
				}
					/*
	*/
	/* 	if err := enc.Encode([]byte{0}); err != nil {
	   		return err
	   	}
	*/
	return nil
}

type Extension []interface{}

//implements TypeMarshaller interface
func (p Extension) Marshal(enc *util.TypeEncoder) error {
	if err := enc.EncodeUVarint(uint64(len(p))); err != nil {
		return err
	}

	return nil
}
