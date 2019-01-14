package main

import (
	"fmt"
	"log"

	"github.com/denkhaus/bitshares"
	"github.com/denkhaus/bitshares/config"
	"github.com/denkhaus/bitshares/crypto"
	"github.com/denkhaus/bitshares/tests"
	"github.com/denkhaus/bitshares/types"
	"github.com/juju/errors"
)

var (
	from    = tests.TestAccount1ID
	to      = tests.TestAccount2ID
	test    = tests.AssetTEST
	memoWIF = tests.TestAccount1PrivKeyActive
	memo    = "my super secret memo message"

	keyBag *crypto.KeyBag
)

const (
	wsTestApiUrl = "wss://node.testnet.bitshares.eu/ws"
)

func init() {
	// init is called before the API is initialized,
	// hence must define current chain config explicitly.
	config.SetCurrent(config.ChainIDTest)
	keyBag = crypto.NewKeyBag()
	if err := keyBag.Add(memoWIF); err != nil {
		log.Fatal(errors.Annotate(err, "Add [wif]"))
	}
}

func main() {
	api := bitshares.NewWebsocketAPI(wsTestApiUrl)
	if err := api.Connect(); err != nil {
		log.Fatal(errors.Annotate(err, "OnConnect"))
	}

	api.OnError(func(err error) {
		log.Fatal(errors.Annotate(err, "OnError"))
	})

	amt := types.AssetAmount{
		Asset:  types.AssetIDFromObject(test),
		Amount: 10000,
	}

	err := api.Transfer(keyBag, from, to, test, amt, memo)
	if err != nil {
		log.Fatal(errors.Annotate(err, "Transfer"))
	}

	fmt.Println("transfer successful")
}
