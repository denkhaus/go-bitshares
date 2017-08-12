package api

import (
	"github.com/denkhaus/bitshares/objects"
	"github.com/juju/errors"
	"github.com/pquerna/ffjson/ffjson"
)

//GetAccountByName returns a Account object by username.
func (p *BitsharesApi) GetAccountByName(name string) (*objects.Account, error) {
	resp, err := p.client.CallApi(0, "get_account_by_name", name)
	if err != nil {
		return nil, errors.Annotate(err, "get_account_by_name")
	}

	//spew.Dump(resp)
	ret := objects.Account{}
	if err := ffjson.Unmarshal(toBytes(resp), &ret); err != nil {
		return nil, errors.Annotate(err, "unmarshal Account")
	}

	return &ret, nil
}
