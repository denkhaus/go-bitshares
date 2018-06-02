package latency

import (
	"sync"

	"github.com/denkhaus/bitshares/api"
	"github.com/juju/errors"
)

var (
	statsMutex sync.Mutex
	statsMap   = map[string]*NodeStats{
		"wss://eu.openledger.info/ws":        nil,
		"wss://bitshares.openledger.info/ws": nil,
		"wss://dex.rnglab.org":               nil,
		"wss://api.bitshares.bhuz.info/ws":   nil,
		"wss://bitshares.crypto.fans/ws":     nil,
		"wss://node.market.rudex.org":        nil,
		"wss://api.bts.blckchnd.com":         nil,
		"wss://eu.nodes.bitshares.ws":        nil,
		"wss://btsws.roelandp.nl/ws":         nil,
		"wss://btsfullnode.bangzi.info/ws":   nil,
		"wss://api-ru.bts.blckchnd.com":      nil,
		"wss://kc-us-dex.xeldal.com/ws":      nil,
		"wss://api.btsxchng.com":             nil,
		"wss://api.bts.network":              nil,
		"wss://dexnode.net/ws":               nil,
		"wss://us.nodes.bitshares.ws":        nil,
		"wss://api.bts.mobi/ws":              nil,
		"wss://blockzms.xyz/ws":              nil,
		"wss://bts-api.lafona.net/ws":        nil,
		"wss://api.bts.ai/":                  nil,
		"wss://la.dexnode.net/ws":            nil,
		"wss://openledger.hk/ws":             nil,
		"wss://sg.nodes.bitshares.ws":        nil,
		"wss://bts.open.icowallet.net/ws":    nil,
		"wss://ws.gdex.io":                   nil,
		"wss://bitshares-api.wancloud.io/ws": nil,
		"wss://ws.hellobts.com/":             nil,
		"wss://bitshares.dacplay.org/ws":     nil,
		"wss://crazybit.online":              nil,
		"wss://kimziv.com/ws":                nil,
		"wss://wss.ioex.top":                 nil,
		"wss://node.btscharts.com/ws":        nil,
		"wss://bts-seoul.clockwork.gr/":      nil,
		"wss://bitshares.cyberit.io/":        nil,
		"wss://api.btsgo.net/ws":             nil,
		"wss://ws.winex.pro":                 nil,
		"wss://bts.to0l.cn:4443/ws":          nil,
		"wss://bitshares.bts123.cc:15138/":   nil,
		"wss://bit.btsabc.org/ws":            nil,
		"wss://ws.gdex.top":                  nil,
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
}

func NewLatencyTester(walletRPCEndpoint string) (*LatencyTester, error) {
	statsMutex.Lock()
	defer statsMutex.Unlock()

	lat := LatencyTester{}

	for ep := range statsMap {
		stat, err := NewNodeStats(ep, walletRPCEndpoint)
		if err != nil {
			return nil, errors.Annotate(err, "Connect")
		}
		statsMap[ep] = stat
	}

	return &lat, nil
}

func (*LatencyTester) test() {

}
