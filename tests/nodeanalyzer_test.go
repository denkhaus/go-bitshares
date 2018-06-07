package tests

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/denkhaus/bitshares/latency"
	"github.com/stretchr/testify/assert"
)

func Test_LatencyAnalyzerWithTimeout(t *testing.T) {
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, 60*time.Second)
	lat, err := latency.NewLatencyTesterWithContext(ctx, WsFullApiUrl)
	if err != nil {
		assert.FailNow(t, err.Error(), "NewLatencyTester")
	}

	lat.Start()
	<-lat.Done()
	log.Print(lat.String())
}

func Test_LatencyAnalyzerWithStop(t *testing.T) {
	lat, err := latency.NewLatencyTester(WsFullApiUrl)
	if err != nil {
		assert.FailNow(t, err.Error(), "NewLatencyTester")
	}

	lat.Start()
	time.Sleep(16 * time.Second)

	if err := lat.Close(); err != nil {
		assert.FailNow(t, err.Error(), "Close")
	}

	log.Print(lat.String())
}
