package types

import (
	"encoding/hex"
	"testing"

	"github.com/btcsuite/btcd/btcec"
	"github.com/denkhaus/bitshares/config"
	"github.com/stretchr/testify/assert"
)

var (
	digSig = [][]string{
		{
			"aa2036ea32d635ef8094883a75a7aa7e9b7f034e0087c43c08872ccd344a486f",
			"1f49954c9f1df0a8d9f73833173eb40625036b4c718723772c4e3bfb687442637027b9f4aef10c282251fc2a3ad5f77deed5e369f2ac75903e3670818609a3f7ad",
			"TEST5Z3vsgH6xj6HLXcsU38yo4TyoZs9AUzpfbaXbuxsAYPbutWvEP",
		},
	}
)

func Test_SignatureVerify(t *testing.T) {
	config.SetCurrentConfig(config.ChainIDTest)

	for _, ds := range digSig {
		digHex := ds[0]
		sigHex := ds[1]
		pub := ds[2]

		dig, err := hex.DecodeString(digHex)
		if err != nil {
			assert.FailNow(t, err.Error(), "DecodeString")
		}

		sig, err := hex.DecodeString(sigHex)
		if err != nil {
			assert.FailNow(t, err.Error(), "DecodeString")
		}

		p, _, err := btcec.RecoverCompact(btcec.S256(), sig, dig)
		if err != nil {
			assert.FailNow(t, err.Error(), "RecoverCompact")
		}

		pubKey, err := NewPublicKey(p)
		if err != nil {
			assert.FailNow(t, err.Error(), "NewPublicKey")
		}

		assert.Equal(t, pub, pubKey.String())
	}
}
