// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package tlsopts_test

import (
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"storj.io/storj/internal/testcontext"
	"storj.io/storj/internal/testplanet"
	"storj.io/storj/pkg/peertls"
	"storj.io/storj/pkg/peertls/tlsopts"
)

var pregeneratedIdentities = testplanet.NewPregeneratedIdentities()

func TestNewOptions(t *testing.T) {
	ctx := testcontext.New(t)
	defer ctx.Cleanup()

	fi, err := pregeneratedIdentities.NewIdentity()
	require.NoError(t, err)

	whitelistPath := ctx.File("whitelist.pem")

	chainData, err := peertls.ChainBytes(fi.CA)
	assert.NoError(t, err)

	err = ioutil.WriteFile(whitelistPath, chainData, 0644)
	assert.NoError(t, err)

	cases := []struct {
		testID      string
		config      tlsopts.Config
		pcvFuncsLen int
	}{
		{
			"default",
			tlsopts.Config{},
			0,
		}, {
			"revocation processing",
			tlsopts.Config{
				RevocationDBURL: "bolt://" + ctx.File("revocation1.db"),
				Extensions: peertls.TLSExtConfig{
					Revocation: true,
				},
			},
			2,
		}, {
			"ca whitelist verification",
			tlsopts.Config{
				PeerCAWhitelistPath: whitelistPath,
				UsePeerCAWhitelist:  true,
			},
			1,
		}, {
			"ca whitelist verification and whitelist signed leaf verification",
			tlsopts.Config{
				// NB: file doesn't actually exist
				PeerCAWhitelistPath: whitelistPath,
				UsePeerCAWhitelist:  true,
				Extensions: peertls.TLSExtConfig{
					WhitelistSignedLeaf: true,
				},
			},
			2,
		}, {
			"revocation processing and whitelist verification",
			tlsopts.Config{
				// NB: file doesn't actually exist
				PeerCAWhitelistPath: whitelistPath,
				UsePeerCAWhitelist:  true,
				RevocationDBURL:     "bolt://" + ctx.File("revocation2.db"),
				Extensions: peertls.TLSExtConfig{
					Revocation: true,
				},
			},
			3,
		}, {
			"revocation processing, whitelist, and signed leaf verification",
			tlsopts.Config{
				// NB: file doesn't actually exist
				PeerCAWhitelistPath: whitelistPath,
				UsePeerCAWhitelist:  true,
				RevocationDBURL:     "bolt://" + ctx.File("revocation3.db"),
				Extensions: peertls.TLSExtConfig{
					Revocation:          true,
					WhitelistSignedLeaf: true,
				},
			},
			3,
		},
	}

	for _, c := range cases {
		t.Log(c.testID)
		opts, err := tlsopts.NewOptions(fi, c.config)
		assert.NoError(t, err)
		assert.True(t, reflect.DeepEqual(fi, opts.Ident))
		assert.Equal(t, c.config, opts.Config)
		assert.Len(t, opts.PCVFuncs, c.pcvFuncsLen)
	}
}

func TestPeerCAWhitelist(t *testing.T) {
	t.Run("all nodes signed", func(t *testing.T) {
		testplanet.Run(t, testplanet.Config{
			SatelliteCount:   1,
			StorageNodeCount: 10,
			UplinkCount:      0,
			UsePeerCAWhitelist: true,
		}, func(t *testing.T, ctx *testcontext.Context, planet *testplanet.Planet) {
			err := planet.Ping(ctx)
			assert.NoError(t, err)
		})
	})

	t.Run("all nodes unsigned", func(t *testing.T) {
		testplanet.Run(t, testplanet.Config{
			SatelliteCount:   1,
			StorageNodeCount: 0,
			UplinkCount:      0,
			UsePeerCAWhitelist: true,
			Identities: testplanet.NewPregeneratedIdentities(),
		}, func(t *testing.T, ctx *testcontext.Context, planet *testplanet.Planet) {
			//_ = planet.Ping(ctx)
			//fmt.Println(err)
			//assert.Error(t, err)
			//assert.True(t, true)
		})
	})

	t.Run("unsigned satellite", func(t *testing.T) {
		testIdentities, err := testplanet.MixedIdentities([]int{1,2,3}, testplanet.MixedIndexesUnsigned)
		require.NoError(t, err)

		testplanet.Run(t, testplanet.Config{
			SatelliteCount:   1,
			StorageNodeCount: 1,
			UplinkCount:      0,
			UsePeerCAWhitelist: true,
			Identities:       testIdentities,
			//Reconfigure: testplanet.Reconfigure{
			//	Satellite: func(_ int, cfg *satellite.Config) {
			//		cfg.Server.UsePeerCAWhitelist = false
			//	},
			//},
		}, func(t *testing.T, ctx *testcontext.Context, planet *testplanet.Planet) {
			err := planet.Ping(ctx)
			assert.NoError(t, err)
		})
	})

	// wip - more mixed scenarios...
}

