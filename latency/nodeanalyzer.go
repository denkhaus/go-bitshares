package latency

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/denkhaus/bitshares/api"
	"github.com/juju/errors"
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
		"wss://la.dexnode.net/ws",
		"wss://openledger.hk/ws",
		"wss://sg.nodes.bitshares.ws",
		"wss://bts.open.icowallet.net/ws",
		"wss://ws.gdex.io",
		"wss://bitshares-api.wancloud.io/ws",
		"wss://ws.hellobts.com/",
		"wss://bitshares.dacplay.org/ws",
		"wss://crazybit.online",
		"wss://kimziv.com/ws",
		"wss://wss.ioex.top",
		"wss://node.btscharts.com/ws",
		"wss://bts-seoul.clockwork.gr/",
		"wss://bitshares.cyberit.io/",
		"wss://api.btsgo.net/ws",
		"wss://ws.winex.pro",
		"wss://bts.to0l.cn:4443/ws",
		"wss://bitshares.bts123.cc:15138/",
		"wss://bit.btsabc.org/ws",
		"wss://ws.gdex.top",
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

func (p *NodeStats) Score() time.Duration {
	if p.attempts > 0 {
		return time.Duration(int64(p.latency) / p.attempts)
	}

	return -1
}

func (p *NodeStats) String() string {
	return fmt.Sprintf("ep: %s | attempts: %d | errors: %d | score: %s",
		p.endpoint, p.attempts, p.errors, p.Score())
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

func (p *NodeStats) check() {
	tm := time.Now()
	_, err := p.api.GetDynamicGlobalProperties()
	if err != nil {
		p.errors++
	}

	p.latency += time.Since(tm)
	p.attempts++
}

type LatencyTester struct {
	sync.Mutex
	wg                sync.WaitGroup
	ticker            *time.Ticker
	ctx               context.Context
	cancel            context.CancelFunc
	statsMap          map[string]*NodeStats
	walletRPCEndpoint string
}

func NewLatencyTester(ctx context.Context, walletRPCEndpoint string, checkDur time.Duration) (*LatencyTester, error) {
	cctx, cancel := context.WithCancel(ctx)
	lat := LatencyTester{
		statsMap:          make(map[string]*NodeStats),
		walletRPCEndpoint: walletRPCEndpoint,
		ticker:            time.NewTicker(checkDur),
		cancel:            cancel,
		ctx:               cctx,
	}

	for _, ep := range knownEndpoints {
		stat, err := NewNodeStats(ep, walletRPCEndpoint)
		if err == nil {
			lat.Lock()
			lat.statsMap[ep] = stat
			lat.Unlock()
		}
	}

	return &lat, nil
}

func (p *LatencyTester) String() string {
	p.Lock()
	defer p.Unlock()
	builder := strings.Builder{}
	for _, stats := range p.statsMap {
		builder.WriteString(stats.String())
		builder.WriteString("\n")
	}

	return builder.String()
}

func (p *LatencyTester) AddEndpoint(ep string) error {
	p.Lock()
	defer p.Unlock()

	if _, ok := p.statsMap[ep]; !ok {
		stat, err := NewNodeStats(ep, p.walletRPCEndpoint)
		if err != nil {
			return errors.Annotate(err, "NewNodeStats")
		}
		p.statsMap[ep] = stat
	}

	return nil
}

func (p *LatencyTester) check() {
	p.Lock()
	defer p.Unlock()

	for _, stats := range p.statsMap {
		go stats.check()
	}
}

func (p *LatencyTester) Done() <-chan struct{} {
	return p.ctx.Done()
}

func (p *LatencyTester) Start() {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		select {
		case <-p.ticker.C:
			p.check()
		case <-p.ctx.Done():
			return
		}
	}()
}

func (p *LatencyTester) Stop() {
	p.ticker.Stop()
	p.cancel()
	p.wg.Wait()
}
