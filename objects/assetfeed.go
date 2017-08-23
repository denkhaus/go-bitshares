package objects

import (
	json "encoding/json"

	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
)

type AssetFeed struct {
	ProviderID GrapheneID
	DateTime   Time
	FeedInfo   AssetFeedInfo
}

func (p *AssetFeed) UnmarshalJSON(data []byte) error {
	if data == nil || len(data) == 0 {
		return nil
	}

	var res interface{}
	if err := json.Unmarshal(data, &res); err != nil {
		return errors.Annotate(err, "unmarshal AssetFeed [unmarshal]")
	}

	d := res.([]interface{})
	if err := p.ProviderID.UnmarshalJSON(util.ToBytes(d[0])); err != nil {
		return errors.Annotate(err, "unmarshal AssetFeed [provider id]")
	}

	feedData := d[1].([]interface{})
	if err := p.DateTime.UnmarshalJSON(util.ToBytes(feedData[0])); err != nil {
		util.Dump("time1", feedData[0])
		return errors.Annotate(err, "unmarshal AssetFeed [feed time]")
	}

	//comment this out while generate ffjson logic -> circular dependencies!!
	if err := p.FeedInfo.UnmarshalJSON(util.ToBytes(feedData[1])); err != nil {
		return errors.Annotate(err, "unmarshal AssetFeed [feed info]")
	}

	return nil
}
