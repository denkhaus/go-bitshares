package types

//go:generate ffjson $GOFILE

import (
	"encoding/json"

	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

type Extensions json.RawMessage

//TODO refactor and test
func (p Extensions) Marshal(enc *util.TypeEncoder) error {
	if err := enc.EncodeUVarint(uint64(0)); err != nil {
		return errors.Annotate(err, "encode length")
	}

	// if err := enc.EncodeUVarint(uint64(len(p))); err != nil {
	// 	return errors.Annotate(err, "encode length")
	// }

	// for _, ex := range p {
	// 	if err := enc.Encode(ex); err != nil {
	// 		return errors.Annotate(err, "encode Extension")
	// 	}
	// }

	return nil
}
