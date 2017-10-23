package objects

import "github.com/denkhaus/bitshares/util"

type Memo struct {
}

//implements Operation interface
func (p Memo) Marshal(enc *util.TypeEncoder) error {
	return nil
}
