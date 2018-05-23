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

	d, ok := res.([]interface{})
	if !ok {
		return ErrInvalidInputType
	}

	if len(d) != 2 {
		return ErrInvalidInputLength
	}

	if err := p.ProviderID.UnmarshalJSON(util.ToBytes(d[0])); err != nil {
		return errors.Annotate(err, "unmarshal AssetFeed [provider id]")
	}

	feedData, ok := d[1].([]interface{})
	if !ok {
		return ErrInvalidInputType
	}

	if err := p.DateTime.UnmarshalJSON(util.ToBytes(feedData[0])); err != nil {
		return errors.Annotate(err, "unmarshal AssetFeed [feed time]")
	}

	//this gives us an error if we generate ffjson logic for the first time
	//for now comment this out to generate and in again
	if err := p.FeedInfo.UnmarshalJSON(util.ToBytes(feedData[1])); err != nil {
		return errors.Annotate(err, "unmarshal AssetFeed [feed info]")
	}

	return nil
}
