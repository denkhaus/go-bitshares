package types

import (
	"encoding/json"

	"github.com/juju/errors"
	"github.com/pquerna/ffjson/ffjson"
)

type AssetFeeds []AssetFeed

type AssetFeed struct {
	ProviderID AccountID
	DateTime   Time
	FeedInfo   PriceFeed
}

func (p AssetFeed) MarshalJSON() ([]byte, error) {
	data := make([]interface{}, 2)
	data[0] = p.DateTime
	data[1] = p.FeedInfo

	ret := make([]interface{}, 2)
	ret[0] = p.ProviderID
	ret[1] = data

	return ffjson.Marshal(ret)
}

func (p *AssetFeed) UnmarshalJSON(data []byte) error {

	raw := make([]json.RawMessage, 2)
	if err := ffjson.Unmarshal(data, &raw); err != nil {
		return errors.Annotate(err, "unmarshal AssetFeed [unmarshal]")
	}

	if err := ffjson.Unmarshal(raw[0], &p.ProviderID); err != nil {
		return errors.Annotate(err, "unmarshal AssetFeed [provider id]")
	}

	feedData := make([]json.RawMessage, 2)
	if err := ffjson.Unmarshal(raw[1], &feedData); err != nil {
		return errors.Annotate(err, "unmarshal AssetFeed [feed data]")
	}

	if err := ffjson.Unmarshal(feedData[0], &p.DateTime); err != nil {
		return errors.Annotate(err, "unmarshal AssetFeed [feed time]")
	}

	if err := ffjson.Unmarshal(feedData[1], &p.FeedInfo); err != nil {
		return errors.Annotate(err, "unmarshal AssetFeed [feed info]")
	}

	return nil
}
