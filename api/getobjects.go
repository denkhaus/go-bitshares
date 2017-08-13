package api

import (
	"github.com/denkhaus/bitshares/objects"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
	"github.com/pquerna/ffjson/ffjson"
)

//GetObjects returns a list of Graphene Objects by ID.
func (p *BitsharesApi) GetObjects(ids ...*objects.GrapheneID) ([]objects.GrapheneObject, error) {

	params := []interface{}{}
	for _, id := range ids {
		params = append(params, id.Id())
	}

	resp, err := p.client.CallApi(0, "get_objects", params)
	if err != nil {
		return nil, errors.Annotate(err, "get_objects")
	}

	//	util.Dump("objects in", resp)
	data := resp.([]interface{})
	ret := make([]objects.GrapheneObject, len(data))

	for idx, obj := range data {
		b := util.ToBytes(obj)
		if err := ffjson.Unmarshal(b, &ret[idx]); err != nil {
			return nil, errors.Annotate(err, "unmarshal GrapheneObject")
		}

		switch ret[idx].ID.Type() {
		case objects.ObjectTypeAsset:
			ass := objects.Asset{}
			if err := ffjson.Unmarshal(b, &ass); err != nil {
				return nil, errors.Annotate(err, "unmarshal Asset")
			}
			ret[idx].Data = ass

		case objects.ObjectTypeAccount:
			acc := objects.Account{}
			if err := ffjson.Unmarshal(b, &acc); err != nil {
				return nil, errors.Annotate(err, "unmarshal Account")
			}
			ret[idx].Data = acc

		case objects.ObjectTypeAssetBitAssetData:

			bit := objects.BitAssetData{}
			if err := ffjson.Unmarshal(b, &bit); err != nil {
				return nil, errors.Annotate(err, "unmarshal BitAssetData")
			}
			ret[idx].Data = bit

		default:
			return nil, errors.New("unable to parse GrapheneObject")
		}
	}

	return ret, nil
}
