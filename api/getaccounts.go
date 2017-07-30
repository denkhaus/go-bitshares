package api

import (
	"github.com/denkhaus/bitshares/objects"
	"github.com/juju/errors"
	"github.com/mailru/easyjson"
)

//Get a list of accounts by ID.
func (p *BitsharesApi) GetAccounts(accountIDs ...objects.GrapheneObject) ([]objects.Account, error) {
	params := []interface{}{}
	for _, ai := range accountIDs {
		params = append(params, ai.Id())
	}

	resp, err := p.client.CallApi(0, "get_accounts", params)
	if err != nil {
		return nil, errors.Annotate(err, "get_accounts")
	}

	data := resp.([]interface{})
	ret := make([]objects.Account, len(data))

	for idx, acct := range data {
		if err := easyjson.Unmarshal(toBytes(acct), &ret[idx]); err != nil {
			return nil, errors.Annotate(err, "unmarshal Account")
		}
	}

	return ret, nil
}
