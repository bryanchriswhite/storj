// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package overlay

import (
	"github.com/spf13/cobra"
	"storj.io/storj/pkg/cfgstruct"
)

var (
	rootCmd = &cobra.Command{
		Use:   "overlay",
		Short: "Overlay cache management",
	}
	addCmd = &cobra.Command{
		Use:   "add",
		Short: "Add nodes to the overlay cache",
		RunE:  cmdAdd,
	}
	clearCmd = &cobra.Command{
		Use:   "clear",
		Short: "Clear the overlay cache",
		RunE:  cmdClear,
	}

	addCfg struct {
		nodeListPath string
	}

	clearCfg struct {}
)

func init() {
	rootCmd.AddCommand(newCACmd)
	cfgstruct.Bind(newCACmd.Flags(), &newCACfg)
	cfgstruct.Bind(idCmd.Flags(), &idCfg)
}

func cmdAdd(cmd *cobra.Command, args []string) (err error) {
	// TODO
	return nil
}

func cmdClear(cmd *cobra.Command, args []string) (err error) {
	// TODO
	return nil
}
