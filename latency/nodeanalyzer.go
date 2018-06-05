package latency

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/denkhaus/bitshares/api"
	"github.com/juju/errors"
	"gopkg.in/tomb.v2"
)

var (
	knownEndpoints = []string{
		"wss://eu.openledger.info/ws",
		"wss://bitshares.openledger.info/ws",
		"wss://dex.rnglab.org",
		"wss://api.bitshares.bhuz.info/ws",
		"wss://bitshares.crypto.fans/ws",
		"wss://node.market.rudex.org",
		"wss://api.bts.blckchnd.com",
		"wss://eu.nodes.bitshares.ws",
		"wss://btsws.roelandp.nl/ws",
		"wss://btsfullnode.bangzi.info/ws",
		"wss://api-ru.bts.blckchnd.com",
		"wss://kc-us-dex.xeldal.com/ws",
		"wss://api.btsxchng.com",
		"wss://api.bts.network",
		"wss://dexnode.net/ws",
		"wss://us.nodes.bitshares.ws",
		"wss://api.bts.mobi/ws",
		"wss://blockzms.xyz/ws",
		"wss://bts-api.lafona.net/ws",
		"wss://api.bts.ai/",
		// "wss://la.dexnode.net/ws",
		// "wss://openledger.hk/ws",
		// "wss://sg.nodes.bitshares.ws",
		// "wss://bts.open.icowallet.net/ws",
		// "wss://ws.gdex.io",
		// "wss://bitshares-api.wancloud.io/ws",
		// "wss://ws.hellobts.com/",
		// "wss://bitshares.dacplay.org/ws",
		// "wss://crazybit.online",
		// "wss://kimziv.com/ws",
		// "wss://wss.ioex.top",
		// "wss://node.btscharts.com/ws",
		// "wss://bts-seoul.clockwork.gr/",
		// "wss://bitshares.cyberit.io/",
		// "wss://api.btsgo.net/ws",
		// "wss://ws.winex.pro",
		// "wss://bts.to0l.cn:4443/ws",
		// "wss://bitshares.bts123.cc:15138/",
		// "wss://bit.btsabc.org/ws",
		// "wss://ws.gdex.top",
	}
)

type NodeStats struct {
	api      api.BitsharesAPI
	latency  time.Duration
	attempts int64
	errors   int64
	endpoint string
}

func (p *NodeStats) onError(err error) {
	fmt.Println(err)
	p.errors++
}

func (p *NodeStats) Latency() time.Duration {
	if p.attempts > 0 {
		return time.Duration(int64(p.latency) / p.attempts)
	}

	return 0
}

func (p *NodeStats) Score() int64 {
	lat := int64(p.Latency())
	if p.errors == 0 {
		return lat
	}

	return lat * p.errors
}

func (p *NodeStats) String() string {
	return fmt.Sprintf("ep: %s | attempts: %d | errors: %d | latency: %s | score: %d",
		p.endpoint, p.attempts, p.errors, p.Latency(), p.Score())
}

func NewNodeStats(wsRPCEndpoint, walletRPCEndpoint string) (*NodeStats, error) {
	api := api.New(wsRPCEndpoint, walletRPCEndpoint)
	if err := api.Connect(); err != nil {
		return nil, errors.Annotate(err, "Connect")
	}

	stats := &NodeStats{
		endpoint: wsRPCEndpoint,
		api:      api,
	}

	api.OnError(stats.onError)
	return stats, nil
}

func (p *NodeStats) check() error {
	tm := time.Now()
	_, err := p.api.GetDynamicGlobalProperties()
	if err != nil {
		p.errors++
		return errors.Annotate(err, "GetDynamicGlobalProperties")
	}

	p.attempts++
	p.latency += time.Since(tm)
	fmt.Printf("checked stats for %s\n", p)
	return nil
}

type LatencyTester struct {
	sync.Mutex
	tmb               *tomb.Tomb
	statsMap          map[string]*NodeStats
	walletRPCEndpoint string
}

func NewLatencyTester(ctx context.Context, walletRPCEndpoint string) (*LatencyTester, error) {
	tmb, _ := tomb.WithContext(ctx)
	lat := LatencyTester{
		statsMap:          make(map[string]*NodeStats),
		walletRPCEndpoint: walletRPCEndpoint,
		tmb:               tmb,
	}

	for _, ep := range knownEndpoints {
		stat, err := NewNodeStats(ep, walletRPCEndpoint)
		if err == nil {
			lat.statsMap[ep] = stat
		}
	}

	return &lat, nil
}

func (p *LatencyTester) String() string {
	builder := strings.Builder{}

	p.Lock()
	defer p.Unlock()
	for _, stats := range p.statsMap {
		builder.WriteString(stats.String())
		builder.WriteString("\n")
	}

	return builder.String()
}

// func (p *LatencyTester) AddEndpoint(ep string) error {
// 	p.Lock()
// 	defer p.Unlock()

// 	if _, ok := p.statsMap[ep]; !ok {
// 		stat, err := NewNodeStats(ep, p.walletRPCEndpoint)
// 		if err != nil {
// 			return errors.Annotate(err, "NewNodeStats")
// 		}
// 		p.statsMap[ep] = stat
// 	}

// 	return nil
// }

func (p *LatencyTester) Done() <-chan struct{} {
	return p.tmb.Dead()
}

func (p *LatencyTester) Start() {
	p.tmb.Go(func() error {
		var cnt int32

		for {
			for _, stats := range p.statsMap {
				select {
				case <-p.tmb.Dying():
					return tomb.ErrDying
				default:
				}

				for atomic.LoadInt32(&cnt) > 3 {
					time.Sleep(1 * time.Second)
				}

				st := stats
				p.tmb.Go(func() error {
					atomic.AddInt32(&cnt, 1)
					defer atomic.AddInt32(&cnt, -1)
					return st.check()
				})
			}
		}
	})
}

func (p *LatencyTester) Stop() error {
	p.tmb.Kill(nil)
	return p.tmb.Wait()
}
