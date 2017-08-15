package api

import (
	"github.com/denkhaus/bitshares/objects"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
	"github.com/pquerna/ffjson/ffjson"
)

//GetObjects returns a list of Graphene Objects by ID.
func (p *bitsharesAPI) GetObjects(ids ...objects.GrapheneObject) ([]interface{}, error) {

	params := objects.GrapheneObjects(ids).ToObjectIDs()
	resp, err := p.client.CallApi(0, "get_objects", params)
	if err != nil {
		return nil, errors.Annotate(err, "get_objects")
	}

	//	util.Dump("objects in", resp)

	data := resp.([]interface{})
	ret := make([]interface{}, len(data))
	id := objects.GrapheneID{}

	for idx, obj := range data {
		if obj == nil {
			continue
		}

		if err := id.FromRawData(obj); err != nil {
			return nil, errors.Annotate(err, "from raw data")
		}

		b := util.ToBytes(obj)

		switch id.Type() {
		case objects.ObjectTypeAsset:
			ass := objects.Asset{}
			if err := ffjson.Unmarshal(b, &ass); err != nil {
				return nil, errors.Annotate(err, "unmarshal Asset")
			}
			ret[idx] = ass

		case objects.ObjectTypeAccount:
			acc := objects.Account{}
			if err := ffjson.Unmarshal(b, &acc); err != nil {
				return nil, errors.Annotate(err, "unmarshal Account")
			}
			ret[idx] = acc

		case objects.ObjectTypeAssetBitAssetData:
			bit := objects.BitAssetData{}
			if err := ffjson.Unmarshal(b, &bit); err != nil {
				return nil, errors.Annotate(err, "unmarshal BitAssetData")
			}
			ret[idx] = bit

		case objects.ObjectTypeLimitOrder:
			lim := objects.LimitOrder{}
			if err := ffjson.Unmarshal(b, &lim); err != nil {
				return nil, errors.Annotate(err, "unmarshal LimitOrder")
			}
			ret[idx] = lim

		case objects.ObjectTypeCallOrder:
			cal := objects.CallOrder{}
			if err := ffjson.Unmarshal(b, &cal); err != nil {
				return nil, errors.Annotate(err, "unmarshal CallOrder")
			}
			ret[idx] = cal

		case objects.ObjectTypeForceSettlement:
			set := objects.SettleOrder{}
			if err := ffjson.Unmarshal(b, &set); err != nil {
				return nil, errors.Annotate(err, "unmarshal SettleOrder")
			}
			ret[idx] = set

		default:
			return nil, errors.Errorf("unable to parse GrapheneObject with ID %s", id)
		}
	}

	return ret, nil
}
