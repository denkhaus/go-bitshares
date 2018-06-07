package latency

import (
	"context"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/denkhaus/bitshares/client"
	sort "github.com/emirpasic/gods/utils"
	"gopkg.in/tomb.v2"
)

var (
	//LoopSeconds = time for one pass to calc dynamic delay
	LoopSeconds = 60
	//our known node endpoints
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

type LatencyTester interface {
	Start()
	Close() error
	String() string
	AddEndpoint(ep string)
	OnTopNodeChanged(fn func(topNodeEndpoint string))
	TopNodeEndpoint() string
	TopNodeClient() client.WebsocketClient
	Done() <-chan struct{}
}

//NodeStats holds stat data for each endpoint
type NodeStats struct {
	cli      client.WebsocketClient
	latency  time.Duration
	attempts int64
	errors   int64
	endpoint string
}

func (p *NodeStats) onError(err error) {
	p.errors++
}

//Latency returns the nodes latency
func (p *NodeStats) Latency() time.Duration {
	if p.attempts > 0 {
		return time.Duration(int64(p.latency) / p.attempts)
	}

	return 0
}

//Score returns reliability score for each node. The less the better.
func (p *NodeStats) Score() int64 {
	lat := int64(p.Latency())
	if lat == 0 {
		return math.MaxInt64
	}
	if p.errors == 0 {
		return lat
	}

	return lat * p.errors
}

// String returns the stats string representation
func (p *NodeStats) String() string {
	return fmt.Sprintf("ep: %s | attempts: %d | errors: %d | latency: %s | score: %d",
		p.endpoint, p.attempts, p.errors, p.Latency(), p.Score())
}

//NewNodeStats creates a new stat object
func NewNodeStats(wsRPCEndpoint string) *NodeStats {
	stats := &NodeStats{
		endpoint: wsRPCEndpoint,
		cli:      client.NewWebsocketClient(wsRPCEndpoint),
	}

	stats.cli.OnError(stats.onError)
	return stats
}

func (p *NodeStats) Equals(n *NodeStats) bool {
	return p.endpoint == n.endpoint
}

func (p *NodeStats) check() {
	if err := p.cli.Connect(); err != nil {
		p.errors++
		return
	}
	defer p.cli.Close()

	tm := time.Now()
	_, err := p.cli.CallAPI(1, "login", "", "")
	if err != nil {
		p.errors++
		return
	}

	p.latency += time.Since(tm)
	p.attempts++
}

type latencyTester struct {
	mu               sync.Mutex
	tmb              *tomb.Tomb
	toApply          []string
	fallbackURL      string
	onTopNodeChanged func(string)
	stats            []interface{}
	pass             int
}

func NewLatencyTester(fallbackURL string) (LatencyTester, error) {
	return NewLatencyTesterWithContext(context.Background(), fallbackURL)
}

func NewLatencyTesterWithContext(ctx context.Context, fallbackURL string) (LatencyTester, error) {
	tmb, _ := tomb.WithContext(ctx)
	lat := latencyTester{
		fallbackURL: fallbackURL,
		stats:       make([]interface{}, 0, len(knownEndpoints)),
		tmb:         tmb,
	}

	lat.createStats(knownEndpoints)
	return &lat, nil
}

func (p *latencyTester) String() string {
	builder := strings.Builder{}

	p.mu.Lock()
	defer p.mu.Unlock()
	for _, s := range p.stats {
		stat := s.(*NodeStats)
		builder.WriteString(stat.String())
		builder.WriteString("\n")
	}

	return builder.String()
}

func (p *latencyTester) OnTopNodeChanged(fn func(string)) {
	p.onTopNodeChanged = fn
}

//AddEndpoint adds a new endpoint while the latencyTester is running
func (p *latencyTester) AddEndpoint(ep string) {
	p.toApply = append(p.toApply, ep)
}

func (p *latencyTester) sortResults() {
	p.mu.Lock()
	defer p.mu.Unlock()

	oldTop := p.stats[0].(*NodeStats)
	sort.Sort(p.stats, func(a, b interface{}) int {
		sa := a.(*NodeStats).Score()
		sb := b.(*NodeStats).Score()
		if sa > sb {
			return 1
		}

		if sa < sb {
			return -1
		}

		return 0
	})

	newTop := p.stats[0].(*NodeStats)
	if !oldTop.Equals(newTop) {
		if p.onTopNodeChanged != nil {
			//use goroutine here to avoid deadlock
			go p.onTopNodeChanged(newTop.endpoint)
		}
	}
}

func (p *latencyTester) createStats(eps []string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, ep := range eps {
		found := false
		for _, s := range p.stats {
			stat := s.(*NodeStats)
			if stat.endpoint == ep {
				found = true
			}
		}

		if !found {
			p.stats = append(
				p.stats,
				NewNodeStats(ep),
			)
		}
	}
}

//TopNodeEndpoint returns the fastest endpoint URL. If the tester has no validated results
//your given fallback endpoint is returned.
func (p *latencyTester) TopNodeEndpoint() string {
	if p.pass > 0 {
		p.mu.Lock()
		defer p.mu.Unlock()
		st := p.stats[0].(*NodeStats)
		return st.endpoint
	}

	return p.fallbackURL
}

//TopNodeClient returns a new WebsocketClient to connect to the fastest node.
//If the tester has no validated results, a client with your given
//fallback endpoint is returned. You need to call Connect for yourself.
func (p *latencyTester) TopNodeClient() client.WebsocketClient {
	return client.NewWebsocketClient(
		p.TopNodeEndpoint(),
	)
}

// Done returns the channel that can be used to wait until
// the tester has finished.
func (p *latencyTester) Done() <-chan struct{} {
	return p.tmb.Dead()
}

//Start starts the testing process
func (p *latencyTester) Start() {
	p.tmb.Go(func() error {
		for {
			//apply later incoming endpoints
			p.createStats(p.toApply)
			// dynamic sleep time
			slp := time.Duration(LoopSeconds/len(p.stats)) * time.Second
			for i := 0; i < len(p.stats); i++ {
				select {
				case <-p.tmb.Dying():
					p.sortResults()
					return tomb.ErrDying
				default:
					idx := i
					time.Sleep(slp)
					p.tmb.Go(func() error {
						p.mu.Lock()
						defer p.mu.Unlock()

						st := p.stats[idx].(*NodeStats)
						st.check()
						return nil
					})
				}
			}

			p.sortResults()
			p.pass++
		}
	})
}

//Close stops the tester and waits until all goroutines have finished.
func (p *latencyTester) Close() error {
	p.tmb.Kill(nil)
	return p.tmb.Wait()
}
