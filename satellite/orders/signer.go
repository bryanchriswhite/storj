// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package orders

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"time"

	"github.com/zeebo/errs"

	"storj.io/common/pb"
	"storj.io/common/signing"
	"storj.io/common/storj"
)

// ErrSigner is default error class for Signer.
var ErrSigner = errs.Class("signer")

// Signer implements signing of order limits.
type Signer struct {
	// TODO: should this be a ref to the necessary pieces instead of the service?
	Service *Service

	// TODO: use a Template pb.OrderLimit here?

	RootPieceID storj.PieceID

	PieceExpiration time.Time
	OrderCreation   time.Time
	OrderExpiration time.Time

	PublicKey  storj.PiecePublicKey
	PrivateKey storj.PiecePrivateKey

	Serial storj.SerialNumber
	Action pb.PieceAction
	Limit  int64

	AddressedLimits []*pb.AddressedOrderLimit
}

// createSerial creates a timestamped serial number.
func createSerial(orderExpiration time.Time) (_ storj.SerialNumber, err error) {
	var serial storj.SerialNumber

	binary.BigEndian.PutUint64(serial[0:8], uint64(orderExpiration.Unix()))
	_, err = rand.Read(serial[8:])
	if err != nil {
		return storj.SerialNumber{}, ErrSigner.Wrap(err)
	}

	return serial, nil
}

// NewSigner creates an order limit signer.
func NewSigner(service *Service, rootPieceID storj.PieceID, pieceExpiration time.Time, orderCreation time.Time, limit int64, action pb.PieceAction) (*Signer, error) {
	signer := &Signer{}
	signer.Service = service

	signer.RootPieceID = rootPieceID

	signer.PieceExpiration = pieceExpiration
	signer.OrderCreation = orderCreation
	signer.OrderExpiration = orderCreation.Add(service.orderExpiration)

	var err error
	signer.PublicKey, signer.PrivateKey, err = storj.NewPieceKey()
	if err != nil {
		return nil, ErrSigner.Wrap(err)
	}

	signer.Serial, err = createSerial(signer.OrderExpiration)
	if err != nil {
		return nil, ErrSigner.Wrap(err)
	}

	signer.Action = action
	signer.Limit = limit

	return signer, nil
}

// NewSignerGet creates a new signer for get orders.
func NewSignerGet(service *Service, rootPieceID storj.PieceID, orderCreation time.Time, limit int64) (*Signer, error) {
	return NewSigner(service, rootPieceID, time.Time{}, orderCreation, limit, pb.PieceAction_GET)
}

// NewSignerPut creates a new signer for put orders.
func NewSignerPut(service *Service, pieceExpiration time.Time, orderCreation time.Time, limit int64) (*Signer, error) {
	rootPieceID := storj.NewPieceID()
	return NewSigner(service, rootPieceID, pieceExpiration, orderCreation, limit, pb.PieceAction_PUT)
}

// NewSignerDelete creates a new signer for delete orders.
func NewSignerDelete(service *Service, rootPieceID storj.PieceID, orderCreation time.Time) (*Signer, error) {
	return NewSigner(service, rootPieceID, time.Time{}, orderCreation, 0, pb.PieceAction_DELETE)
}

// NewSignerRepairGet creates a new signer for get repair orders.
func NewSignerRepairGet(service *Service, rootPieceID storj.PieceID, orderCreation time.Time, pieceSize int64) (*Signer, error) {
	return NewSigner(service, rootPieceID, time.Time{}, orderCreation, pieceSize, pb.PieceAction_GET_REPAIR)
}

// NewSignerRepairPut creates a new signer for put repair orders.
func NewSignerRepairPut(service *Service, rootPieceID storj.PieceID, pieceExpiration time.Time, orderCreation time.Time, pieceSize int64) (*Signer, error) {
	return NewSigner(service, rootPieceID, pieceExpiration, orderCreation, pieceSize, pb.PieceAction_PUT_REPAIR)
}

// NewSignerAudit creates a new signer for audit orders.
func NewSignerAudit(service *Service, rootPieceID storj.PieceID, orderCreation time.Time, pieceSize int64) (*Signer, error) {
	return NewSigner(service, rootPieceID, time.Time{}, orderCreation, pieceSize, pb.PieceAction_GET_AUDIT)
}

// NewSignerGracefulExit creates a new signer for graceful exit orders.
func NewSignerGracefulExit(service *Service, rootPieceID storj.PieceID, orderCreation time.Time, shareSize int32) (*Signer, error) {
	// TODO: we're using zero time.Time for piece expiration for some reason.

	// TODO: we're using `PUT_REPAIR` here even though we should be using `PUT`, such
	// that the storage node cannot distinguish between requests. We can't use `PUT`
	// because we don't want to charge bucket owners for graceful exit bandwidth, and
	// we can't use `PUT_GRACEFUL_EXIT` because storagenode will only accept upload
	// orders with `PUT` or `PUT_REPAIR` as the action. we also don't have a bunch of
	// supporting code/tables to aggregate `PUT_GRACEFUL_EXIT` bandwidth into our rollups
	// and stuff. so, for now, we just use `PUT_REPAIR` because it's the least bad of
	// our options. this should be fixed.
	return NewSigner(service, rootPieceID, time.Time{}, orderCreation, int64(shareSize), pb.PieceAction_PUT_REPAIR)
}

// Sign signs an order limit for the specified node.
func (signer *Signer) Sign(ctx context.Context, node storj.NodeURL, pieceNum int32) (_ *pb.AddressedOrderLimit, err error) {
	defer mon.Task()(&ctx)(&err)

	limit := &pb.OrderLimit{
		SerialNumber:    signer.Serial,
		SatelliteId:     signer.Service.satellite.ID(),
		UplinkPublicKey: signer.PublicKey,
		StorageNodeId:   node.ID,

		PieceId: signer.RootPieceID.Derive(node.ID, pieceNum),
		Limit:   signer.Limit,
		Action:  signer.Action,

		PieceExpiration: signer.PieceExpiration,
		OrderCreation:   signer.OrderCreation,
		OrderExpiration: signer.OrderExpiration,

		SatelliteAddress: signer.Service.satelliteAddress,
	}

	signedLimit, err := signing.SignOrderLimit(ctx, signer.Service.satellite, limit)
	if err != nil {
		return nil, ErrSigner.Wrap(err)
	}

	addressedLimit := &pb.AddressedOrderLimit{
		Limit: signedLimit,
		StorageNodeAddress: &pb.NodeAddress{
			Address: node.Address,
		},
	}

	signer.AddressedLimits = append(signer.AddressedLimits, addressedLimit)

	return addressedLimit, nil
}