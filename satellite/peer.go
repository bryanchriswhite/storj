// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package satellite

import (
	"context"
	"fmt"
	"net"
	"path/filepath"

	"github.com/zeebo/errs"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"storj.io/storj/pkg/accounting"
	"storj.io/storj/pkg/bwagreement"
	"storj.io/storj/pkg/datarepair/checker"
	"storj.io/storj/pkg/datarepair/irreparable"
	"storj.io/storj/pkg/datarepair/queue"
	"storj.io/storj/pkg/datarepair/repairer"
	"storj.io/storj/pkg/discovery"
	"storj.io/storj/pkg/identity"
	"storj.io/storj/pkg/kademlia"
	"storj.io/storj/pkg/node"
	"storj.io/storj/pkg/overlay"
	"storj.io/storj/pkg/pb"
	"storj.io/storj/pkg/pointerdb"
	"storj.io/storj/pkg/server"
	"storj.io/storj/pkg/statdb"
	"storj.io/storj/pkg/storj"
	"storj.io/storj/satellite/console"
	"storj.io/storj/storage"
	"storj.io/storj/storage/boltdb"
	"storj.io/storj/storage/storelogger"
)

// DB is the master database for the satellite
type DB interface {
	// CreateTables initializes the database
	CreateTables() error
	// Close closes the database
	Close() error

	// BandwidthAgreement returns database for storing bandwidth agreements
	BandwidthAgreement() bwagreement.DB
	// StatDB returns database for storing node statistics
	StatDB() statdb.DB
	// OverlayCache returns database for caching overlay information
	OverlayCache() overlay.DB
	// Accounting returns database for storing information about data use
	Accounting() accounting.DB
	// RepairQueue returns queue for segments that need repairing
	RepairQueue() queue.RepairQueue
	// Irreparable returns database for failed repairs
	Irreparable() irreparable.DB
	// Console returns database for satellite console
	Console() console.DB
}

// Config is the global config satellite
type Config struct {
	Identity identity.Config

	// TODO: switch to using server.Config when Identity has been removed from it
	Database      string `help:"satellite database connection string" default:"sqlite3://$CONFDIR/master.db"`
	PublicAddress string `help:"public address to listen on" default:":7777"`

	Kademlia  kademlia.Config
	Overlay   overlay.Config
	Discovery discovery.Config

	PointerDB   pointerdb.Config
	BwAgreement bwagreement.Config

	Checker  checker.Config
	Repairer repairer.Config
	// TODO: Audit    audit.Config
}

// Peer is the satellite
type Peer struct {
	// core dependencies
	Log      *zap.Logger
	Identity *identity.FullIdentity
	DB       DB

	// servers
	Public struct {
		Listener net.Listener
		Server   *server.Server
	}

	// services and endpoints
	Kademlia struct {
		RoutingTable *kademlia.RoutingTable
		Service      *kademlia.Kademlia
		Endpoint     *node.Server
	}

	Overlay struct {
		Service  *overlay.Cache
		Endpoint *overlay.Server
	}

	Discovery struct {
		Service *discovery.Discovery
	}

	Metainfo struct {
		Database storage.KeyValueStore // TODO: move into pointerDB
		Service  *pointerdb.Service
		Endpoint *pointerdb.Server
	}

	Agreements struct {
		Endpoint *bwagreement.Server
	}

	Repair struct {
		Checker  checker.Checker // TODO: convert to actual struct
		Repairer *repairer.Service
	}
	Audit struct {
		// TODO: Service *audit.Service
	}

	// TODO: add console
}

// New creates a new satellite
func New(log *zap.Logger, full *identity.FullIdentity, db DB, config *Config) (*Peer, error) {
	peer := &Peer{
		Log:      log,
		Identity: full,
		DB:       db,
	}

	var err error

	{ // setup listener and server
		peer.Public.Listener, err = net.Listen("tcp", config.PublicAddress)
		if err != nil {
			return nil, errs.Combine(err, peer.Close())
		}

		publicConfig := server.Config{Address: peer.Public.Listener.Addr().String()}
		publicOptions, err := server.NewOptions(peer.Identity, publicConfig)
		if err != nil {
			return nil, errs.Combine(err, peer.Close())
		}

		peer.Public.Server, err = server.NewServer(publicOptions, peer.Public.Listener, nil)
		if err != nil {
			return nil, errs.Combine(err, peer.Close())
		}
	}

	{ // setup kademlia
		config := config.Kademlia
		// TODO: move this setup logic into kademlia package
		if config.ExternalAddress == "" {
			config.ExternalAddress = peer.Public.Server.Addr().String()
		}

		self := pb.Node{
			Id:   peer.ID(),
			Type: pb.NodeType_SATELLITE,
			Address: &pb.NodeAddress{
				Address: config.ExternalAddress,
			},
			Metadata: &pb.NodeMetadata{
				Email:  config.Operator.Email,
				Wallet: config.Operator.Wallet,
			},
		}

		{ // setup routing table
			// TODO: clean this up
			bucketIdentifier := peer.ID().String()[:5] // need a way to differentiate between nodes if running more than one simultaneously
			dbpath := filepath.Join(config.DBPath, fmt.Sprintf("kademlia_%s.db", bucketIdentifier))

			dbs, err := boltdb.NewShared(dbpath, kademlia.KademliaBucket, kademlia.NodeBucket)
			if err != nil {
				return nil, errs.Combine(err, peer.Close())
			}
			kdb, ndb := dbs[0], dbs[1]

			peer.Kademlia.RoutingTable, err = kademlia.NewRoutingTable(peer.Log.Named("routing"), self, kdb, ndb)
			if err != nil {
				return nil, errs.Combine(err, peer.Close())
			}
		}

		// TODO: reduce number of arguments
		peer.Kademlia.Service, err = kademlia.NewWith(peer.Log.Named("kademlia"), self, nil, peer.Identity, config.Alpha, peer.Kademlia.RoutingTable)
		if err != nil {
			return nil, errs.Combine(err, peer.Close())
		}

		peer.Kademlia.Endpoint = node.NewServer(peer.Log.Named("kademlia:endpoint"), peer.Kademlia.Service)
		pb.RegisterNodesServer(peer.Public.Server.GRPC(), peer.Kademlia.Endpoint)
	}

	{ // setup overlay
		config := config.Overlay
		peer.Overlay.Service = overlay.NewCache(peer.DB.OverlayCache(), peer.DB.StatDB())

		ns := &pb.NodeStats{
			UptimeCount:       config.Node.UptimeCount,
			UptimeRatio:       config.Node.UptimeRatio,
			AuditSuccessRatio: config.Node.AuditSuccessRatio,
			AuditCount:        config.Node.AuditCount,
		}

		peer.Overlay.Endpoint = overlay.NewServer(peer.Log.Named("overlay:endpoint"), peer.Overlay.Service, ns)
		pb.RegisterOverlayServer(peer.Public.Server.GRPC(), peer.Overlay.Endpoint)
	}

	{ // setup discovery
		config := config.Discovery
		peer.Discovery.Service = discovery.New(peer.Log.Named("discovery"), peer.Overlay.Service, peer.Kademlia.Service, peer.DB.StatDB(), config.RefreshInterval)
	}

	{ // setup metainfo
		db, err := pointerdb.NewStore(config.PointerDB.DatabaseURL)
		if err != nil {
			return nil, errs.Combine(err, peer.Close())
		}

		peer.Metainfo.Database = storelogger.New(peer.Log.Named("pdb"), db)
		peer.Metainfo.Service = pointerdb.NewService(peer.Log.Named("pointerdb"), peer.Metainfo.Database)
		peer.Metainfo.Endpoint = pointerdb.NewServer(peer.Log.Named("pointerdb:endpoint"), peer.Metainfo.Service, peer.Overlay.Service, config.PointerDB, peer.Identity)
		pb.RegisterPointerDBServer(peer.Public.Server.GRPC(), peer.Metainfo.Endpoint)
	}

	{ // setup agreements
		peer.Agreements.Endpoint = bwagreement.NewServer(peer.DB.BandwidthAgreement(), peer.Log.Named("agreements"), peer.Identity.Leaf.PublicKey)
		pb.RegisterBandwidthServer(peer.Public.Server.GRPC(), peer.Agreements.Endpoint)
	}

	{ // setup datarepair
		// TODO: simplify argument list somehow
		peer.Repair.Checker = checker.NewChecker(
			peer.Metainfo.Endpoint,
			peer.DB.StatDB(), peer.DB.RepairQueue(),
			peer.Overlay.Endpoint, peer.DB.Irreparable(),
			0, peer.Log.Named("checker"),
			config.Checker.Interval)

		// TODO: close segment repairer, currently this leaks connections
		segmentRepairer, err := config.Repairer.GetSegmentRepairer(context.TODO(), peer.Identity)
		if err != nil {
			return nil, errs.Combine(err, peer.Close())
		}

		peer.Repair.Repairer = repairer.NewService(peer.DB.RepairQueue(), segmentRepairer, config.Repairer.Interval, config.Repairer.MaxRepair)
	}

	{ // setup audit
		// TODO: audit needs many fixes
	}

	return peer, nil
}

func ignoreCancel(err error) error {
	if err == context.Canceled || err == grpc.ErrServerStopped {
		return nil
	}
	return err
}

// Run runs storage node until it's either closed or it errors.
func (peer *Peer) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var group errgroup.Group
	group.Go(func() error {
		return ignoreCancel(peer.Kademlia.Service.Bootstrap(ctx))
	})
	group.Go(func() error {
		return ignoreCancel(peer.Kademlia.Service.RunRefresh(ctx))
	})
	group.Go(func() error {
		return ignoreCancel(peer.Discovery.Service.Run(ctx))
	})
	group.Go(func() error {
		return ignoreCancel(peer.Repair.Checker.Run(ctx))
	})
	group.Go(func() error {
		return ignoreCancel(peer.Repair.Repairer.Run(ctx))
	})
	group.Go(func() error {
		return ignoreCancel(peer.Public.Server.Run(ctx))
	})

	return group.Wait()
}

// Close closes all the resources.
func (peer *Peer) Close() error {
	var errlist errs.Group

	// TODO: ensure that Close can be called on nil-s that way this code won't need the checks.

	// close services in reverse initialization order
	if peer.Repair.Repairer != nil {
		errlist.Add(peer.Repair.Repairer.Close())
	}
	if peer.Repair.Checker != nil {
		errlist.Add(peer.Repair.Checker.Close())
	}

	if peer.Agreements.Endpoint != nil {
		errlist.Add(peer.Agreements.Endpoint.Close())
	}

	if peer.Metainfo.Endpoint != nil {
		errlist.Add(peer.Metainfo.Endpoint.Close())
	}
	if peer.Metainfo.Database != nil {
		errlist.Add(peer.Metainfo.Database.Close())
	}

	if peer.Discovery.Service != nil {
		errlist.Add(peer.Discovery.Service.Close())
	}

	if peer.Overlay.Endpoint != nil {
		errlist.Add(peer.Overlay.Endpoint.Close())
	}
	if peer.Overlay.Service != nil {
		errlist.Add(peer.Overlay.Service.Close())
	}

	// TODO: add kademlia.Endpoint for consistency
	if peer.Kademlia.Service != nil {
		errlist.Add(peer.Kademlia.Service.Close())
	}
	if peer.Kademlia.RoutingTable != nil {
		errlist.Add(peer.Kademlia.RoutingTable.SelfClose())
	}

	// close servers
	if peer.Public.Server != nil {
		errlist.Add(peer.Public.Server.Close())
	} else {
		// peer.Public.Server automatically closes listener
		if peer.Public.Listener != nil {
			errlist.Add(peer.Public.Listener.Close())
		}
	}
	return errlist.Err()
}

// ID returns the peer ID.
func (peer *Peer) ID() storj.NodeID { return peer.Identity.ID }

// Local returns the peer local node info.
func (peer *Peer) Local() pb.Node { return peer.Kademlia.RoutingTable.Local() }

// Addr returns the public address.
func (peer *Peer) Addr() string { return peer.Public.Server.Addr().String() }
