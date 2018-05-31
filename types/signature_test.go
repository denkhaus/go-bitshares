package types

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/btcsuite/btcd/btcec"
	"github.com/denkhaus/bitshares/config"
	"github.com/stretchr/testify/assert"
)

var (
	digHex = "236e76bc735645982202207db3338fb79cece62de4970bfa09e1094ed79434b4"
	sigHex = "204f2f225577867483b0fbb675ade7b80ea62b9d9761857b369b0a9796cb2efc227d7828997904f9569423a3119b4b05ee5796e16cfcf2d599eb63c378a52c6e2b"
)

func Test_SignatureVerify(t *testing.T) {
	config.SetCurrentConfig(config.ChainIDTest)

	dig, err := hex.DecodeString(digHex)
	if err != nil {
		assert.FailNow(t, err.Error(), "DecodeString")
	}

	sig, err := hex.DecodeString(sigHex)
	if err != nil {
		assert.FailNow(t, err.Error(), "DecodeString")
	}

	s, err := btcec.ParseSignature(sig[1:], btcec.S256())
	if err != nil {
		assert.FailNow(t, err.Error(), "ParseSignature")
	}

	_ = s

	p, comp, err := btcec.RecoverCompact(btcec.S256(), sig, dig)
	if err != nil {
		assert.FailNow(t, err.Error(), "RecoverCompact")
	}

	fmt.Println("compressed: ", comp)

	pub, err := NewPublicKey(p)
	if err != nil {
		assert.FailNow(t, err.Error(), "NewPublicKey")
	}

	fmt.Println("public key: ", pub.String())

}
