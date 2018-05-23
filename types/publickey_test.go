package types

import (
	"testing"

	"github.com/juju/errors"
	"github.com/stretchr/testify/assert"
)

var keys = []string{
	"BTS6K35Bajw29N4fjP4XADHtJ7bEj2xHJ8CoY2P2s1igXTB5oMBhR",
	"BTS4txNeAoSWcDX7oWceKppMb956z5oRx6mQyCJXCUB7aUh1EJp5y",
	"BTS6iUXJDmAPNbHWHtDDcmPTQ6F3nMBqi6pUHdhSkzWNd6grob2JP",
	"BTS5KCRzL27VLBvhPJ1DaXViuUPxyEXjDvVtWaifUkouNr2MkMGSH",
	"BTS6ThjMq97v6dLQUAmdsZfWG9ENq8nghVUhmLMQi52MDqXvtRGNc",
}

func TestNewPublicKey(t *testing.T) {
	for _, k := range keys {
		key, err := NewPublicKey(k)
		if err != nil {
			assert.FailNow(t, errors.Annotate(err, "NewPublicKey").Error())
		}

		assert.Equal(t, key.String(), k)
	}
}
