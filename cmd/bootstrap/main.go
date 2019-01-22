// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"storj.io/storj/internal/fpath"
	"storj.io/storj/pkg/cfgstruct"
	"storj.io/storj/pkg/kademlia"
	"storj.io/storj/pkg/process"
	"storj.io/storj/pkg/server"
)

// Bootstrap defines a bootstrap node configuration
type Bootstrap struct {
	Server   server.Config
	Kademlia kademlia.BootstrapConfig
}

var (
	rootCmd = &cobra.Command{
		Use:   "bootstrap",
		Short: "bootstrap",
	}
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run the bootstrap server",
		RunE:  cmdRun,
	}
	setupCmd = &cobra.Command{
		Use:         "setup",
		Short:       "Create config files",
		RunE:        cmdSetup,
		Annotations: map[string]string{"type": "setup"},
	}

	cfg Bootstrap

	defaultConfDir     = fpath.ApplicationDir("storj", "bootstrap")
	defaultIdentityDir = fpath.ApplicationDir("storj", "identity", "bootstrap")
	confDir            string
	identityDir        string
)

const (
	defaultServerAddr = ":28967"
)

func init() {
	confDirParam := cfgstruct.FindConfigDirParam()
	if confDirParam != "" {
		defaultConfDir = confDirParam
	}
	identityDirParam := cfgstruct.FindCredsDirParam()
	if identityDirParam != "" {
		defaultIdentityDir = identityDirParam
	}

	rootCmd.PersistentFlags().StringVar(&confDir, "config-dir", defaultConfDir, "main directory for bootstrap configuration")
	err := rootCmd.PersistentFlags().SetAnnotation("config-dir", "setup", []string{"true"})
	if err != nil {
		zap.S().Error("Failed to set 'setup' annotation for 'config-dir'")
	}
	rootCmd.PersistentFlags().StringVar(&identityDir, "identity-dir", defaultIdentityDir, "main directory for bootstrap identity credentials")
	err = rootCmd.PersistentFlags().SetAnnotation("identity-dir", "setup", []string{"true"})
	if err != nil {
		zap.S().Error("Failed to set 'setup' annotation for 'config-dir'")
	}

	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(setupCmd)
	cfgstruct.Bind(runCmd.Flags(), &cfg, cfgstruct.ConfDir(defaultConfDir), cfgstruct.CredsDir(defaultIdentityDir))
	cfgstruct.BindSetup(setupCmd.Flags(), &cfg, cfgstruct.ConfDir(defaultConfDir), cfgstruct.CredsDir(defaultIdentityDir))
}

func cmdRun(cmd *cobra.Command, args []string) (err error) {
	ctx := process.Ctx(cmd)
	if _, err := cfg.Server.Identity.Load(); err != nil {
		zap.S().Fatal(err)
	}

	if err := process.InitMetricsWithCertPath(ctx, nil, cfg.Server.Identity.CertPath); err != nil {
		zap.S().Errorf("Failed to initialize telemetry batcher: %+v", err)
	}
	return cfg.Server.Run(ctx, nil, cfg.Kademlia)
}

func cmdSetup(cmd *cobra.Command, args []string) (err error) {
	setupDir, err := filepath.Abs(confDir)
	if err != nil {
		return err
	}

	valid, _ := fpath.IsValidSetupDir(setupDir)
	if !valid {
		return fmt.Errorf("bootstrap configuration already exists (%v)", setupDir)
	}

	err = os.MkdirAll(setupDir, 0700)
	if err != nil {
		return err
	}

	overrides := map[string]interface{}{
		"identity.server.address": defaultServerAddr,
		"kademlia.bootstrap-addr": "localhost" + defaultServerAddr,
	}

	return process.SaveConfigWithAllDefaults(cmd.Flags(), filepath.Join(setupDir, "config.yaml"), overrides)
}

func main() {
	process.Exec(rootCmd)
}
