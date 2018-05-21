package api

import (
	"fmt"
	"testing"

	"github.com/denkhaus/bitshares/tests"
	"github.com/stretchr/testify/assert"
	// register operations
	_ "github.com/denkhaus/bitshares/operations"
)

func TestBlockRange(t *testing.T) {
	api := New(tests.WsFullApiUrl, tests.RpcApiUrl)
	if err := api.Connect(); err != nil {
		assert.FailNow(t, err.Error(), "Connect")
	}

	api.OnError(func(err error) {
		assert.FailNow(t, err.Error(), "OnError")
	})

	block := uint64(26878298)

	for {
		fmt.Println("get block: ", block)
		res, err := api.GetBlock(block)
		if err != nil {
			assert.FailNow(t, err.Error(), "GetBlock")
		}

		assert.NotNil(t, res)
		block++
	}

	//util.Dump("get_block >", res)
}
