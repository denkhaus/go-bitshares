package tests

import (
	"context"
	"testing"
	"time"

	"github.com/denkhaus/bitshares/api"
	"github.com/stretchr/testify/assert"
)

//long running test
func Test_LatencyAnalyzerWithTimeout(t *testing.T) {
	t.Log("create tester")

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	lat, err := api.NewLatencyTesterWithContext(ctx, WsFullApiUrl)
	if err != nil {
		assert.FailNow(t, err.Error(), "NewLatencyTester")
	}

	t.Log("start tester")
	lat.Start()
	t.Log("wait")
	<-lat.Done()

	t.Log("\n", lat.String())
}

//long running test
func Test_LatencyAnalyzerWithClose(t *testing.T) {
	t.Log("create tester")
	lat, err := api.NewLatencyTester(WsFullApiUrl)
	if err != nil {
		assert.FailNow(t, err.Error(), "NewLatencyTester")
	}

	t.Log("start tester")
	lat.Start()
	t.Log("wait")
	time.Sleep(16 * time.Second)

	t.Log("close tester")
	if err := lat.Close(); err != nil {
		assert.FailNow(t, err.Error(), "Close")
	}

	t.Log("\n", lat.String())
}

//long running test
func Test_BestNodeAPI(t *testing.T) {
	t.Log("startup API")
	api, err := api.NewWebsocketAPIWithAutoEndpoint(WsFullApiUrl)
	if err != nil {
		assert.FailNow(t, err.Error(), "NewWithAutoEndpoint")
	}

	api.OnError(func(err error) {
		assert.FailNow(t, err.Error(), "OnError")
	})

	// run this test at least 70 seconds to force
	//the latency tester to provide its first new results.
	for i := 0; i < 70; i++ {
		t.Log("invoke API")
		_, err = api.CallWsAPI(1, "login", "", "")
		if err != nil {
			assert.FailNow(t, err.Error(), "CallWsAPI")
		}
		time.Sleep(1 * time.Second)
	}

	t.Log("close API")
	if err := api.Close(); err != nil {
		assert.FailNow(t, err.Error(), "Close")
	}
}
