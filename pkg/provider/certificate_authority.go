// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package provider

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"reflect"

	"github.com/zeebo/errs"

	"crypto/x509/pkix"
	"encoding/asn1"
	"go.uber.org/zap"
	"os"
	"storj.io/storj/pkg/peertls"
	"storj.io/storj/pkg/utils"
	"storj.io/storj/storage"
	"storj.io/storj/storage/boltdb"
)

// PeerCertificateAuthority represents the CA which is used to validate peer identities
type PeerCertificateAuthority struct {
	// Cert is the x509 certificate of the CA
	Cert *x509.Certificate
	// The ID is calculated from the CA public key.
	ID nodeID
}

// FullCertificateAuthority represents the CA which is used to author and validate full identities
type FullCertificateAuthority struct {
	// Cert is the x509 certificate of the CA
	Cert *x509.Certificate
	// The ID is calculated from the CA public key.
	ID nodeID
	// Key is the private key of the CA
	Key crypto.PrivateKey
}

// RevocationDB is a key/value store for keeping track of certificate revocations
type RevocationDB interface {
	Get(*PeerIdentity) (*pkix.CertificateList, error)
}
type BoltRevocationDB struct {
	db storage.KeyValueStore
}

// CASetupConfig is for creating a CA
type CASetupConfig struct {
	CertPath    string `help:"path to the certificate chain for this identity" default:"$CONFDIR/ca.cert"`
	KeyPath     string `help:"path to the private key for this identity" default:"$CONFDIR/ca.key"`
	Difficulty  uint64 `help:"minimum difficulty for identity generation" default:"12"`
	Timeout     string `help:"timeout for CA generation; golang duration string (0 no timeout)" default:"5m"`
	Overwrite   bool   `help:"if true, existing CA certs AND keys will overwritten" default:"false"`
	Concurrency uint   `help:"number of concurrent workers for certificate authority generation" default:"4"`
}

// PeerCAConfig is for locating a CA certificate without a private key
type PeerCAConfig struct {
	CertPath string `help:"path to the certificate chain for this identity" default:"$CONFDIR/ca.cert"`
}

// FullCAConfig is for locating a CA certificate and it's private key
type FullCAConfig struct {
	CertPath     string `help:"path to the certificate chain for this identity" default:"$CONFDIR/ca.cert"`
	KeyPath      string `help:"path to the private key for this identity" default:"$CONFDIR/ca.key"`
	RevocationDB RevocationDBConfig
}

type RevocationDBConfig struct {
	Path    string `help:"path to the certificate revocation database" default:"$CONFDIR/revocations.db"`
	MaxSize string `help:"maximum size before overwriting old revocations(\"0\" is unlimited)" default:"100M"`
}

// NewCA creates a new full identity with the given difficulty
func NewCA(ctx context.Context, difficulty uint16, concurrency uint) (*FullCertificateAuthority, error) {
	if concurrency < 1 {
		concurrency = 1
	}
	ctx, cancel := context.WithCancel(ctx)

	eC := make(chan error)
	caC := make(chan FullCertificateAuthority, 1)
	for i := 0; i < int(concurrency); i++ {
		go newCAWorker(ctx, difficulty, caC, eC)
	}

	select {
	case ca := <-caC:
		cancel()
		return &ca, nil
	case err := <-eC:
		cancel()
		return nil, err
	}
}

// Open creates or loads the revocation database described by the config
func (dc RevocationDBConfig) Open() (RevocationDB, error) {
	// TODO: newChecker db size against `dc.MaxSize? Warn?
	b, err := boltdb.NewClient(zap.L(), dc.Path, "revocations")
	if err != nil {
		return nil, err
	}

	return &BoltRevocationDB{
		db: b,
	}, nil
}

// TODO: comment
func (d BoltRevocationDB) Get(pi *PeerIdentity) (r *pkix.CertificateList, err error) {
	v, err := d.db.Get(pi.ID.Bytes())
	if err != nil {
		return nil, err
	}
	if v.IsZero() {
		return nil, nil
	}

	_, err = asn1.Unmarshal(v, r)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	return r, nil
}

// TODO: comment
func (d BoltRevocationDB) Put(pi *PeerIdentity, r *pkix.CertificateList) (error) {
	v, err := asn1.Marshal(r)
	if err != nil {
		return err
	}
	return d.db.Put(pi.ID.Bytes(), v)
}

// Status returns the status of the CA cert/key files for the config
func (caS CASetupConfig) Status() TLSFilesStatus {
	return statTLSFiles(caS.CertPath, caS.KeyPath)
}

// Create generates and saves a CA using the config
func (caS CASetupConfig) Create(ctx context.Context) (*FullCertificateAuthority, error) {
	ca, err := NewCA(ctx, uint16(caS.Difficulty), caS.Concurrency)
	if err != nil {
		return nil, err
	}
	caC := FullCAConfig{
		CertPath: caS.CertPath,
		KeyPath:  caS.KeyPath,
	}
	return ca, caC.Save(ca)
}

// Load loads a CA from the given configuration
func (fc FullCAConfig) Load() (*FullCertificateAuthority, error) {
	p, err := fc.PeerConfig().Load()
	if err != nil {
		return nil, err
	}

	kb, err := ioutil.ReadFile(fc.KeyPath)
	if err != nil {
		return nil, peertls.ErrNotExist.Wrap(err)
	}
	kp, _ := pem.Decode(kb)
	k, err := x509.ParseECPrivateKey(kp.Bytes)
	if err != nil {
		return nil, errs.New("unable to parse EC private key", err)
	}

	ec, ok := p.Cert.PublicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, peertls.ErrUnsupportedKey.New("certificate public key type not supported: %T", k)
	}

	if !reflect.DeepEqual(k.PublicKey, *ec) {
		return nil, errs.New("certificate public key and loaded")
	}

	return &FullCertificateAuthority{
		Cert: p.Cert,
		Key:  k,
		ID:   p.ID,
	}, nil
}

// Save saves a CA with the given configuration
func (fc FullCAConfig) Save(ca *FullCertificateAuthority) error {
	f := os.O_WRONLY | os.O_CREATE
	c, err := openCert(fc.CertPath, f)
	if err != nil {
		return err
	}
	defer utils.LogClose(c)
	k, err := openKey(fc.KeyPath, f)
	if err != nil {
		return err
	}
	defer utils.LogClose(k)

	if err = peertls.WriteChain(c, ca.Cert); err != nil {
		return err
	}
	if err = peertls.WriteKey(k, ca.Key); err != nil {
		return err
	}
	return nil
}

// PeerConfig converts a full ca config to a peer ca config
func (fc FullCAConfig) PeerConfig() PeerCAConfig {
	return PeerCAConfig{
		CertPath: fc.CertPath,
	}
}

// Load loads a CA from the given configuration
func (pc PeerCAConfig) Load() (*PeerCertificateAuthority, error) {
	cd, err := ioutil.ReadFile(pc.CertPath)
	if err != nil {
		return nil, peertls.ErrNotExist.Wrap(err)
	}

	var cb [][]byte
	for {
		var cp *pem.Block
		cp, cd = pem.Decode(cd)
		if cp == nil {
			break
		}
		cb = append(cb, cp.Bytes)
	}
	c, err := ParseCertChain(cb)
	if err != nil {
		return nil, errs.New("failed to load identity %#v: %v",
			pc.CertPath, err)
	}

	i, err := idFromKey(c[len(c)-1].PublicKey)
	if err != nil {
		return nil, err
	}

	return &PeerCertificateAuthority{
		Cert: c[0],
		ID:   i,
	}, nil
}

// NewIdentity generates a new `FullIdentity` based on the CA. The CA
// cert is included in the identity's cert chain and the identity's leaf cert
// is signed by the CA.
func (ca FullCertificateAuthority) NewIdentity() (*FullIdentity, error) {
	lT, err := peertls.LeafTemplate()
	if err != nil {
		return nil, err
	}
	k, err := peertls.NewKey()
	if err != nil {
		return nil, err
	}
	pk, ok := k.(*ecdsa.PrivateKey)
	if !ok {
		return nil, peertls.ErrUnsupportedKey.New("%T", k)
	}
	l, err := peertls.NewCert(&lT, ca.Cert, &pk.PublicKey, ca.Key)
	if err != nil {
		return nil, err
	}

	return &FullIdentity{
		CA:   ca.Cert,
		Leaf: l,
		Key:  k,
		ID:   ca.ID,
	}, nil
}

// TODO: comment
func (ca FullCertificateAuthority) RevokeIdentity(ri *FullIdentity) (*pkix.CertificateList, error) {
	r, err := peertls.RevokeCert(ri.Leaf, ca.Key)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// TODO: comment
func (ca FullCertificateAuthority) RotateIdentity(ri *FullIdentity) (*FullIdentity, error) {
	r, err := ca.RevokeIdentity(ri)
	if err != nil {
		return nil, err
	}
	v, err := asn1.Marshal(*r)
	if err != nil {
		return nil, err
	}
	e := pkix.Extension{
		Id:       peertls.OIDExample,
		Critical: true,
		Value:    v,
	}

	lT, err := peertls.LeafTemplate()
	if err != nil {
		return nil, err
	}
	k, err := peertls.NewKey()
	if err != nil {
		return nil, err
	}
	pk, ok := k.(*ecdsa.PrivateKey)
	if !ok {
		return nil, peertls.ErrUnsupportedKey.New("%T", k)
	}

	lT.ExtraExtensions = append(lT.ExtraExtensions, e)
	l, err := peertls.NewCert(&lT, ca.Cert, &pk.PublicKey, ca.Key)
	if err != nil {
		return nil, err
	}

	fi := &FullIdentity{
		CA:   ca.Cert,
		Leaf: l,
		Key:  k,
		ID:   ca.ID,
	}

	return fi, nil
}
