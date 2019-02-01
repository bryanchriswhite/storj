// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package utils

type Migration interface {
	Version() (int, error)
	Up(interface{}) error
	Down(interface{}) error
}
