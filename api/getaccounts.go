package api

import (
	"github.com/denkhaus/bitshares/objects"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
	"github.com/pquerna/ffjson/ffjson"
)

//GetAccounts returns a list of accounts by ID.
func (p *BitsharesApi) GetAccounts(accounts ...objects.GrapheneObject) ([]objects.Account, error) {

	ids := objects.GrapheneObjects(accounts).ToObjectIDs()
	resp, err := p.client.CallApi(0, "get_accounts", ids)
	if err != nil {
		return nil, errors.Annotate(err, "get_accounts")
	}

	//spew.Dump(resp)
	data := resp.([]interface{})
	ret := make([]objects.Account, len(data))

	for idx, acct := range data {
		if err := ffjson.Unmarshal(util.ToBytes(acct), &ret[idx]); err != nil {
			return nil, errors.Annotate(err, "unmarshal Account")
		}
	}

	return ret, nil
}
