package api

import (
	"github.com/juju/errors"
)

//GetChainID returns the ID of the chain we are connected to.
func (p *bitsharesAPI) GetChainID() (string, error) {
	/* if err := p.ensureInitialized(); err != nil {
		return "", errors.Annotate(err, "ensure initialized")
	} */

	resp, err := p.client.CallApi(p.databaseApiID, "get_chain_id", EmptyParams)
	if err != nil {
		return "", errors.Annotate(err, "get_chain_id")
	}

	//util.Dump("get_chain_id in", resp)
	return resp.(string), nil
}
