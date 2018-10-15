// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package provider

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"storj.io/storj/pkg/peertls"
	"encoding/asn1"
)

func TestNewCA(t *testing.T) {
	expectedDifficulty := uint16(4)

	ca, err := NewCA(context.Background(), expectedDifficulty, 5)
	assert.NoError(t, err)
	assert.NotEmpty(t, ca)

	actualDifficulty := ca.ID.Difficulty()
	assert.True(t, actualDifficulty >= expectedDifficulty)
}

func TestFullCertificateAuthority_NewIdentity(t *testing.T) {
	check := newChecker(t)

	ca, err := NewCA(context.Background(), 12, 5)
	check.e(err).v(ca)
	fi, err := ca.NewIdentity()
	check.e(err).v(fi)

	assert.Equal(t, ca.Cert, fi.CA)
	assert.Equal(t, ca.ID, fi.ID)
	assert.NotEqual(t, ca.Key, fi.Key)
	assert.NotEqual(t, ca.Cert, fi.Leaf)

	err = fi.Leaf.CheckSignatureFrom(ca.Cert)
	assert.NoError(t, err)
}

func TestFullCertificateAuthority_RevokeIdentity(t *testing.T) {
	check := newChecker(t)
	ca, err := NewCA(context.Background(), 12, 5)
	check.e(err).v(ca)
	fi, err := ca.NewIdentity()
	check.e(err).v(fi)

	fi2, err := ca.RotateIdentity(fi)
	check.e(err).v(fi2)
	assert.NotEqual(t, fi.Leaf.SerialNumber, fi2.Leaf.SerialNumber)

	before := time.Now()
	r := revokedFromCert(fi2.Leaf)
	assert.NotNil(t, r)
	assert.NotEqual(t, fi2.Leaf.SerialNumber, r.SerialNumber)
	assert.Equal(t, fi.Leaf.SerialNumber, r.SerialNumber)
	assert.True(t, before.After(r.RevocationTime))
	assert.True(t, time.Now().After(r.RevocationTime))
}

func tempCA(t *testing.T) *FullCertificateAuthority {
	check := func(err error) {
		if !assert.NoError(t, err) {
			t.Fail()
		}
	}
	cert := `-----BEGIN CERTIFICATE-----
MIIBOTCB4KADAgECAhArQDR0b3E9H0kAenvqD9/2MAoGCCqGSM49BAMCMAAwIhgP
MDAwMTAxMDEwMDAwMDBaGA8wMDAxMDEwMTAwMDAwMFowADBZMBMGByqGSM49AgEG
CCqGSM49AwEHA0IABD9MFFY4lAb5pAj1iSYimiC+XjSRV4dgrumO8tQZ0WPaHbu/
uRwGgsLuGWWXZYy036a8+YSsBqi3Cej75yBVCHWjODA2MA4GA1UdDwEB/wQEAwIC
BDATBgNVHSUEDDAKBggrBgEFBQcDATAPBgNVHRMBAf8EBTADAQH/MAoGCCqGSM49
BAMCA0gAMEUCIQCmZfST3nQDW6VNzunx2KV6TRVDR8lUYjLD95+awT2O6wIgWXDF
eo9J61lRxNcjBGzHQ2B0SLbcGhDnkG+AjpF5t/8=
-----END CERTIFICATE-----`
	key := `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEILJD5CLOMQN+irJn8cZ/vvvwwDNvqIptvRIOJQTuUbm2oAoGCCqGSM49
AwEHoUQDQgAEP0wUVjiUBvmkCPWJJiKaIL5eNJFXh2Cu6Y7y1BnRY9odu7+5HAaC
wu4ZZZdljLTfprz5hKwGqLcJ6PvnIFUIdQ==
-----END EC PRIVATE KEY-----`

	var cb [][]byte
	cp, _ := pem.Decode([]byte(cert))
	cb = append(cb, cp.Bytes)

	c, err := ParseCertChain(cb)
	check(err)

	i, err := idFromKey(c[len(c)-1].PublicKey)
	check(err)

	kp, _ := pem.Decode([]byte(key))
	k, err := x509.ParseECPrivateKey(kp.Bytes)
	check(err)

	return &FullCertificateAuthority{
		Cert: c[0],
		Key:  k,
		ID:   i,
	}
}

func tempBoltRevocationDB(t *testing.T) (*BoltRevocationDB, func()) {
	tempDir := os.TempDir()
	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	dbPath := filepath.Join(tempDir, "revocations.db")
	dbConfig := RevocationDBConfig{
		Path: dbPath,
	}
	d, err := dbConfig.Open()
	if !assert.NoError(t, err) {
		t.Errorf("couldn't open revocation db")
	}

	b, ok := d.(*BoltRevocationDB)
	if !ok {
		t.Errorf("couldn't cast revocation db")
	}
	return b, cleanup
}

// func TestRevocationDB_Put(t *testing.T) {
// 	d, done := tempBoltRevocationDB(t)
// 	defer done()
// 	check := newChecker(t)
//
// 	ca := tempCA(t)
// 	fi, err := ca.NewIdentity()
// 	check.e(err).v(fi)
//
// 	r, err := ca.RotateIdentity(fi)
// 	rb, err := asn1.Marshal(r)
// 	check.e(err).v(rb)
//
// 	d.Put(fi, rb)
// 	rb, err := d.db.Get(fi.ID.Bytes())
// }

func TestRevocationDB_Get(t *testing.T) {
	d, done := tempBoltRevocationDB(t)
	defer done()
	check := newChecker(t)

	ca := tempCA(t)
	fi, err := ca.NewIdentity()
	check.e(err).v(fi)

	r, err := ca.RotateIdentity(fi)
	rb, err := asn1.Marshal(r)
	check.e(err).v(rb)

	d.db.Put(fi.ID.Bytes(), rb)
}

func TestCheckRevocations(t *testing.T) {
	d, done := tempBoltRevocationDB(t)
	defer done()
	check := newChecker(t)

	ca, err := NewCA(context.Background(), 12, 5)
	check.e(err).v(ca)
	fi, err := ca.NewIdentity()
	check.e(err).v(fi)

	fi2, err := ca.RotateIdentity(fi)
	check.e(err).v(fi2)

	err = peertls.VerifyPeerFunc(
		CheckRevocations(d),
	)([][]byte{fi2.Leaf.Raw, ca.Cert.Raw}, nil)
	assert.NoError(t, err)

	err = peertls.VerifyPeerFunc(
		CheckRevocations(d),
	)([][]byte{fi.Leaf.Raw, ca.Cert.Raw}, nil)
	assert.True(t, ErrRevoked.Has(err))
}

func NewCABenchmark(b *testing.B, difficulty uint16, concurrency uint) {
	for i := 0; i < b.N; i++ {
		NewCA(context.Background(), difficulty, concurrency)
	}
}

func BenchmarkNewCA_Difficulty8_Concurrency1(b *testing.B) {
	NewCABenchmark(b, 8, 1)
}

func BenchmarkNewCA_Difficulty8_Concurrency2(b *testing.B) {
	NewCABenchmark(b, 8, 2)
}

func BenchmarkNewCA_Difficulty8_Concurrency5(b *testing.B) {
	NewCABenchmark(b, 8, 5)
}

func BenchmarkNewCA_Difficulty8_Concurrency10(b *testing.B) {
	NewCABenchmark(b, 8, 10)
}

type checker struct {
	t *testing.T
}
func (c *checker) e(e error) *checker {
	if !assert.NoError(c.t, e) {
		c.t.FailNow()
	}
	return c
}
func (c *checker) v(v interface {}) *checker {
	if !assert.NotNil(c.t, v) {
		c.t.FailNow()
	}
	return c
}
func newChecker(t *testing.T) (*checker) {
	return &checker{
		t: t,
	}
}
