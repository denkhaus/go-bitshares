package latency

import (
	"sync"

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
	latency  uint64
	attempts uint64
	endpoint string
}

func (p *NodeStats) onError(err error) {

}

func NewNodeStats(wsRPCEndpoint, walletRPCEndpoint string) (*NodeStats, error) {
	stats := NodeStats{
		endpoint: wsRPCEndpoint,
	}

	api := api.New(stats.endpoint, "walletRPCEndpoint")
	if err := api.Connect(); err != nil {
		return nil, errors.Annotate(err, "Connect")
	}

	api.OnError(stats.onError)
	stats.api = api

	return &stats, nil
}

type LatencyTester struct {
	sync.Mutex
	statsMap          map[string]*NodeStats
	walletRPCEndpoint string
}

func NewLatencyTester(walletRPCEndpoint string) (*LatencyTester, error) {
	lat := LatencyTester{
		statsMap:          make(map[string]*NodeStats),
		walletRPCEndpoint: walletRPCEndpoint,
	}

	for _, ep := range knownEndpoints {
		stat, err := NewNodeStats(ep, walletRPCEndpoint)
		if err != nil {
			return nil, errors.Annotate(err, "NewNodeStats")
		}
		lat.statsMap[ep] = stat
	}

	return &lat, nil
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

func (p *LatencyTester) Start() error {

	return nil
}
