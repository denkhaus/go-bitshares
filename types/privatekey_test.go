package types

import (
	"testing"

	"github.com/denkhaus/bitshares/config"
	"github.com/stretchr/testify/assert"
)

var privKeys = [][]string{
	[]string{"5Hx8KiHLnc3pDLkwe2jujkTTJev72n3Qx7xtyaRNBsJDuejzh9u", "BTS5zzvbDtkbUVU1gFFsKqCE55U7JbjTp6mTh1usFv7KGgXL7HDQk"},
	[]string{"5JyuWmopuyxFyvM9xm8fxXyujzfVnsg9cvE6z3pcib5NW1Av4rP", "BTS5yXqEBShUgcVm7Mve8Fg4RzQ2ftPpmo77aMbz884eX9aeGVvwD"},
	[]string{"5KRZv3ZmkcE71K9KwEKG6pV6pyufkMQgCJrCu8xKLf2y7R7J8YK", "BTS5Z3vsgH6xj6HLXcsU38yo4TyoZs9AUzpfbaXbuxsAYPbutWvEP"},
	[]string{"5JTge2oTwFqfNPhUrrm6upheByG2VXvaXBAqWdDUvK2CsygMG3Z", "BTS5shffTjVoT4J8Zrj3f2mQJw4UVKrnbx5FWYhVgov45EpBf2NYi"},
	[]string{"5JqmjeakPoTz3ComQ7Jgg11jHxywfkJHZPhMJoBomZLrZSfRAvr", "BTS56fy8qpkLzNoguGMPgPNkkznxnx2woEg1qPq7E6gF2SeGSRyK5"},
}

func TestPrivateKey_SharedSecret(t *testing.T) {
	config.SetCurrentConfig(config.ChainIDBTS)

	for _, k := range privKeys {
		wif := k[0]
		pub := k[1]

		key, err := NewPrivateKeyFromWif(wif)
		if err != nil {
			assert.FailNow(t, err.Error(), "NewPrivateKeyFromWif")
		}

		assert.Equal(t, pub, key.PublicKey().String())
	}
}
