// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package ranger

import (
	"io"
	"os"

	"go.uber.org/zap"
)

// FileHandleRanger returns a RangeCloser from a file handle. The
// Closer's Close method will call fh.Close().
// Footgun: If FileHandleRanger fails, fh.Close will not have been called.
func FileHandleRanger(fh *os.File) (RangeCloser, error) {
	stat, err := fh.Stat()
	if err != nil {
		return nil, Error.Wrap(err)
	}
	return struct {
		Ranger
		io.Closer
	}{
		Ranger: ReaderAtRanger(fh, stat.Size()),
		Closer: fh,
	}, nil
}

// FileRanger returns a RangeCloser from a path.
func FileRanger(path string) (RangeCloser, error) {
	fh, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	r, err := FileHandleRanger(fh)
	if err != nil {
		if closeErr := fh.Close(); closeErr != nil {
			zap.S().Error(closeErr)
		}

		return nil, err
	}
	return r, nil
}
