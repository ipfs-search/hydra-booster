package hydra

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/axiomhq/hyperloglog"
	"github.com/ipfs/go-datastore"
	levelds "github.com/ipfs/go-ds-leveldb"
	"github.com/libp2p/go-libp2p-core/network"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/hydra-booster/metrics"
	"github.com/libp2p/hydra-booster/sybil"
	"github.com/libp2p/hydra-booster/sybil/opts"
	"github.com/multiformats/go-multiaddr"
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
)

// Hydra is a container for heads (sybils) and their shared belly bits.
type Hydra struct {
	Sybils          []*sybil.Sybil
	SharedDatastore datastore.Datastore
	// SharedRoutingTable *kbucket.RoutingTable

	hyperLock       *sync.Mutex
	hyperlog        *hyperloglog.Sketch
	periodicMetrics *PeriodicMetrics
	stopped         bool
}

// Options are configuration for a new hydra
type Options struct {
	DatastorePath string
	GetPort       func() int
	NSybils       int
	BucketSize    int
	BsCon         int
	Relay         bool
	Stagger       time.Duration
}

// NewHydra creates a new Hydra with the passed options.
func NewHydra(options Options) (*Hydra, error) {
	datastore, err := levelds.NewDatastore(options.DatastorePath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create datastore: %w", err)
	}

	var sybils []*sybil.Sybil

	fmt.Fprintf(os.Stderr, "Running %d DHT Sybils:\n", options.NSybils)

	var hyperLock sync.Mutex
	hyperlog := hyperloglog.New()

	// What is a limiter?
	limiter := make(chan struct{}, options.BsCon)

	for i := 0; i < options.NSybils; i++ {
		time.Sleep(options.Stagger)
		fmt.Fprintf(os.Stderr, ".")

		addr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", options.GetPort()))
		syb, bsCh, err := sybil.NewSybil(
			opts.Datastore(datastore),
			opts.Addr(addr),
			opts.Relay(options.Relay),
			opts.BucketSize(options.BucketSize),
			opts.Limiter(limiter),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to spawn node with swarm address %v: %w", addr, err)
		}

		ctx, err := tag.New(context.Background(), tag.Insert(metrics.KeyPeerID, syb.Host.ID().String()))
		if err != nil {
			return nil, err
		}

		stats.Record(ctx, metrics.Sybils.M(1))

		syb.Host.Network().Notify(&network.NotifyBundle{
			ConnectedF: func(n network.Network, v network.Conn) {
				hyperLock.Lock()
				hyperlog.Insert([]byte(v.RemotePeer()))
				hyperLock.Unlock()
				stats.Record(ctx, metrics.ConnectedPeers.M(1))
			},
			DisconnectedF: func(n network.Network, v network.Conn) {
				stats.Record(ctx, metrics.ConnectedPeers.M(-1))
			},
		})

		go handleBootstrapStatus(ctx, bsCh)

		sybils = append(sybils, syb)
	}
	fmt.Fprintf(os.Stderr, "\n")

	hydra := Hydra{
		Sybils:          sybils,
		SharedDatastore: datastore,
		hyperLock:       &hyperLock,
		hyperlog:        hyperlog,
	}

	hydra.periodicMetrics = NewPeriodicMetrics(&hydra, time.Second*5)

	return &hydra, nil
}

// Stop halts the hydra node
func (hy *Hydra) Stop() {
	hy.periodicMetrics.Stop()
	for _, syb := range hy.Sybils {
		err := syb.Routing.(*dht.IpfsDHT).Close()
		if err != nil {
			fmt.Println(err)
		}
		err = syb.Host.Close()
		if err != nil {
			fmt.Println(err)
		}
	}
	hy.SharedDatastore.Close()
}

func handleBootstrapStatus(ctx context.Context, ch chan sybil.BootstrapStatus) {
	for status := range ch {
		if status.Err != nil {
			fmt.Println(status.Err)
		}
		if status.Done {
			stats.Record(ctx, metrics.BootstrappedSybils.M(1))
		}
	}
}

// GetUniquePeersCount retrieves the current total for unique peers
func (hy *Hydra) GetUniquePeersCount() uint64 {
	hy.hyperLock.Lock()
	defer hy.hyperLock.Unlock()
	return hy.hyperlog.Estimate()
}