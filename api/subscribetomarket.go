package api

import (
	"github.com/denkhaus/bitshares/objects"
	"github.com/juju/errors"
)

func (p *bitsharesAPI) SubscribeToMarket(notifyID int, base objects.GrapheneObject, quote objects.GrapheneObject) error {

	// returns nil if successfull
	_, err := p.client.CallAPI(p.databaseApiID, "subscribe_to_market", notifyID, base.Id(), quote.Id())
	if err != nil {
		return errors.Annotate(err, "subscribe_to_market")
	}

	return nil
}

func (p *bitsharesAPI) UnsubscribeFromMarket(base objects.GrapheneObject, quote objects.GrapheneObject) error {

	// returns nil if successfull
	_, err := p.client.CallAPI(p.databaseApiID, "unsubscribe_from_market", base.Id(), quote.Id())
	if err != nil {
		return errors.Annotate(err, "unsubscribe_from_market")
	}

	return nil
}
