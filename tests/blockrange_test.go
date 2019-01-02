package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/bradhe/stopwatch"
	"github.com/denkhaus/logging"
	"github.com/stretchr/testify/assert"

	// register operations

	_ "github.com/denkhaus/bitshares/operations"
)

func TestBlockRange(t *testing.T) {
	websocketAPI := NewWebsocketTestAPI(t, WsFullApiUrl)
	defer websocketAPI.Close()

	walletAPI := NewWalletTestAPI(t, RpcFullApiUrl)
	defer walletAPI.Close()

	block := uint64(26880513)

	for {
		bl, err := websocketAPI.GetBlock(block)
		if err != nil {
			assert.FailNow(t, err.Error(), "GetBlock")
		}

		nTrx := int64(len(bl.Transactions))
		fmt.Printf("block %d: binary compare %d transactions\n", block, nTrx)
		watch := stopwatch.Start()

		for _, tx := range bl.Transactions {
			ref, test, err := CompareTransactions(walletAPI, &tx, false)
			if err != nil {
				logging.Dump("trx", tx)
				assert.FailNow(t, err.Error(), "CompareTransactions")
				return
			}

			if !assert.Equal(t, ref, test) {
				logging.DumpJSON("trx", tx)
				return
			}
		}

		watch.Stop()
		fmt.Printf("time/trx:%v\n",
			time.Duration(int64(watch.Milliseconds()*time.Millisecond)/nTrx))
		block++
	}

	//util.Dump("get_block >", res)
}
