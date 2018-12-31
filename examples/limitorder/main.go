package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/denkhaus/bitshares"
	"github.com/denkhaus/bitshares/config"
	"github.com/denkhaus/bitshares/crypto"
	"github.com/denkhaus/bitshares/operations"
	"github.com/denkhaus/bitshares/types"
	"github.com/juju/errors"
)

var (
	bts    = types.NewGrapheneID("1.3.0")
	cny    = types.NewGrapheneID("1.3.113")
	keyBag *crypto.KeyBag
	seller *types.GrapheneID
)

const (
	wsFullApiUrl = "wss://bitshares.openledger.info/ws"
)

func init() {
	// init is called before the API is initialized,
	// hence must define current chain config explicitly.
	config.SetCurrentConfig(config.ChainIDBTS)
	seller = types.NewGrapheneID(
		os.Getenv("BTS_TEST_ACCOUNT"),
	)
	keyBag = crypto.NewKeyBag()
	if err := keyBag.Add(os.Getenv("BTS_TEST_WIF")); err != nil {
		log.Fatal(errors.Annotate(err, "Add [wif]"))
	}
}

func main() {
	api := bitshares.NewWebsocketAPI(wsFullApiUrl)
	if err := api.Connect(); err != nil {
		log.Fatal(errors.Annotate(err, "OnConnect"))
	}

	api.OnError(func(err error) {
		log.Fatal(errors.Annotate(err, "OnError"))
	})

	op := operations.LimitOrderCreateOperation{
		FillOrKill: false,
		Seller:     *seller,
		Extensions: types.Extensions{},
		AmountToSell: types.AssetAmount{
			Amount: 100,
			Asset:  *bts,
		},
		MinToReceive: types.AssetAmount{
			Amount: 1000000,
			Asset:  *cny,
		},
	}

	op.Expiration.Set(24 * time.Hour)

	tx, err := api.BuildSignedTransaction(keyBag, bts, &op)
	if err != nil {
		log.Fatal(errors.Annotate(err, "BuildSignedTransaction"))
	}

	if err := api.BroadcastTransaction(tx); err != nil {
		log.Fatal(errors.Annotate(err, "BroadcastTransaction"))
	}

	fmt.Println("operation successfull broadcasted")
}
