// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/spf13/cobra"
	"github.com/urfave/cli"
	"github.com/zeebo/errs"

	"storj.io/storj/pkg/piecestore"
	"storj.io/storj/pkg/process"
)

var argError = errs.Class("argError")

func run(_ *cobra.Command, args []string) error {
	app := cli.NewApp()

	app.Name = "Piece Store CLI"
	app.Usage = "Store data in hash folder structure"
	app.Version = "1.0.0"

	app.Flags = []cli.Flag{}

	app.Commands = []cli.Command{
		{
			Name:      "store",
			Aliases:   []string{"s"},
			Usage:     "Store data by id",
			ArgsUsage: "[id] [dataPath] [storeDir]",
			Action: func(c *cli.Context) error {
				if c.Args().Get(0) == "" {
					return argError.New("No id specified")
				}

				id := c.Args().Get(0)

				if c.Args().Get(1) == "" {
					return argError.New("No input file specified")
				}

				path := c.Args().Get(1)

				if c.Args().Get(2) == "" {
					return argError.New("No output directory specified")
				}

				outputDir := c.Args().Get(2)

				file, err := os.Open(path)
				if err != nil {
					return err
				}

				// Close the file when we are done
				defer file.Close()

				fileInfo, err := os.Stat(path)
				if err != nil {
					return err
				}

				if fileInfo.IsDir() {
					return argError.New(fmt.Sprintf("Path (%s) is a directory, not a file", path))
				}

				dataFileChunk, err := pstore.StoreWriter(id, outputDir)
				if err != nil {
					return err
				}

				// Close when finished
				defer dataFileChunk.Close()

				_, err = io.Copy(dataFileChunk, file)

				return err
			},
		},
		{
			Name:      "retrieve",
			Aliases:   []string{"r"},
			Usage:     "Retrieve data by id and print to Stdout",
			ArgsUsage: "[id] [storeDir]",
			Action: func(c *cli.Context) error {
				if c.Args().Get(0) == "" {
					return argError.New("Missing data id")
				}
				if c.Args().Get(1) == "" {
					return argError.New("Missing file path")
				}
				fileInfo, err := os.Stat(c.Args().Get(1))
				if err != nil {
					return err
				}

				if fileInfo.IsDir() != true {
					return argError.New(fmt.Sprintf("Path (%s) is a file, not a directory", c.Args().Get(1)))
				}

				dataFileChunk, err := pstore.RetrieveReader(context.Background(),
					c.Args().Get(0), 0, -1, c.Args().Get(1))
				if err != nil {
					return err
				}

				// Close when finished
				defer dataFileChunk.Close()

				_, err = io.Copy(os.Stdout, dataFileChunk)
				return err
			},
		},
		{
			Name:      "delete",
			Aliases:   []string{"d"},
			Usage:     "Delete data by id",
			ArgsUsage: "[id] [storeDir]",
			Action: func(c *cli.Context) error {
				if c.Args().Get(0) == "" {
					return argError.New("Missing data id")
				}
				if c.Args().Get(1) == "" {
					return argError.New("No directory specified")
				}
				err := pstore.Delete(c.Args().Get(0), c.Args().Get(1))

				return err
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	return app.Run(append([]string{os.Args[0]}, args...))
}

func main() {
	process.Exec(&cobra.Command{
		Use:   "piecestore-cli",
		Short: "piecestore example cli",
		RunE:  run,
	})
}
