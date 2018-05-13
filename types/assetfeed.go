package types

import (
	"github.com/denkhaus/bitshares/util"
	"github.com/juju/errors"
	"github.com/pquerna/ffjson/ffjson"
)

type AssetFeeds []AssetFeed

type AssetFeed struct {
	ProviderID GrapheneID
	DateTime   Time
	FeedInfo   PriceFeed
}

func (p *AssetFeed) UnmarshalJSON(data []byte) error {
	if data == nil || len(data) == 0 {
		return nil
	}

	var res interface{}
	if err := ffjson.Unmarshal(data, &res); err != nil {
		return errors.Annotate(err, "unmarshal AssetFeed [unmarshal]")
	}

	d := res.([]interface{})
	if err := p.ProviderID.UnmarshalJSON(util.ToBytes(d[0])); err != nil {
		return errors.Annotate(err, "unmarshal AssetFeed [provider id]")
	}

	feedData := d[1].([]interface{})
	if err := p.DateTime.UnmarshalJSON(util.ToBytes(feedData[0])); err != nil {
		return errors.Annotate(err, "unmarshal AssetFeed [feed time]")
	}

	//this gives us circular dependencies when we generate ffjson logic for the first time
	//meanwhile comment this out to generate
	// if err := p.FeedInfo.UnmarshalJSON(util.ToBytes(feedData[1])); err != nil {
	// 	return errors.Annotate(err, "unmarshal AssetFeed [feed info]")
	// }

	return nil
}
