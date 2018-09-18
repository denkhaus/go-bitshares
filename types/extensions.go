package types

//go:generate ffjson $GOFILE

import (
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

type Extensions []interface{}

func (p Extensions) Marshal(enc *util.TypeEncoder) error {
	if err := enc.EncodeUVarint(uint64(len(p))); err != nil {
		return errors.Annotate(err, "encode length")
	}

	for _, ex := range p {
		if err := enc.Encode(ex); err != nil {
			return errors.Annotate(err, "encode Extension")
		}
	}

	return nil
}

// func (p *Extensions) UnmarshalJSON(s []byte) error {
// 	var val interface{}
// 	if err := ffjson.Unmarshal(s, &val); err != nil {
// 		return errors.Annotate(err, "Unmarshal")
// 	}

// 	fmt.Printf("extension %+v", val)
// 	return nil
// }
