// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package main

import (
	"io/ioutil"
	"encoding/json"

	"github.com/spf13/cobra"
	"github.com/zeebo/errs"

	"storj.io/storj/pkg/cfgstruct"
	"storj.io/storj/pkg/process"
	"fmt"
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
		NodesPath string
	}

	clearCfg struct {
		ExceptPath string
	}
)

func init() {
	rootCmd.AddCommand(addCmd)
	cfgstruct.Bind(addCmd.Flags(), &addCfg)
	cfgstruct.Bind(clearCmd.Flags(), &clearCfg)
}

func cmdAdd(cmd *cobra.Command, args []string) (err error) {
	j, err := ioutil.ReadFile(addCfg.NodesPath)
	if err != nil {
		return errs.Wrap(err)
	}

	type id string
	type address string
	var nodes map[id]address
	if err := json.Unmarshal(j, &nodes); err != nil {
		return errs.Wrap(err)
	}

	// TODO add records to cache
	fmt.Println(nodes)

	return nil
}

func cmdClear(cmd *cobra.Command, args []string) (err error) {
	// TODO
	return nil
}

func main() {
	process.Exec(rootCmd)
}
