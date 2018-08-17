// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package overlay

import (
	"context"

	"google.golang.org/grpc"

	"storj.io/storj/pkg/dht"
	proto "storj.io/storj/protos/overlay"
)

// Client is the interface that defines an overlay client.
//
// Choose returns a list of storage NodeID's that fit the provided criteria.
// 	limit is the maximum number of nodes to be returned.
// 	space is the storage and bandwidth requested consumption in bytes.
//
// Lookup finds a Node with the provided identifier.
type Client interface {
	Choose(ctx context.Context, limit int, space int64) ([]*proto.Node, error)
	Lookup(ctx context.Context, nodeID dht.NodeID) (*proto.Node, error)
}

// Overlay is the overlay concrete implementation of the client interface
type Overlay struct {
	client proto.OverlayClient
}

// NewOverlayClient returns a new intialized Overlay Client
func NewOverlayClient(address string) (*Overlay, error) {
	c, err := NewClient(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &Overlay{
		client: c,
	}, nil
}

// a compiler trick to make sure *Overlay implements Client
var _ Client = (*Overlay)(nil)

// Choose implements the client.Choose interface
func (o *Overlay) Choose(ctx context.Context, amount int, space int64) ([]*proto.Node, error) {
	// TODO(coyle): We will also need to communicate with the reputation service here
	resp, err := o.client.FindStorageNodes(ctx, &proto.FindStorageNodesRequest{
		Opts: &proto.OverlayOptions{Amount: int64(amount), Restrictions: &proto.NodeRestrictions{
			FreeDisk: space,
		}},
	})
	if err != nil {
		return nil, Error.Wrap(err)
	}

	return resp.GetNodes(), nil
}

// Lookup provides a Node with the given address
func (o *Overlay) Lookup(ctx context.Context, nodeID dht.NodeID) (*proto.Node, error) {
	resp, err := o.client.Lookup(ctx, &proto.LookupRequest{NodeID: nodeID.String()})
	if err != nil {
		return nil, err
	}

	return resp.GetNode(), nil
}
