package types

import (
	"testing"

	"github.com/denkhaus/bitshares/config"
	"github.com/juju/errors"
	"github.com/stretchr/testify/assert"
)

const (
	BTSNullKey = "BTS1111111111111111111111111111111114T1Anm"
)

var keys = []string{
	"BTS6K35Bajw29N4fjP4XADHtJ7bEj2xHJ8CoY2P2s1igXTB5oMBhR",
	"BTS4txNeAoSWcDX7oWceKppMb956z5oRx6mQyCJXCUB7aUh1EJp5y",
	"BTS6iUXJDmAPNbHWHtDDcmPTQ6F3nMBqi6pUHdhSkzWNd6grob2JP",
	"BTS5KCRzL27VLBvhPJ1DaXViuUPxyEXjDvVtWaifUkouNr2MkMGSH",
	"BTS6ThjMq97v6dLQUAmdsZfWG9ENq8nghVUhmLMQi52MDqXvtRGNc",
	"BTS5zzvbDtkbUVU1gFFsKqCE55U7JbjTp6mTh1usFv7KGgXL7HDQk",
	"BTS5yXqEBShUgcVm7Mve8Fg4RzQ2ftPpmo77aMbz884eX9aeGVvwD",
	"BTS5Z3vsgH6xj6HLXcsU38yo4TyoZs9AUzpfbaXbuxsAYPbutWvEP",
	"BTS5shffTjVoT4J8Zrj3f2mQJw4UVKrnbx5FWYhVgov45EpBf2NYi",
	"BTS56fy8qpkLzNoguGMPgPNkkznxnx2woEg1qPq7E6gF2SeGSRyK5",
}

func TestNewPublicKey(t *testing.T) {
	config.SetCurrentConfig(config.ChainIDBTS)
	for _, k := range keys {
		key, err := NewPublicKeyFromString(k)
		if err != nil {
			assert.FailNow(t, errors.Annotate(err, "NewPublicKeyFromString").Error())
		}

		assert.Equal(t, k, key.String())
	}
}

func TestNullPublicKey(t *testing.T) {
	config.SetCurrentConfig(config.ChainIDBTS)
	key, err := NewPublicKeyFromString(BTSNullKey)
	if err != nil {
		assert.FailNow(t, errors.Annotate(err, "NewPublicKeyFromString").Error())
	}

	assert.NotNil(t, key)
	assert.Equal(t, BTSNullKey, key.String())

}
