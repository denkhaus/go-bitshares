package api

import (
	"github.com/denkhaus/bitshares/objects"
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
	"github.com/pquerna/ffjson/ffjson"
)

const (
	AssetsListAll      = -1
	AssetsMaxBatchSize = 100
)

//ListAssets retrieves assets
//lowerBoundSymbol: Lower bound of symbol names to retrieve
//limit: Maximum number of assets to fetch, if the constant AssetsListAll
//is passed, all existing assets will be retrieved.
func (p *BitsharesApi) ListAssets(lowerBoundSymbol string, limit int) ([]objects.Asset, error) {

	lim := limit
	if limit > AssetsMaxBatchSize || limit == AssetsListAll {
		lim = AssetsMaxBatchSize
	}

	resp, err := p.client.CallApi(0, "list_assets", lowerBoundSymbol, lim)
	if err != nil {
		return nil, errors.Annotate(err, "list_assets")
	}
	//spew.Dump(resp)
	data := resp.([]interface{})
	ret := make([]objects.Asset, len(data))

	for idx, a := range data {
		if err := ffjson.Unmarshal(util.ToBytes(a), &ret[idx]); err != nil {
			return nil, errors.Annotate(err, "unmarshal Asset")
		}
	}

	return ret, nil
}
