// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package statdb

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"storj.io/storj/internal/storj"
	"storj.io/storj/pkg/storj"

	dbx "storj.io/storj/pkg/statdb/dbx"
	pb "storj.io/storj/pkg/statdb/proto"
)

var (
	ctx = context.Background()
)

func TestCreateDoesNotExist(t *testing.T) {
	dbPath := getDBPath()
	statdb, db, err := getServerAndDB(dbPath)
	assert.NoError(t, err)

	apiKey := []byte("")
	nodeID := teststorj.NodeIDFromString("testnodeid")
	node := &pb.Node{Id: nodeID}
	createReq := &pb.CreateRequest{
		Node:   node,
		APIKey: apiKey,
	}
	resp, err := statdb.Create(ctx, createReq)
	assert.NoError(t, err)
	stats := resp.Stats
	assert.EqualValues(t, 0, stats.AuditSuccessRatio)
	assert.EqualValues(t, 0, stats.UptimeRatio)

	nodeInfo, err := db.Get_Node_By_Id(ctx, dbx.Node_Id(nodeID.Bytes()))
	assert.NoError(t, err)

	assert.EqualValues(t, nodeID, nodeInfo.Id)
	assert.EqualValues(t, 0, nodeInfo.AuditSuccessRatio)
	assert.EqualValues(t, 0, nodeInfo.UptimeRatio)
}

func TestCreateExists(t *testing.T) {
	dbPath := getDBPath()
	statdb, db, err := getServerAndDB(dbPath)
	assert.NoError(t, err)

	apiKey := []byte("")
	nodeID := teststorj.NodeIDFromString("testnodeid")

	auditSuccessCount, totalAuditCount, auditRatio := getRatio(4, 10)
	uptimeSuccessCount, totalUptimeCount, uptimeRatio := getRatio(8, 25)
	err = createNode(ctx, db, nodeID, auditSuccessCount, totalAuditCount, auditRatio,
		uptimeSuccessCount, totalUptimeCount, uptimeRatio)
	assert.NoError(t, err)

	node := &pb.Node{Id: nodeID}
	createReq := &pb.CreateRequest{
		Node:   node,
		APIKey: apiKey,
	}
	_, err = statdb.Create(ctx, createReq)
	assert.Error(t, err)
}
func TestCreateWithStats(t *testing.T) {
	dbPath := getDBPath()
	statdb, db, err := getServerAndDB(dbPath)
	assert.NoError(t, err)

	auditSuccessCount, totalAuditCount, auditRatio := getRatio(4, 10)
	uptimeSuccessCount, totalUptimeCount, uptimeRatio := getRatio(8, 25)
	apiKey := []byte("")
	nodeID := teststorj.NodeIDFromString("testnodeid")
	node := &pb.Node{Id: nodeID}
	stats := &pb.NodeStats{
		AuditCount:         totalAuditCount,
		AuditSuccessCount:  auditSuccessCount,
		UptimeCount:        totalUptimeCount,
		UptimeSuccessCount: uptimeSuccessCount,
	}
	createReq := &pb.CreateRequest{
		Node:   node,
		Stats:  stats,
		APIKey: apiKey,
	}
	resp, err := statdb.Create(ctx, createReq)
	assert.NoError(t, err)
	s := resp.Stats
	assert.EqualValues(t, auditRatio, s.AuditSuccessRatio)
	assert.EqualValues(t, uptimeRatio, s.UptimeRatio)

	nodeInfo, err := db.Get_Node_By_Id(ctx, dbx.Node_Id(nodeID.Bytes()))
	assert.NoError(t, err)

	assert.EqualValues(t, nodeID, nodeInfo.Id, nodeID)
	assert.EqualValues(t, auditRatio, nodeInfo.AuditSuccessRatio)
	assert.EqualValues(t, uptimeRatio, nodeInfo.UptimeRatio)
}

func TestGetExists(t *testing.T) {
	dbPath := getDBPath()
	statdb, db, err := getServerAndDB(dbPath)
	assert.NoError(t, err)

	apiKey := []byte("")
	nodeID := teststorj.NodeIDFromString("testnodeid")

	auditSuccessCount, totalAuditCount, auditRatio := getRatio(4, 10)
	uptimeSuccessCount, totalUptimeCount, uptimeRatio := getRatio(8, 25)

	err = createNode(ctx, db, nodeID, auditSuccessCount, totalAuditCount, auditRatio,
		uptimeSuccessCount, totalUptimeCount, uptimeRatio)
	assert.NoError(t, err)

	getReq := &pb.GetRequest{
		NodeId: nodeID,
		APIKey: apiKey,
	}
	resp, err := statdb.Get(ctx, getReq)
	assert.NoError(t, err)

	stats := resp.Stats
	assert.EqualValues(t, auditRatio, stats.AuditSuccessRatio)
	assert.EqualValues(t, uptimeRatio, stats.UptimeRatio)
}

func TestGetDoesNotExist(t *testing.T) {
	dbPath := getDBPath()
	statdb, _, err := getServerAndDB(dbPath)
	assert.NoError(t, err)

	apiKey := []byte("")
	nodeID := teststorj.NodeIDFromString("testnodeid")

	getReq := &pb.GetRequest{
		NodeId: nodeID,
		APIKey: apiKey,
	}
	_, err = statdb.Get(ctx, getReq)
	assert.Error(t, err)
}

func TestFindValidNodes(t *testing.T) {
	NodeIDs := teststorj.NodeIDsFromStrings(
		"id1", "id2",
		"id3", "id4",
		"id5", "id7",
	)
	dbPath := getDBPath()
	statdb, db, err := getServerAndDB(dbPath)
	assert.NoError(t, err)

	apiKey := []byte("")

	for _, tt := range []struct {
		nodeID             storj.NodeID
		auditSuccessCount  int64
		totalAuditCount    int64
		auditRatio         float64
		uptimeSuccessCount int64
		totalUptimeCount   int64
		uptimeRatio        float64
	}{
		{NodeIDs[0], 10, 20, 0.5, 10, 20, 0.5},   // bad ratios
		{NodeIDs[1], 20, 20, 1, 20, 20, 1},       // good ratios
		{NodeIDs[2], 20, 20, 1, 10, 20, 0.5},     // good audit success bad uptime
		{NodeIDs[3], 10, 20, 0.5, 20, 20, 1},     // good uptime bad audit success
		{NodeIDs[4], 5, 5, 1, 5, 5, 1},           // good ratios not enough audits
		{NodeIDs[5], 20, 20, 1, 20, 20, 1},       // good ratios, excluded from query
		{NodeIDs[6], 19, 20, 0.95, 19, 20, 0.95}, // borderline ratios
	} {
		err = createNode(ctx, db, tt.nodeID, tt.auditSuccessCount, tt.totalAuditCount, tt.auditRatio,
			tt.uptimeSuccessCount, tt.totalUptimeCount, tt.uptimeRatio)
		assert.NoError(t, err)
	}

	findValidNodesReq := &pb.FindValidNodesRequest{
		NodeIds: storj.NodeIDList{
			NodeIDs[0], NodeIDs[1],
			NodeIDs[2], NodeIDs[3],
			NodeIDs[4], NodeIDs[6],
		},
		MinStats: &pb.NodeStats{
			AuditSuccessRatio: 0.95,
			UptimeRatio:       0.95,
			AuditCount:        15,
		},
		APIKey: apiKey,
	}

	resp, err := statdb.FindValidNodes(ctx, findValidNodesReq)
	assert.NoError(t, err)

	passed := resp.PassedIds

	assert.Contains(t, passed, NodeIDs[1])
	assert.Contains(t, passed, NodeIDs[6])
	assert.Len(t, passed, 2)
}

func TestUpdateExists(t *testing.T) {
	dbPath := getDBPath()
	statdb, db, err := getServerAndDB(dbPath)
	assert.NoError(t, err)

	apiKey := []byte("")
	nodeID := teststorj.NodeIDFromString("testnodeid")

	auditSuccessCount, totalAuditCount, auditRatio := getRatio(4, 10)
	uptimeSuccessCount, totalUptimeCount, uptimeRatio := getRatio(8, 25)
	err = createNode(ctx, db, nodeID, auditSuccessCount, totalAuditCount, auditRatio,
		uptimeSuccessCount, totalUptimeCount, uptimeRatio)
	assert.NoError(t, err)

	node := &pb.Node{
		Id:                 nodeID,
		UpdateAuditSuccess: true,
		AuditSuccess:       true,
		UpdateUptime:       true,
		IsUp:               false,
	}
	updateReq := &pb.UpdateRequest{
		Node:   node,
		APIKey: apiKey,
	}
	resp, err := statdb.Update(ctx, updateReq)
	assert.NoError(t, err)

	_, _, newAuditRatio := getRatio(int(auditSuccessCount+1), int(totalAuditCount+1))
	_, _, newUptimeRatio := getRatio(int(uptimeSuccessCount), int(totalUptimeCount+1))
	stats := resp.Stats
	assert.EqualValues(t, newAuditRatio, stats.AuditSuccessRatio)
	assert.EqualValues(t, newUptimeRatio, stats.UptimeRatio)
}

func TestUpdateBatchExists(t *testing.T) {
	dbPath := getDBPath()
	statdb, db, err := getServerAndDB(dbPath)
	assert.NoError(t, err)

	apiKey := []byte("")
	nodeID1 := teststorj.NodeIDFromString("testnodeid1")
	nodeID2 := teststorj.NodeIDFromString("testnodeid2")

	auditSuccessCount1, totalAuditCount1, auditRatio1 := getRatio(4, 10)
	uptimeSuccessCount1, totalUptimeCount1, uptimeRatio1 := getRatio(8, 25)
	err = createNode(ctx, db, nodeID1, auditSuccessCount1, totalAuditCount1, auditRatio1,
		uptimeSuccessCount1, totalUptimeCount1, uptimeRatio1)
	assert.NoError(t, err)
	auditSuccessCount2, totalAuditCount2, auditRatio2 := getRatio(7, 10)
	uptimeSuccessCount2, totalUptimeCount2, uptimeRatio2 := getRatio(8, 20)
	err = createNode(ctx, db, nodeID2, auditSuccessCount2, totalAuditCount2, auditRatio2,
		uptimeSuccessCount2, totalUptimeCount2, uptimeRatio2)
	assert.NoError(t, err)

	node1 := &pb.Node{
		Id:                 nodeID1,
		UpdateAuditSuccess: true,
		AuditSuccess:       true,
		UpdateUptime:       true,
		IsUp:               false,
	}
	node2 := &pb.Node{
		Id:                 nodeID2,
		UpdateAuditSuccess: true,
		AuditSuccess:       true,
		UpdateUptime:       false,
	}
	updateBatchReq := &pb.UpdateBatchRequest{
		NodeList: []*pb.Node{node1, node2},
		APIKey:   apiKey,
	}
	resp, err := statdb.UpdateBatch(ctx, updateBatchReq)
	assert.NoError(t, err)

	_, _, newAuditRatio1 := getRatio(int(auditSuccessCount1+1), int(totalAuditCount1+1))
	_, _, newUptimeRatio1 := getRatio(int(uptimeSuccessCount1), int(totalUptimeCount1+1))
	_, _, newAuditRatio2 := getRatio(int(auditSuccessCount2+1), int(totalAuditCount2+1))
	stats1 := resp.StatsList[0]
	stats2 := resp.StatsList[1]
	assert.EqualValues(t, newAuditRatio1, stats1.AuditSuccessRatio)
	assert.EqualValues(t, newUptimeRatio1, stats1.UptimeRatio)
	assert.EqualValues(t, newAuditRatio2, stats2.AuditSuccessRatio)
	assert.EqualValues(t, uptimeRatio2, stats2.UptimeRatio)
}

func TestUpdateBatchDoesNotExist(t *testing.T) {
	dbPath := getDBPath()
	statdb, db, err := getServerAndDB(dbPath)
	assert.NoError(t, err)

	apiKey := []byte("")
	nodeID1 := teststorj.NodeIDFromString("testnodeid1")
	nodeID2 := teststorj.NodeIDFromString("testnodeid2")

	auditSuccessCount1, totalAuditCount1, auditRatio1 := getRatio(4, 10)
	uptimeSuccessCount1, totalUptimeCount1, uptimeRatio1 := getRatio(8, 25)
	err = createNode(ctx, db, nodeID1, auditSuccessCount1, totalAuditCount1, auditRatio1,
		uptimeSuccessCount1, totalUptimeCount1, uptimeRatio1)
	assert.NoError(t, err)

	node1 := &pb.Node{
		Id:                 nodeID1,
		UpdateAuditSuccess: true,
		AuditSuccess:       true,
		UpdateUptime:       true,
		IsUp:               false,
	}
	node2 := &pb.Node{
		Id:                 nodeID2,
		UpdateAuditSuccess: true,
		AuditSuccess:       true,
		UpdateUptime:       false,
	}
	updateBatchReq := &pb.UpdateBatchRequest{
		NodeList: []*pb.Node{node1, node2},
		APIKey:   apiKey,
	}
	_, err = statdb.UpdateBatch(ctx, updateBatchReq)
	assert.NoError(t, err)
}

func TestUpdateBatchEmpty(t *testing.T) {
	dbPath := getDBPath()
	statdb, db, err := getServerAndDB(dbPath)
	assert.NoError(t, err)

	apiKey := []byte("")
	nodeID1 := teststorj.NodeIDFromString("testnodeid1")

	auditSuccessCount1, totalAuditCount1, auditRatio1 := getRatio(4, 10)
	uptimeSuccessCount1, totalUptimeCount1, uptimeRatio1 := getRatio(8, 25)
	err = createNode(ctx, db, nodeID1, auditSuccessCount1, totalAuditCount1, auditRatio1,
		uptimeSuccessCount1, totalUptimeCount1, uptimeRatio1)
	assert.NoError(t, err)

	updateBatchReq := &pb.UpdateBatchRequest{
		NodeList: []*pb.Node{},
		APIKey:   apiKey,
	}
	resp, err := statdb.UpdateBatch(ctx, updateBatchReq)
	assert.NoError(t, err)
	assert.Equal(t, len(resp.StatsList), 0)
}

func TestCreateEntryIfNotExists(t *testing.T) {
	dbPath := getDBPath()
	statdb, db, err := getServerAndDB(dbPath)
	assert.NoError(t, err)

	apiKey := []byte("")
	nodeID1 := teststorj.NodeIDFromString("testnodeid1")
	nodeID2 := teststorj.NodeIDFromString("testnodeid2")

	auditSuccessCount1, totalAuditCount1, auditRatio1 := getRatio(4, 10)
	uptimeSuccessCount1, totalUptimeCount1, uptimeRatio1 := getRatio(8, 25)
	err = createNode(ctx, db, nodeID1, auditSuccessCount1, totalAuditCount1, auditRatio1,
		uptimeSuccessCount1, totalUptimeCount1, uptimeRatio1)
	assert.NoError(t, err)

	node1 := &pb.Node{Id: nodeID1}
	createIfNotExistsReq1 := &pb.CreateEntryIfNotExistsRequest{
		Node:   node1,
		APIKey: apiKey,
	}
	_, err = statdb.CreateEntryIfNotExists(ctx, createIfNotExistsReq1)
	assert.NoError(t, err)

	nodeInfo1, err := db.Get_Node_By_Id(ctx, dbx.Node_Id(nodeID1.Bytes()))
	assert.NoError(t, err)
	assert.EqualValues(t, nodeID1, nodeInfo1.Id)
	assert.EqualValues(t, auditRatio1, nodeInfo1.AuditSuccessRatio)
	assert.EqualValues(t, uptimeRatio1, nodeInfo1.UptimeRatio)

	node2 := &pb.Node{Id: nodeID2}
	createIfNotExistsReq2 := &pb.CreateEntryIfNotExistsRequest{
		Node:   node2,
		APIKey: apiKey,
	}
	_, err = statdb.CreateEntryIfNotExists(ctx, createIfNotExistsReq2)
	assert.NoError(t, err)

	nodeInfo2, err := db.Get_Node_By_Id(ctx, dbx.Node_Id(nodeID2.Bytes()))
	assert.NoError(t, err)
	assert.EqualValues(t, nodeID2, nodeInfo2.Id)
	assert.EqualValues(t, 0, nodeInfo2.AuditSuccessRatio)
	assert.EqualValues(t, 0, nodeInfo2.UptimeRatio)
}

func getDBPath() string {
	return fmt.Sprintf("file:memdb%d?mode=memory&cache=shared", rand.Int63())
}

func getServerAndDB(path string) (statdb *Server, db *dbx.DB, err error) {
	statdb, err = NewServer("sqlite3", path, zap.NewNop())
	if err != nil {
		return &Server{}, &dbx.DB{}, err
	}
	db, err = dbx.Open("sqlite3", path)
	if err != nil {
		return &Server{}, &dbx.DB{}, err
	}
	return statdb, db, err
}

func createNode(ctx context.Context, db *dbx.DB, nodeID storj.NodeID,
	auditSuccessCount, totalAuditCount int64, auditRatio float64,
	uptimeSuccessCount, totalUptimeCount int64, uptimeRatio float64) error {
	_, err := db.Create_Node(
		ctx,
		dbx.Node_Id(nodeID.Bytes()),
		dbx.Node_AuditSuccessCount(auditSuccessCount),
		dbx.Node_TotalAuditCount(totalAuditCount),
		dbx.Node_AuditSuccessRatio(auditRatio),
		dbx.Node_UptimeSuccessCount(uptimeSuccessCount),
		dbx.Node_TotalUptimeCount(totalUptimeCount),
		dbx.Node_UptimeRatio(uptimeRatio),
	)
	return err
}

func getRatio(s, t int) (success, total int64, ratio float64) {
	ratio = float64(s) / float64(t)
	return int64(s), int64(t), ratio
}
