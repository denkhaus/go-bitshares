package objects

import (
	"strconv"

	"github.com/juju/errors"
)

type UInt64 uint64

func (p *UInt64) UnmarshalJSON(s []byte) error {
	str := string(s)

	if len(str) > 0 && str != "null" {
		q, err := strconv.Unquote(str)
		if err != nil {
			return err
		}

		*(*uint64)(p), err = strconv.ParseUint(q, 10, 64)
		if err != nil {
			return errors.Annotate(err, "parse uint64")
		}

		return nil
	}

	return errors.Errorf("unable to unmarshal UInt64 from %s", str)
}
