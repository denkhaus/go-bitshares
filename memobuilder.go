package bitshares

import (
	"math/rand"

	"github.com/denkhaus/bitshares/crypto"
	"github.com/denkhaus/bitshares/types"
	"github.com/juju/errors"
)

// MemoBuilder is a memo factory
type MemoBuilder struct {
	from types.GrapheneObject
	to   types.GrapheneObject
	api  WebsocketAPI
	memo string
}

//NewMemoBuilder creates a new MemoBuilder
func (p *websocketAPI) NewMemoBuilder(from, to types.GrapheneObject, memo string) *MemoBuilder {
	builder := MemoBuilder{
		from: from,
		to:   to,
		memo: memo,
		api:  p,
	}

	return &builder
}

// Encrypt encrypts the memo message with the corresponding private key found in keyBag
func (p *MemoBuilder) Encrypt(keyBag *crypto.KeyBag) (*types.Memo, error) {
	accts, err := p.api.GetAccounts(p.from, p.to)
	if err != nil {
		return nil, errors.Annotate(err, "GetAccounts")
	}

	from := accts.Lookup(p.from)
	if from == nil {
		return nil, errors.Errorf("can't retrieve account infos for %s", p.from)
	}

	to := accts.Lookup(p.to)
	if to == nil {
		return nil, errors.Errorf("can't retrieve account infos for %s", p.to)
	}

	memo := types.Memo{
		From:  from.Options.MemoKey,
		To:    to.Options.MemoKey,
		Nonce: types.UInt64(rand.Int63()),
	}

	if err := keyBag.EncryptMemo(&memo, p.memo); err != nil {
		return nil, errors.Annotate(err, "EncryptMemo")
	}

	return &memo, nil
}
