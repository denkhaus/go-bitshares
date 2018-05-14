package types

import (
	"fmt"

	"strconv"
	"strings"

	"github.com/denkhaus/bitshares/util"
	sort "github.com/emirpasic/gods/utils"
	"github.com/juju/errors"
)

type Votes []VoteID

func (p *Votes) UnmarshalJSON(data []byte) error {
	return ErrNotImplemented
}

//TODO: define this
func (p Votes) Marshal(enc *util.TypeEncoder) error {
	if err := enc.EncodeUVarint(uint64(len(p))); err != nil {
		return errors.Annotate(err, "encode length")
	}

	//TODO: remove duplicates
	//copy votes and sort
	votes := make([]interface{}, len(p))
	for idx, id := range p {
		votes[idx] = id
	}

	sort.Sort(votes, voteIDComparator)
	for _, v := range votes {
		if err := enc.Encode(v); err != nil {
			return errors.Annotate(err, "encode VoteID")
		}
	}

	return nil
}

type VoteID struct {
	typ      int
	instance int
}

func (p *VoteID) UnmarshalJSON(data []byte) error {
	str := string(data)

	q, err := util.SafeUnquote(str)
	if err != nil {
		return errors.Annotate(err, "unquote VoteID")
	}

	tk := strings.Split(q, ":")
	if len(tk) != 2 {
		return errors.Errorf("unable to unmarshal Vote from %s", str)
	}

	t, err := strconv.Atoi(tk[0])
	if err != nil {
		return errors.Annotate(err, "Atoi VoteID [type]")
	}
	p.typ = t

	in, err := strconv.Atoi(tk[1])
	if err != nil {
		return errors.Annotate(err, "Atoi VoteID [instance]")
	}
	p.instance = in

	return nil
}

func (p VoteID) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%d:%d"`, p.typ, p.instance)), nil
}

//TODO: define this
func (p VoteID) Marshal(enc *util.TypeEncoder) error {
	bin := (p.typ & 0xff) | (p.instance << 8)
	if err := enc.Encode(uint32(bin)); err != nil {
		return errors.Annotate(err, "encode ID")
	}

	return nil
}

func NewVoteID(id string) *VoteID {
	v := VoteID{}
	if err := v.UnmarshalJSON([]byte(id)); err != nil {
		panic(errors.Annotatef(err, "unmarshal VoteID from %v", id))
	}

	return &v
}

func voteIDComparator(a, b interface{}) int {
	aID := a.(VoteID)
	bID := b.(VoteID)

	switch {
	case aID.instance > bID.instance:
		return 1
	case aID.instance < bID.instance:
		return -1
	default:
		return 0
	}
}
