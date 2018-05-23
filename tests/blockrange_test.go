package tests

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	// register operations
	"github.com/denkhaus/bitshares/api"
	_ "github.com/denkhaus/bitshares/operations"
	"github.com/denkhaus/bitshares/util"
)

func TestBlockRange(t *testing.T) {
	api := api.New(WsFullApiUrl, RpcApiUrl)
	if err := api.Connect(); err != nil {
		assert.FailNow(t, err.Error(), "Connect")
	}

	api.OnError(func(err error) {
		assert.FailNow(t, err.Error(), "OnError")
	})

	block := uint64(26878298)

	for {
		fmt.Println("get block: ", block)
		bl, err := api.GetBlock(block)
		if err != nil {
			assert.FailNow(t, err.Error(), "GetBlock")
		}

		for _, tx := range bl.Transactions {
			ref, test, err := CompareTransactions(api, &tx, false)
			if err != nil {
				util.Dump("trx", tx)
				assert.FailNow(t, err.Error(), "compareTry")
			}

			assert.Equal(t, ref, test)
		}

		block++
	}

	//util.Dump("get_block >", res)
}
