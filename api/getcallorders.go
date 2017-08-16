package api

import (
	"github.com/denkhaus/bitshares/objects"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
	"github.com/pquerna/ffjson/ffjson"
)

const (
	GetCallOrdersLimit = 100
)

//GetCallOrders returns a slice of CallOrder objects.
func (p *bitsharesAPI) GetCallOrders(assetID objects.GrapheneObject, limit int) ([]objects.CallOrder, error) {
	if limit > GetCallOrdersLimit {
		limit = GetCallOrdersLimit
	}

	resp, err := p.client.CallAPI(0, "get_call_orders", assetID.Id(), limit)
	if err != nil {
		return nil, errors.Annotate(err, "get_call_orders")
	}

	//util.Dump("callorders in", resp)

	data := resp.([]interface{})
	ret := make([]objects.CallOrder, len(data))

	for idx, a := range data {
		if err := ffjson.Unmarshal(util.ToBytes(a), &ret[idx]); err != nil {
			return nil, errors.Annotate(err, "unmarshal CallOrder")
		}
	}

	return ret, nil
}
