// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package provider

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	"storj.io/storj/pkg/peertls"
)

func TestPeerIdentityFromCertChain(t *testing.T) {
	k, err := peertls.NewKey()
	assert.NoError(t, err)

	caT, err := peertls.CATemplate()
	assert.NoError(t, err)

	cp, _ := k.(*ecdsa.PrivateKey)
	c, err := peertls.NewCert(&caT, nil, &cp.PublicKey, k)
	assert.NoError(t, err)

	lT, err := peertls.LeafTemplate()
	assert.NoError(t, err)

	lk, err := peertls.NewKey()
	assert.NoError(t, err)

	lp, _ := lk.(*ecdsa.PrivateKey)
	l, err := peertls.NewCert(&lT, &caT, &lp.PublicKey, k)
	assert.NoError(t, err)

	pi, err := PeerIdentityFromCerts(l, c)
	assert.NoError(t, err)
	assert.Equal(t, c, pi.CA)
	assert.Equal(t, l, pi.Leaf)
	assert.NotEmpty(t, pi.ID)
}

func TestFullIdentityFromPEM(t *testing.T) {
	ck, err := peertls.NewKey()
	assert.NoError(t, err)

	caT, err := peertls.CATemplate()
	assert.NoError(t, err)

	cp, _ := ck.(*ecdsa.PrivateKey)
	c, err := peertls.NewCert(&caT, nil, &cp.PublicKey, ck)
	assert.NoError(t, err)
	assert.NoError(t, err)
	assert.NotEmpty(t, c)

	lT, err := peertls.LeafTemplate()
	assert.NoError(t, err)

	lk, err := peertls.NewKey()
	assert.NoError(t, err)

	lp, _ := lk.(*ecdsa.PrivateKey)
	l, err := peertls.NewCert(&lT, &caT, &lp.PublicKey, ck)
	assert.NoError(t, err)
	assert.NotEmpty(t, l)

	chainPEM := bytes.NewBuffer([]byte{})
	pem.Encode(chainPEM, peertls.NewCertBlock(l.Raw))
	pem.Encode(chainPEM, peertls.NewCertBlock(c.Raw))

	lkE, ok := lk.(*ecdsa.PrivateKey)
	assert.True(t, ok)
	assert.NotEmpty(t, lkE)

	lkB, err := x509.MarshalECPrivateKey(lkE)
	assert.NoError(t, err)
	assert.NotEmpty(t, lkB)

	keyPEM := bytes.NewBuffer([]byte{})
	pem.Encode(keyPEM, peertls.NewKeyBlock(lkB))

	fi, err := FullIdentityFromPEM(chainPEM.Bytes(), keyPEM.Bytes())
	assert.NoError(t, err)
	assert.Equal(t, l.Raw, fi.Leaf.Raw)
	assert.Equal(t, c.Raw, fi.CA.Raw)
	assert.Equal(t, lk, fi.Key)
}

func TestIdentityConfig_SaveIdentity(t *testing.T) {
	ic, fi, _, done := tempIdentity(t)
	defer done()

	chainPEM := bytes.NewBuffer([]byte{})
	pem.Encode(chainPEM, peertls.NewCertBlock(fi.Leaf.Raw))
	pem.Encode(chainPEM, peertls.NewCertBlock(fi.CA.Raw))

	privateKey, ok := fi.Key.(*ecdsa.PrivateKey)
	assert.True(t, ok)
	assert.NotEmpty(t, privateKey)

	keyBytes, err := x509.MarshalECPrivateKey(privateKey)
	assert.NoError(t, err)
	assert.NotEmpty(t, keyBytes)

	keyPEM := bytes.NewBuffer([]byte{})
	pem.Encode(keyPEM, peertls.NewKeyBlock(keyBytes))

	err = ic.Save(fi)
	assert.NoError(t, err)

	if runtime.GOOS != "windows" {
		//TODO (windows): ignoring for windows due to different default permissions
		certInfo, err := os.Stat(ic.CertPath)
		assert.NoError(t, err)
		assert.Equal(t, os.FileMode(0644), certInfo.Mode())

		keyInfo, err := os.Stat(ic.KeyPath)
		assert.NoError(t, err)
		assert.Equal(t, os.FileMode(0600), keyInfo.Mode())
	}

	savedChainPEM, err := ioutil.ReadFile(ic.CertPath)
	assert.NoError(t, err)

	savedKeyPEM, err := ioutil.ReadFile(ic.KeyPath)
	assert.NoError(t, err)

	assert.Equal(t, chainPEM.Bytes(), savedChainPEM)
	assert.Equal(t, keyPEM.Bytes(), savedKeyPEM)
}

func tempIdentityConfig() (*IdentityConfig, error, func()) {
	tmpDir, err := ioutil.TempDir("", "tempIdentity")
	if err != nil {
		return nil, err, nil
	}

	cleanup := func() { os.RemoveAll(tmpDir) }

	return &IdentityConfig{
		CertPath: filepath.Join(tmpDir, "chain.pem"),
		KeyPath:  filepath.Join(tmpDir, "key.pem"),
	}, nil, cleanup
}

func tempIdentity(t *testing.T) (*IdentityConfig, *FullIdentity, uint16, func()) {
	// NB: known difficulty
	difficulty := uint16(12)

	chain := `-----BEGIN CERTIFICATE-----
MIIBQDCB56ADAgECAhB+u3d03qyW/ROgwy/ZsPccMAoGCCqGSM49BAMCMAAwIhgP
MDAwMTAxMDEwMDAwMDBaGA8wMDAxMDEwMTAwMDAwMFowADBZMBMGByqGSM49AgEG
CCqGSM49AwEHA0IABIZrEPV/ExEkF0qUF0fJ3qSeGt5oFUX231v02NSUywcQ/Ve0
v3nHbmcJdjWBis2AkfL25mYDVC25jLl4tylMKumjPzA9MA4GA1UdDwEB/wQEAwIF
oDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwDAYDVR0TAQH/BAIwADAK
BggqhkjOPQQDAgNIADBFAiEA2ZvsR0ncw4mHRIg2Isavd+XVEoMo/etXQRAkDy9n
wyoCIDykUsqjshc9kCrXOvPSN8GuO2bNoLu5C7K1GlE/HI2X
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIBODCB4KADAgECAhAOcvhKe5TWT44LqFfgA1f8MAoGCCqGSM49BAMCMAAwIhgP
MDAwMTAxMDEwMDAwMDBaGA8wMDAxMDEwMTAwMDAwMFowADBZMBMGByqGSM49AgEG
CCqGSM49AwEHA0IABIZrEPV/ExEkF0qUF0fJ3qSeGt5oFUX231v02NSUywcQ/Ve0
v3nHbmcJdjWBis2AkfL25mYDVC25jLl4tylMKumjODA2MA4GA1UdDwEB/wQEAwIC
BDATBgNVHSUEDDAKBggrBgEFBQcDATAPBgNVHRMBAf8EBTADAQH/MAoGCCqGSM49
BAMCA0cAMEQCIGAZfPT1qvlnkTacojTtP20ZWf6XbnSztJHIKlUw6AE+AiB5Vcjj
awRaC5l1KBPGqiKB0coVXDwhW+K70l326MPUcg==
-----END CERTIFICATE-----`

	key := `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIKGjEetrxKrzl+AL1E5LXke+1ElyAdjAmr88/1Kx09+doAoGCCqGSM49
AwEHoUQDQgAEoLy/0hs5deTXZunRumsMkiHpF0g8wAc58aXANmr7Mxx9tzoIYFnx
0YN4VDKdCtUJa29yA6TIz1MiIDUAcB5YCA==
-----END EC PRIVATE KEY-----`

	ic, err, cleanup := tempIdentityConfig()
	assert.NoError(t, err)

	fi, err := FullIdentityFromPEM([]byte(chain), []byte(key))
	assert.NoError(t, err)

	return ic, fi, difficulty, cleanup
}

func TestIdentityConfig_LoadIdentity(t *testing.T) {
	ic, expectedFI, _, done := tempIdentity(t)
	defer done()

	err := ic.Save(expectedFI)
	assert.NoError(t, err)

	fi, err := ic.Load()
	assert.NoError(t, err)
	assert.NotEmpty(t, fi)
	assert.NotEmpty(t, fi.Key)
	assert.NotEmpty(t, fi.Leaf)
	assert.NotEmpty(t, fi.CA)
	assert.NotEmpty(t, fi.ID.Bytes())

	assert.Equal(t, expectedFI.Key, fi.Key)
	assert.Equal(t, expectedFI.Leaf, fi.Leaf)
	assert.Equal(t, expectedFI.CA, fi.CA)
	assert.Equal(t, expectedFI.ID.Bytes(), fi.ID.Bytes())
}

func TestNodeID_Difficulty(t *testing.T) {
	_, fi, knownDifficulty, done := tempIdentity(t)
	defer done()

	difficulty := fi.ID.Difficulty()
	assert.True(t, difficulty >= knownDifficulty)
}

func TestVerifyPeer(t *testing.T) {
	check := newChecker(t)

	ca, err := NewCA(context.Background(), 12, 4)
	check.e(err)
	fi, err := ca.NewIdentity()
	check.e(err)

	err = peertls.VerifyPeerFunc(peertls.VerifyPeerCertChains)([][]byte{fi.Leaf.Raw, fi.CA.Raw}, nil)
	assert.NoError(t, err)
}

func TestCheckPeerRevocations(t *testing.T)  {
	// TODO
	t.SkipNow()
}
