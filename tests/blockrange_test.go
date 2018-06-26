package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/bradhe/stopwatch"
	"github.com/stretchr/testify/assert"
	// register operations

	"github.com/denkhaus/bitshares/logging"
	_ "github.com/denkhaus/bitshares/operations"
)

func TestBlockRange(t *testing.T) {

	api := NewTestAPI(t, WsFullApiUrl, RpcFullApiUrl)
	defer api.Close()

	block := uint64(26878913)

	for {
		bl, err := api.GetBlock(block)
		if err != nil {
			assert.FailNow(t, err.Error(), "GetBlock")
		}

		nTrx := uint64(len(bl.Transactions))
		fmt.Printf("block %d: binary compare %d transactions\n", block, nTrx)
		watch := stopwatch.Start()

		for _, tx := range bl.Transactions {
			time.Sleep(300 * time.Millisecond) // to avoid EOF from client
			ref, test, err := CompareTransactions(api, &tx, false)
			if err != nil {
				logging.DDump("trx", tx)
				assert.FailNow(t, err.Error(), "CompareTransactions")
				return
			}

			if !assert.Equal(t, ref, test) {
				logging.DDump("trx", tx)
				return
			}
		}

		watch.Stop()
		fmt.Printf("ms/trx:%v\n", time.Duration(uint64(watch.Milliseconds())/nTrx))
		block++
	}

	//util.Dump("get_block >", res)
}
