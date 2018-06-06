package latency

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/denkhaus/bitshares/tests"
	"github.com/stretchr/testify/assert"
)

func Test_LatencyAnalyzerWithTimeout(t *testing.T) {

	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, 60*time.Second)
	lat, err := NewLatencyTester(ctx, tests.RpcFullApiUrl)
	if err != nil {
		assert.FailNow(t, err.Error(), "NewLatencyTester")
	}

	lat.Start()
	<-lat.Done()
	fmt.Print(lat.String())
}

func Test_LatencyAnalyzerWithStop(t *testing.T) {
	lat, err := NewLatencyTester(context.Background(), tests.RpcFullApiUrl)
	if err != nil {
		assert.FailNow(t, err.Error(), "NewLatencyTester")
	}

	lat.Start()
	time.Sleep(16 * time.Second)

	if err := lat.Stop(); err != nil {
		assert.FailNow(t, err.Error(), "Stop")
	}

	fmt.Print(lat.String())
}
