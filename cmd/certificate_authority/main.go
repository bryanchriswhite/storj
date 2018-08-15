// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package certificate_authority

import (
	"fmt"

	"github.com/spf13/cobra"

	"storj.io/storj/pkg/cfgstruct"
	"storj.io/storj/pkg/provider"
	"storj.io/storj/pkg/process"
)

var (
	RootCmd = &cobra.Command{
		Use:   "ca",
		Short: "Certificate authority",
	}
	newCmd = &cobra.Command{
		Use:   "new",
		Short: "Create a new certificate authority",
		RunE:  cmdNew,
	}
	idCmd = &cobra.Command{
		Use:   "id",
		Short: "Get the id of a CA",
		RunE:  cmdID,
	}

	newCACfg struct {
		CA provider.CASetupConfig
	}

	idCfg struct {
		CA provider.PeerCAConfig
	}
)

func init() {
	RootCmd.AddCommand(newCmd)
	RootCmd.AddCommand(idCmd)
	cfgstruct.Bind(newCmd.Flags(), &newCACfg)
	cfgstruct.Bind(idCmd.Flags(), &idCfg)
}

func cmdNew(cmd *cobra.Command, args []string) (err error) {
	_, err = provider.SetupCA(process.Ctx(cmd), newCACfg.CA)
	if err != nil {
		return err
	}

	return nil
}

func cmdID(cmd *cobra.Command, args []string) (err error) {
	p, err := idCfg.CA.Load()
	if err != nil {
		return err
	}

	fmt.Println(p.ID.String())
	return nil
}

func main() {
	process.Exec(RootCmd)
}
