package latency

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/denkhaus/bitshares/tests"
	"github.com/stretchr/testify/assert"
)

func Test_LatencyAnalyzer(t *testing.T) {

	ctx := context.Background()
	//ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	lat, err := NewLatencyTester(ctx, tests.RpcFullApiUrl, 10*time.Second)
	if err != nil {
		assert.FailNow(t, err.Error(), "NewLatencyTester")
	}

	lat.Start()
	time.Sleep(3000 * time.Second)
	lat.Stop()
	//<-lat.Done()
	fmt.Print(lat.String())
}
