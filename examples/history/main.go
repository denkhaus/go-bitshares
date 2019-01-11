package main

import (
	"fmt"
	"log"

	"github.com/denkhaus/bitshares"
	"github.com/denkhaus/bitshares/types"
	"github.com/juju/errors"
)

var (
	start   = types.NewOperationHistoryID("1.11.1000000")
	stop    = types.NewOperationHistoryID("1.11.0")
	account = types.NewAccountID("1.2.253")
)

const (
	limit = 100
	//wsFullApiUrl = "wss://bitshares.openledger.info/ws"
	wsFullApiUrl = "wss://node.market.rudex.org"
)

func main() {
	api := bitshares.NewWebsocketAPI(wsFullApiUrl)
	if err := api.Connect(); err != nil {
		log.Fatal(errors.Annotate(err, "OnConnect"))
	}

	defer api.Close()
	api.OnError(func(err error) {
		log.Fatal(errors.Annotate(err, "OnError"))
	})

	nObjectsTotal := 0

	for {
		hist, err := api.GetAccountHistory(account, stop, limit, start)
		if err != nil {
			log.Fatal(errors.Annotate(err, "GetAccountHistory"))
		}

		nObjects := len(hist)
		//don't repeat the last object over and over
		if nObjects == 1 {
			break
		}

		for _, op := range hist {
			if start.Instance() > op.ID.Instance() {
				start = &op.ID
			}

			nObjectsTotal++
			fmt.Printf("object #%d: id = %v , typ = %s\n",
				nObjectsTotal, op.ID, op.Operation.Type,
			)
		}
	}
}
