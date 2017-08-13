package objects

import (
	"strconv"
	"time"

	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

type RFC3339Time time.Time

func (t RFC3339Time) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(time.Time(t).Format(time.RFC3339))), nil
}

func (t *RFC3339Time) UnmarshalJSON(s []byte) error {
	str := string(s)

	if len(str) > 0 && str != "null" {
		q, err := util.SafeUnquote(str)
		if err != nil {
			return errors.Annotate(err, "unquote")
		}

		dt, err := time.Parse("2006-01-02T15:04:05", q)
		if err != nil {
			return errors.Annotate(err, "parse datetime string")
		}

		*(*time.Time)(t) = dt
	}

	return nil
}

func (t RFC3339Time) Unix() int64 {
	return time.Time(t).Unix()
}

func (t RFC3339Time) ToTime() time.Time {
	return time.Time(t)
}

func (t RFC3339Time) String() string {
	return time.Time(t).String()
}
