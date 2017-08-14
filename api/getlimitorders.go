package api

import (
	"github.com/denkhaus/bitshares/objects"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
	"github.com/pquerna/ffjson/ffjson"
)

const (
	GetLimitOrdersLimit = 100
)

//GetLimitOrders returns a slice of LimitOrder objects.
func (p *BitsharesApi) GetLimitOrders(base, quote objects.GrapheneObject, limit int) ([]objects.LimitOrder, error) {
	if limit > GetLimitOrdersLimit {
		limit = GetLimitOrdersLimit
	}

	resp, err := p.client.CallApi(0, "get_limit_orders", base.Id(), quote.Id(), limit)
	if err != nil {
		return nil, errors.Annotate(err, "get_limit_orders")
	}

	//util.Dump("limitorders in", resp)

	data := resp.([]interface{})
	ret := make([]objects.LimitOrder, len(data))

	for idx, a := range data {
		if err := ffjson.Unmarshal(util.ToBytes(a), &ret[idx]); err != nil {
			return nil, errors.Annotate(err, "unmarshal LimitOrder")
		}
	}

	return ret, nil
}
