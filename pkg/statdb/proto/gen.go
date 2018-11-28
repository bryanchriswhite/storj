// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package statdb

import (
	"storj.io/storj/pkg/pb"
	"storj.io/storj/pkg/storj"
)

// NodeID is an alias to storj.NodeID for use in generated protobuf code
type NodeID = storj.NodeID
// NodeIDLiiist is an alias to storj.NodeIDList for use in generated protobuf code
type NodeIDList = storj.NodeIDList
// Node is an alias to storj.Node for use in generated protobuf code
type Node = pb.Node
// NodeStats is an alias to storj.NodeStats for use in generated protobuf code
type NodeStats = pb.NodeStats

//go:generate protoc --gogo_out=plugins=grpc:. -I=. -I=$GOPATH/src/storj.io/storj/pkg/pb statdb.proto
