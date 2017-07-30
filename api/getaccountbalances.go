package api

import (
	"github.com/denkhaus/bitshares/objects"
	"github.com/juju/errors"
	"github.com/mailru/easyjson"
)

//GetAccountBalances retrieves AssetAmount objects by given Account
func (p *BitsharesApi) GetAccountBalances(account objects.GrapheneObject, assets ...objects.GrapheneObject) ([]objects.AssetAmount, error) {
	assetIDs := []interface{}{"1.3.0"}
	for _, asset := range assets {
		assetIDs = append(assetIDs, asset.Id())
	}

	resp, err := p.client.CallApi(0, "get_account_balances", account.Id(), assetIDs)
	if err != nil {
		return nil, errors.Annotate(err, "get_account_balances")
	}

	data := resp.([]interface{})
	ret := make([]objects.AssetAmount, len(data))

	for idx, a := range data {
		if err := easyjson.Unmarshal(toBytes(a), &ret[idx]); err != nil {
			return nil, errors.Annotate(err, "unmarshal AssetAmount")
		}
	}

	return ret, nil
}
