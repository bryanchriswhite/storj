// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package node

import (
	"context"
	"google.golang.org/grpc"

	// "storj.io/storj/internal/test"
	proto "storj.io/storj/protos/overlay"
)

// func NewNodeID(t *testing.T) string {
// 	// NewNodeID returns the string representation of a dht node PeerIdentity
// 	id, err := NewID(1, 38, 5)
// 	assert.NoError(t, err)
//
// 	return id.String()
// }

// func TestLookup(t *testing.T) {
// 	cases := []struct {
// 		self             proto.Node
// 		to               proto.Node
// 		find             proto.Node
// 		expectedErr      error
// 		expectedNumNodes int
// 	}{
// 		{
// 			self:        proto.Node{Id: NewNodeID(t), Address: &proto.NodeAddress{Address: ":7070"}},
// 			to:          proto.Node{Id: NewNodeID(t), Address: &proto.NodeAddress{Address: ":8080"}},
// 			find:        proto.Node{Id: NewNodeID(t), Address: &proto.NodeAddress{Address: ":9090"}},
// 			expectedErr: nil,
// 		},
// 	}
//
// 	// take writers
// 	certPath, keyPath := NewID(8)
//
// 	for _, v := range cases {
// 		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 8080))
// 		assert.NoError(t, err)
//
// 		srv, mock := newTestServer()
// 		go srv.Serve(lis)
// 		defer srv.Stop()
//
// 		// take readers
// 		nc, err := NewNodeClient(v.self, certPath, keyPath)
// 		assert.NoError(t, err)
//
// 		_, err = nc.Lookup(context.Background(), v.to, v.find)
// 		assert.Equal(t, v.expectedErr, err)
// 		assert.Equal(t, 1, mock.queryCalled)
// 	}
// }

func newTestServer() (*grpc.Server, *mockNodeServer) {
	grpcServer := grpc.NewServer()
	mn := &mockNodeServer{queryCalled: 0}

	proto.RegisterNodesServer(grpcServer, mn)

	return grpcServer, mn

}

type mockNodeServer struct {
	queryCalled int
}

func (mn *mockNodeServer) Query(ctx context.Context, req *proto.QueryRequest) (*proto.QueryResponse, error) {
	mn.queryCalled++
	return &proto.QueryResponse{}, nil
}
