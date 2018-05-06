package objects

import (
	"strconv"
	"strings"

	"github.com/juju/errors"
)

type Votes []Vote

type Vote struct {
	typ      uint8
	instance uint32
}

func (p *Vote) UnmarshalJSON(data []byte) error {
	str := string(data)

	if len(str) > 0 && str != "null" {
		q, err := strconv.Unquote(str)
		if err != nil {
			return errors.Annotate(err, "unquote Vote")
		}

		tk := strings.Split(q, ":")
		if len(tk) != 2 {
			return errors.Errorf("unable to unmarshal Vote from %s", str)
		}

		t, err := strconv.Atoi(tk[0])
		if err != nil {
			return errors.Annotate(err, "Atoi Vote [type]")
		}
		p.typ = uint8(t)

		in, err := strconv.Atoi(tk[1])
		if err != nil {
			return errors.Annotate(err, "Atoi Vote [instance]")
		}
		p.instance = uint32(in)

		return nil
	}

	return errors.Errorf("unable to unmarshal Vote from %s", str)
}
