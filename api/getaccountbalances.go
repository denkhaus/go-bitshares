package api

import (
	"github.com/denkhaus/bitshares/objects"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
	"github.com/pquerna/ffjson/ffjson"
)

//GetAccountBalances retrieves AssetAmount objects by given AccountID
func (p *BitsharesApi) GetAccountBalances(account objects.GrapheneObject, assets ...objects.GrapheneObject) ([]objects.AssetAmount, error) {

	ids := objects.GrapheneObjects(assets).ToObjectIDs()
	resp, err := p.client.CallApi(0, "get_account_balances", account.Id(), ids)
	if err != nil {
		return nil, errors.Annotate(err, "get_account_balances")
	}

	data := resp.([]interface{})
	ret := make([]objects.AssetAmount, len(data))

	for idx, a := range data {
		if err := ffjson.Unmarshal(util.ToBytes(a), &ret[idx]); err != nil {
			return nil, errors.Annotate(err, "unmarshal AssetAmount")
		}
	}

	return ret, nil
}
