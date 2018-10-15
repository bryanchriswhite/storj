// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package peertls

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"io"
	"time"

	"github.com/zeebo/errs"
)

const (
	// BlockTypeEcPrivateKey is the value to define a block type of private key
	BlockTypeEcPrivateKey = "EC PRIVATE KEY"
	// BlockTypeCertificate is the value to define a block type of certificate
	BlockTypeCertificate = "CERTIFICATE"
	// BlockTypeIDOptions is the value to define a block type of id options
	// (e.g. `version`)
	BlockTypeIDOptions = "ID OPTIONS"
)

var (
	// ErrNotExist is used when a file or directory doesn't exist
	ErrNotExist = errs.Class("file or directory not found error")
	// ErrGenerate is used when an error occurred during cert/key generation
	ErrGenerate = errs.Class("tls generation error")
	// ErrUnsupportedKey is used when key type is not supported
	ErrUnsupportedKey = errs.Class("unsupported key type")
	// ErrTLSTemplate is used when an error occurs during tls template generation
	ErrTLSTemplate = errs.Class("tls template error")
	// ErrVerifyPeerCert is used when an error occurs during `VerifyPeerCertificate`
	ErrVerifyPeerCert = errs.Class("tls peer certificate verification error")
	// ErrVerifySignature is used when a cert-chain signature verificaion error occurs
	ErrVerifySignature = errs.Class("tls certificate signature verification error")

	// NB: see http://oid-info.com/get/2.999
	OIDExample = asn1.ObjectIdentifier{2, 999}
)

// PeerCertVerificationFunc is the signature for a `*tls.Config{}`'s
// `VerifyPeerCertificate` function.
type PeerCertVerificationFunc func([][]byte, [][]*x509.Certificate) error

// NewKey returns a new PrivateKey
func NewKey() (crypto.PrivateKey, error) {
	k, err := ecdsa.GenerateKey(authECCurve, rand.Reader)
	if err != nil {
		return nil, ErrGenerate.New("failed to generate private key: %v", err)
	}

	return k, nil
}

// NewCert returns a new x509 certificate using the provided templates and
// signed by the `signer` key
func NewCert(template, parentTemplate *x509.Certificate, pubKey crypto.PublicKey, signer crypto.PrivateKey) (*x509.Certificate, error) {
	k, ok := signer.(*ecdsa.PrivateKey)
	if !ok {
		return nil, ErrUnsupportedKey.New("%T", k)
	}

	if parentTemplate == nil {
		parentTemplate = template
	}

	cb, err := x509.CreateCertificate(
		rand.Reader,
		template,
		parentTemplate,
		pubKey,
		k,
	)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	c, err := x509.ParseCertificate(cb)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return c, nil
}

// TODO: comment
func RevokeCert(c *x509.Certificate, signer crypto.PrivateKey) (*pkix.CertificateList, error) {
	r := []pkix.RevokedCertificate{
		{
			SerialNumber:   c.SerialNumber,
			RevocationTime: time.Now(),
		},
	}
	// v, err := asn1.Marshal(r)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// return &pkix.Extension{
	// 	Id:       OIDExample,
	// 	Critical: true,
	// 	Value:    v,
	// }, nil
	t := time.Now()
	crlBytes, err := c.CreateCRL(rand.Reader, signer, r, t, t)
	if err != nil {
		errs.Wrap(err)
	}
	crl, err := x509.ParseCRL(crlBytes)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return crl, nil
}

// VerifyPeerFunc combines multiple `*tls.Config#VerifyPeerCertificate`
// functions and adds certificate parsing.
func VerifyPeerFunc(next ...PeerCertVerificationFunc) PeerCertVerificationFunc {
	return func(chain [][]byte, _ [][]*x509.Certificate) error {
		c, err := parseCertificateChains(chain)
		if err != nil {
			return err
		}

		for _, n := range next {
			if n != nil {
				if err := n(chain, [][]*x509.Certificate{c}); err != nil {
					return err
				}
			}
		}
		return nil
	}
}

// VerifyPeerCertChains verifies that the first certificate chain contains certificates
// which are signed by their respective parents, ending with a self-signed root
func VerifyPeerCertChains(_ [][]byte, parsedChains [][]*x509.Certificate) error {
	return verifyChainSignatures(parsedChains[0])
}

// NewKeyBlock converts an ASN1/DER-encoded byte-slice of a private key into
// a `pem.Block` pointer
func NewKeyBlock(b []byte) *pem.Block {
	return &pem.Block{Type: BlockTypeEcPrivateKey, Bytes: b}
}

// NewCertBlock converts an ASN1/DER-encoded byte-slice of a tls certificate
// into a `pem.Block` pointer
func NewCertBlock(b []byte) *pem.Block {
	return &pem.Block{Type: BlockTypeCertificate, Bytes: b}
}

// TLSCert creates a tls.Certificate from chains, key and leaf
func TLSCert(chain [][]byte, leaf *x509.Certificate, key crypto.PrivateKey) (*tls.Certificate, error) {
	var err error
	if leaf == nil {
		leaf, err = x509.ParseCertificate(chain[0])
		if err != nil {
			return nil, err
		}
	}

	return &tls.Certificate{
		Leaf:        leaf,
		Certificate: chain,
		PrivateKey:  key,
	}, nil
}

// WriteChain writes the certificate chain (leaf-first) to the writer, PEM-encoded.
func WriteChain(w io.Writer, chain ...*x509.Certificate) error {
	if len(chain) < 1 {
		return errs.New("expected at least one certificate for writing")
	}

	for _, c := range chain {
		if err := pem.Encode(w, NewCertBlock(c.Raw)); err != nil {
			return errs.Wrap(err)
		}
	}
	return nil
}

// WriteKey writes the private key to the writer, PEM-encoded.
func WriteKey(w io.Writer, key crypto.PrivateKey) error {
	var (
		kb  []byte
		err error
	)

	switch k := key.(type) {
	case *ecdsa.PrivateKey:
		kb, err = x509.MarshalECPrivateKey(k)
		if err != nil {
			return errs.Wrap(err)
		}
	default:
		return ErrUnsupportedKey.New("%T", k)
	}

	if err := pem.Encode(w, NewKeyBlock(kb)); err != nil {
		return errs.Wrap(err)
	}
	return nil
}
