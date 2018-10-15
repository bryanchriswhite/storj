// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package provider

import (
	"context"
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net"
	"os"

	"github.com/zeebo/errs"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"

	"encoding/base64"
	"fmt"
	"math/bits"

	"crypto/x509/pkix"
	"encoding/asn1"
	"reflect"
	"storj.io/storj/pkg/peertls"
	"storj.io/storj/pkg/utils"
)

const (
	// IdentityLength is the number of bytes required to represent node id
	IdentityLength = uint16(256 / 8) // 256 bits
)

// PeerIdentity represents another peer on the network.
type PeerIdentity struct {
	// CA represents the peer's self-signed CA
	CA *x509.Certificate
	// Leaf represents the leaf they're currently using. The leaf should be
	// signed by the CA. The leaf is what is used for communication.
	Leaf *x509.Certificate
	// The ID taken from the CA public key
	ID nodeID
}

// FullIdentity represents you on the network. In addition to a PeerIdentity,
// a FullIdentity also has a Key, which a PeerIdentity doesn't have.
type FullIdentity struct {
	// CA represents the peer's self-signed CA. The ID is taken from this cert.
	CA *x509.Certificate
	// Leaf represents the leaf they're currently using. The leaf should be
	// signed by the CA. The leaf is what is used for communication.
	Leaf *x509.Certificate
	// The ID taken from the CA public key
	ID nodeID
	// Key is the key this identity uses with the leaf for communication.
	Key crypto.PrivateKey
}

// IdentitySetupConfig allows you to run a set of Responsibilities with the given
// identity. You can also just load an Identity from disk.
type IdentitySetupConfig struct {
	CertPath  string `help:"path to the certificate chain for this identity" default:"$CONFDIR/identity.cert"`
	KeyPath   string `help:"path to the private key for this identity" default:"$CONFDIR/identity.key"`
	Overwrite bool   `help:"if true, existing identity certs AND keys will overwritten for" default:"false"`
	Version   string `help:"semantic version of identity storage format" default:"0"`
}

// IdentityConfig allows you to run a set of Responsibilities with the given
// identity. You can also just load an Identity from disk.
type IdentityConfig struct {
	CertPath string `help:"path to the certificate chain for this identity" default:"$CONFDIR/identity.cert"`
	KeyPath  string `help:"path to the private key for this identity" default:"$CONFDIR/identity.key"`
	Address  string `help:"address to listen on" default:":7777"`
}

// FullIdentityFromPEM loads a FullIdentity from a certificate chain and
// private key file
func FullIdentityFromPEM(chainPEM, keyPEM []byte) (*FullIdentity, error) {
	cb, err := decodePEM(chainPEM)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	if len(cb) < 2 {
		return nil, errs.New("too few certificates in chain")
	}
	kb, err := decodePEM(keyPEM)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	// NB: there shouldn't be multiple keys in the key file but if there
	// are, this uses the first one
	k, err := x509.ParseECPrivateKey(kb[0])
	if err != nil {
		return nil, errs.New("unable to parse EC private key: %v", err)
	}
	ch, err := ParseCertChain(cb)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	i, err := idFromKey(ch[1].PublicKey)
	if err != nil {
		return nil, err
	}

	return &FullIdentity{
		CA:   ch[1],
		Leaf: ch[0],
		Key:  k,
		ID:   i,
	}, nil
}

// ParseCertChain converts a chain of certificate bytes into x509 certs
func ParseCertChain(chain [][]byte) ([]*x509.Certificate, error) {
	c := make([]*x509.Certificate, len(chain))
	for i, ct := range chain {
		cp, err := x509.ParseCertificate(ct)
		if err != nil {
			return nil, errs.Wrap(err)
		}
		c[i] = cp
	}
	return c, nil
}

// PeerIdentityFromCerts loads a PeerIdentity from a pair of leaf and ca x509 certificates
func PeerIdentityFromCerts(leaf, ca *x509.Certificate) (*PeerIdentity, error) {
	i, err := idFromKey(ca.PublicKey.(crypto.PublicKey))
	if err != nil {
		return nil, err
	}

	return &PeerIdentity{
		CA:   ca,
		ID:   i,
		Leaf: leaf,
	}, nil
}

// PeerIdentityFromContext loads a PeerIdentity from a ctx TLS credentials
func PeerIdentityFromContext(ctx context.Context) (*PeerIdentity, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, Error.New("unable to get grpc peer from contex")
	}
	tlsInfo := p.AuthInfo.(credentials.TLSInfo)
	c := tlsInfo.State.PeerCertificates
	if len(c) < 2 {
		return nil, Error.New("invalid certificate chain")
	}
	pi, err := PeerIdentityFromCerts(c[0], c[1])
	if err != nil {
		return nil, err
	}

	return pi, nil
}

// CheckRevocations compares the chain against the revocation DB,
// updating if a new, valid revocation is present
// NB: only checks the first parsed chain
func CheckRevocations(db RevocationDB) func([][]byte, [][]*x509.Certificate) error {
	return func(_ [][]byte, parsedChains [][]*x509.Certificate) error {
		err := checkRevocations(parsedChains[0], db)
		if !ErrRevoked.Has(err) {
		}

		return err
	}
}

func checkRevocations(chain []*x509.Certificate, db RevocationDB) error {
	cl := len(chain)
	if cl < 2 {
		return ErrCertificate.New("chain length expected to be >= 2; got: %d", cl)
	}

	pi, err := PeerIdentityFromCerts(chain[0], chain[1])
	if err != nil {
		return err
	}
	r, err := db.Get(pi)
	if err != nil {
		return err
	}
	if r != nil {
		return ErrRevoked.New("%s", r)
	}

	return nil
}

func revokedFromCert(c *x509.Certificate) *pkix.RevokedCertificate {
	var r pkix.RevokedCertificate
	for _, e := range c.Extensions {
		if reflect.DeepEqual(e.Id, peertls.OIDExample) {
			asn1.Unmarshal(e.Value, &r)
		}
	}
	return &r
}

// Stat returns the status of the identity cert/key files for the config
func (is IdentitySetupConfig) Stat() TLSFilesStatus {
	return statTLSFiles(is.CertPath, is.KeyPath)
}

// Create generates and saves a CA using the config
func (is IdentitySetupConfig) Create(ca *FullCertificateAuthority) (*FullIdentity, error) {
	fi, err := ca.NewIdentity()
	if err != nil {
		return nil, err
	}
	fi.CA = ca.Cert
	ic := IdentityConfig{
		CertPath: is.CertPath,
		KeyPath:  is.KeyPath,
	}
	return fi, ic.Save(fi)
}

// Load loads a FullIdentity from the config
func (ic IdentityConfig) Load() (*FullIdentity, error) {
	c, err := ioutil.ReadFile(ic.CertPath)
	if err != nil {
		return nil, peertls.ErrNotExist.Wrap(err)
	}
	k, err := ioutil.ReadFile(ic.KeyPath)
	if err != nil {
		return nil, peertls.ErrNotExist.Wrap(err)
	}

	fi, err := FullIdentityFromPEM(c, k)
	if err != nil {
		return nil, errs.New("failed to load identity %#v, %#v: %v",
			ic.CertPath, ic.KeyPath, err)
	}
	return fi, nil
}

// Save saves a FullIdentity according to the config
func (ic IdentityConfig) Save(fi *FullIdentity) error {
	f := os.O_WRONLY | os.O_CREATE
	c, err := openCert(ic.CertPath, f)
	if err != nil {
		return err
	}
	defer utils.LogClose(c)
	k, err := openKey(ic.KeyPath, f)
	if err != nil {
		return err
	}
	defer utils.LogClose(k)

	if err = peertls.WriteChain(c, fi.Leaf, fi.CA); err != nil {
		return err
	}
	if err = peertls.WriteKey(k, fi.Key); err != nil {
		return err
	}
	return nil
}

// Run will run the given responsibilities with the configured identity.
func (ic IdentityConfig) Run(ctx context.Context,
	responsibilities ...Responsibility) (
	err error) {
	defer mon.Task()(&ctx)(&err)

	pi, err := ic.Load()
	if err != nil {
		return err
	}

	lis, err := net.Listen("tcp", ic.Address)
	if err != nil {
		return err
	}
	defer func() { _ = lis.Close() }()

	s, err := NewProvider(pi, lis, responsibilities...)
	if err != nil {
		return err
	}
	defer func() { _ = s.Close() }()

	zap.S().Infof("Node %s started", s.Identity().ID)

	return s.Run(ctx)
}

// PeerCA returns a peer certificate authority based on the full identity
func (fi *FullIdentity) PeerCA() *PeerCertificateAuthority {
	return &PeerCertificateAuthority{
		Cert: fi.CA,
		ID: fi.ID,
	}
}

// TODO: comment
func (fi *FullIdentity) Revoked() *pkix.RevokedCertificate {
	return revokedFromCert(fi.Leaf)
}

// ServerOption returns a grpc `ServerOption` for incoming connections
// to the node with this full identity
func (fi *FullIdentity) ServerOption(db RevocationDB) (grpc.ServerOption, error) {
	ch := [][]byte{fi.Leaf.Raw, fi.CA.Raw}
	c, err := peertls.TLSCert(ch, fi.Leaf, fi.Key)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{*c},
		InsecureSkipVerify: true,
		ClientAuth:         tls.RequireAnyClientCert,
		VerifyPeerCertificate: peertls.VerifyPeerFunc(
			peertls.VerifyPeerCertChains,
			// TODO
			// CheckRevocations(db),
		),
	}

	return grpc.Creds(credentials.NewTLS(tlsConfig)), nil
}

// DialOption returns a grpc `DialOption` for making outgoing connections
// to the node with this peer identity
func (fi *FullIdentity) DialOption() (grpc.DialOption, error) {
	ch := [][]byte{fi.Leaf.Raw, fi.CA.Raw}
	c, err := peertls.TLSCert(ch, fi.Leaf, fi.Key)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{*c},
		InsecureSkipVerify: true,
		VerifyPeerCertificate: peertls.VerifyPeerFunc(
			peertls.VerifyPeerCertChains,
		),
	}

	return grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)), nil
}

type nodeID string

func (n nodeID) String() string { return string(n) }
func (n nodeID) Bytes() []byte  { return []byte(n) }
func (n nodeID) Difficulty() uint16 {
	hash, err := base64.URLEncoding.DecodeString(n.String())
	if err != nil {
		zap.S().Error(errs.Wrap(err))
	}

	for i := 1; i < len(hash); i++ {
		b := hash[len(hash)-i]

		if b != 0 {
			zeroBits := bits.TrailingZeros16(uint16(b))
			if zeroBits == 16 {
				zeroBits = 0
			}

			return uint16((i-1)*8 + zeroBits)
		}
	}

	// NB: this should never happen
	reason := fmt.Sprintf("difficulty matches hash length! hash: %s", hash)
	zap.S().Error(reason)
	panic(reason)
}
