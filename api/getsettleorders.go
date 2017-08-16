package api

import (
	"github.com/denkhaus/bitshares/objects"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
	"github.com/pquerna/ffjson/ffjson"
)

const (
	GetSettleOrdersLimit = 100
)

//GetSettleOrders returns a slice of SettleOrder objects.
func (p *bitsharesAPI) GetSettleOrders(assetID objects.GrapheneObject, limit int) ([]objects.SettleOrder, error) {
	if limit > GetSettleOrdersLimit {
		limit = GetSettleOrdersLimit
	}

	resp, err := p.client.CallAPI(0, "get_settle_orders", assetID.Id(), limit)
	if err != nil {
		return nil, errors.Annotate(err, "get_settle_orders")
	}

	//util.Dump("settleorders in", resp)

	data := resp.([]interface{})
	ret := make([]objects.SettleOrder, len(data))

	for idx, a := range data {
		if err := ffjson.Unmarshal(util.ToBytes(a), &ret[idx]); err != nil {
			return nil, errors.Annotate(err, "unmarshal SettleOrder")
		}
	}

	return ret, nil
}
