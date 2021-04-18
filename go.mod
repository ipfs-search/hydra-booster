module github.com/libp2p/hydra-booster

require (
	contrib.go.opencensus.io/exporter/prometheus v0.2.0
	github.com/alanshaw/ipfs-ds-postgres v0.1.0
	github.com/alanshaw/ipfs-hookds v0.3.0
	github.com/alanshaw/prom-metrics-client v0.3.0
	github.com/axiomhq/hyperloglog v0.0.0-20191112132149-a4c4c47bc57f
	github.com/btcsuite/btcd v0.21.0-beta // indirect
	github.com/dustin/go-humanize v1.0.0
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp/go-multierror v1.1.0
	github.com/hnlq715/golang-lru v0.2.1-0.20200422024707-82ba7badf9a6
	github.com/ipfs-search/ipfs-search v0.0.0-20210415163049-0ad29fdf8e0f
	github.com/ipfs/go-cid v0.0.7
	github.com/ipfs/go-datastore v0.4.5
	github.com/ipfs/go-ds-leveldb v0.4.2
	github.com/ipfs/go-ipns v0.0.2
	github.com/jackc/pgx/v4 v4.9.0
	github.com/libp2p/go-libp2p v0.11.0
	github.com/libp2p/go-libp2p-circuit v0.3.1
	github.com/libp2p/go-libp2p-connmgr v0.2.4
	github.com/libp2p/go-libp2p-core v0.6.1
	github.com/libp2p/go-libp2p-kad-dht v0.10.0
	github.com/libp2p/go-libp2p-kbucket v0.4.7
	github.com/libp2p/go-libp2p-noise v0.1.2
	github.com/libp2p/go-libp2p-peerstore v0.2.6
	github.com/libp2p/go-libp2p-quic-transport v0.8.2
	github.com/libp2p/go-libp2p-record v0.1.3
	github.com/libp2p/go-libp2p-tls v0.1.3
	github.com/libp2p/go-tcp-transport v0.2.1
	github.com/multiformats/go-base32 v0.0.3
	github.com/multiformats/go-multiaddr v0.3.1
	github.com/multiformats/go-multihash v0.0.14
	github.com/prometheus/client_golang v1.7.1
	github.com/streadway/amqp v0.0.0-20200108173154-1c71cc93ed71
	github.com/whyrusleeping/timecache v0.0.0-20160911033111-cfcb2f1abfee
	go.opencensus.io v0.22.5
	go.uber.org/zap v1.16.0 // indirect
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a // indirect
	golang.org/x/sys v0.0.0-20200909081042-eff7692f9009 // indirect
)

go 1.14
