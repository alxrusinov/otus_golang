package main

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	tc := []struct {
		name     string
		fromPath string
		toPath   string
		limit    int64
		offset   int64
		error    error
	}{
		{
			name:     "simple case",
			fromPath: "testdata/input.txt",
			toPath:   "out.txt",
			limit:    0,
			offset:   0,
			error:    nil,
		},
		{
			name:     "offset is over file size",
			fromPath: "testdata/input.txt",
			toPath:   "out.txt",
			limit:    0,
			offset:   7000,
			error:    ErrOffsetExceedsFileSize,
		},
		{
			name:     "unsupported file",
			fromPath: "/dev/urandom",
			toPath:   "out.txt",
			limit:    0,
			offset:   0,
			error:    ErrUnsupportedFile,
		},
		{
			name:     "open file error",
			fromPath: "testdata/input",
			toPath:   "out.txt",
			limit:    0,
			offset:   0,
			error:    errors.New("open testdata/input: no such file or directory"),
		},
		{
			name:     "negative offset",
			fromPath: "testdata/input.txt",
			toPath:   "out.txt",
			limit:    0,
			offset:   -2,
			error:    ErrNegativeOffset,
		},
		{
			name:     "negative limit",
			fromPath: "testdata/input.txt",
			toPath:   "out.txt",
			limit:    -2,
			offset:   0,
			error:    ErrNegativeLimit,
		},
		{
			name:     "the same paths",
			fromPath: "input.txt",
			toPath:   "input.txt",
			limit:    0,
			offset:   0,
			error:    ErrPathTheSame,
		},
	}

	for _, tCase := range tc {
		t.Run(tCase.name, func(t *testing.T) {
			err := Copy(tCase.fromPath, tCase.toPath, tCase.offset, tCase.limit)
			if tCase.error == nil {
				require.NoErrorf(t, err, "expected error %v, actual - %v", tCase.error, err)
				return
			}
			require.EqualError(t, err, tCase.error.Error())

			errRm := os.Remove(tCase.toPath)
			_ = errRm
		})
	}
}
