// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package overlay

import (
  "context"
  "flag"
  "fmt"
  "net"
  "net/http"

  "go.uber.org/zap"
  "google.golang.org/grpc"
  monkit "gopkg.in/spacemonkeygo/monkit.v2"

  "storj.io/storj/pkg/kademlia"
  "storj.io/storj/storage/redis"
  "storj.io/storj/pkg/utils"
  proto "storj.io/storj/protos/overlay"
)

var (
  redisAddress  string
  redisPassword string
  db            int
  tlsCertPath   string
  tlsKeyPath    string
  tlsHosts      string
  tlsCreate     bool
  tlsOverwrite  bool
  node          string
  bootstrapIP   string
  bootstrapPort string
  stun          bool
  httpPort      string
  gui           bool
  srvPort       uint
)

func init() {
  flag.StringVar(&redisAddress, "redisAddress", "", "The <IP:PORT> string to use for connection to a redis cache")
  flag.StringVar(&redisPassword, "redisPassword", "", "The password used for authentication to a secured redis instance")
  flag.IntVar(&db, "db", 0, "The network cache database")
  flag.StringVar(&tlsCertPath, "tlsCertPath", "", "TLS Certificate file")
  flag.StringVar(&tlsKeyPath, "tlsKeyPath", "", "TLS Key file")
  flag.StringVar(&tlsHosts, "tlsHosts", "", "TLS Key file")
  flag.BoolVar(&tlsCreate, "tlsCreate", false, "If true, generate a new TLS cert/key files")
  flag.BoolVar(&tlsOverwrite, "tlsOverwrite", false, "If true, overwrite existing TLS cert/key files")
  flag.StringVar(&httpPort, "httpPort", "", "The port for the health endpoint")
  flag.UintVar(&srvPort, "srvPort", 8080, "Port to listen on")
}

// NewServer creates a new Overlay Service Server
func NewServer() (*grpc.Server, error) {
  t := &utils.TLSFileOptions{
    CertRelPath: tlsCertPath,
    KeyRelPath:  tlsKeyPath,
    Create:      tlsCreate,
    Overwrite:   tlsOverwrite,
    Hosts:       tlsHosts,
  }

  creds, err := utils.NewServerTLSFromFile(t)
  if err != nil {
    return nil, err
  }

  credsOption := grpc.Creds(creds)
  grpcServer := grpc.NewServer(credsOption)
  proto.RegisterOverlayServer(grpcServer, &Overlay{})

  return grpcServer, nil
}

// NewClient connects to grpc server at the provided address with the provided options
// returns a new instance of an overlay Client
func NewClient(serverAddr *string, opts ...grpc.DialOption) (proto.OverlayClient, error) {
  t := &utils.TLSFileOptions{
    CertRelPath: tlsCertPath,
    Create:      tlsCreate,
    Overwrite:   tlsOverwrite,
    Hosts:       tlsHosts,
    Client:      true,
  }

  creds, err := utils.NewClientTLSFromFile(t)
  if err != nil {
    return nil, err
  }

  credsOption := grpc.WithTransportCredentials(creds)
  opts = append(opts, credsOption)
  conn, err := grpc.Dial(*serverAddr, opts...)
  if err != nil {
    return nil, err
  }

  return proto.NewOverlayClient(conn), nil
}

// Service contains all methods needed to implement the process.Service interface
type Service struct {
  logger  *zap.Logger
  metrics *monkit.Registry
}

// Process is the main function that executes the service
func (s *Service) Process(ctx context.Context) error {
	// TODO
	// 1. Boostrap a node on the network
	// 2. Start up the overlay gRPC service
	// 3. Connect to Redis
	// 4. Boostrap Redis Cache

	// TODO(coyle): Should add the ability to pass a configuration to change the bootstrap node
	in := kademlia.GetIntroNode()

	kad, err := kademlia.NewKademlia([]proto.Node{in}, "127.0.0.1", "8080")
	if err != nil {
		s.logger.Error("Failed to instantiate new Kademlia", zap.Error(err))
		return err
	}

	if err := kad.ListenAndServe(); err != nil {
		s.logger.Error("Failed to ListenAndServe on new Kademlia", zap.Error(err))
		return err
	}

	if err := kad.Bootstrap(ctx); err != nil {
		s.logger.Error("Failed to Bootstrap on new Kademlia", zap.Error(err))
		return err
	}

	// bootstrap cache
	cache, err := redis.NewOverlayClient(redisAddress, redisPassword, db, kad)
	if err != nil {
		s.logger.Error("Failed to create a new redis overlay client", zap.Error(err))
		return err
	}

	if err := cache.Bootstrap(ctx); err != nil {
		s.logger.Error("Failed to boostrap cache", zap.Error(err))
		return err
	}

	// send off cache refreshes concurrently
	go cache.Refresh(ctx)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", srvPort))
	if err != nil {
		s.logger.Error("Failed to initialize TCP connection", zap.Error(err))
		return err
	}

	grpcServer := grpc.NewServer()
	proto.RegisterOverlayServer(grpcServer, &Overlay{
		kad:     kad,
		DB:      cache,
		logger:  s.logger,
		metrics: s.metrics,
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintln(w, "OK") })
	go func() { http.ListenAndServe(fmt.Sprintf(":%s", httpPort), nil) }()

  defer grpcServer.GracefulStop()
  return grpcServer.Serve(lis)
}

// SetLogger adds the initialized logger to the Service
func (s *Service) SetLogger(l *zap.Logger) error {
  s.logger = l
  return nil
}

// SetMetricHandler adds the initialized metric handler to the Service
func (s *Service) SetMetricHandler(m *monkit.Registry) error {
  s.metrics = m
  return nil
}

// InstanceID implements Service.InstanceID
func (s *Service) InstanceID() string { return "" }
