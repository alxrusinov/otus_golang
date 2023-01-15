package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	tc := []struct {
		name      string
		dir       string
		err       error
		resultLen int
	}{
		{
			name:      "dir is not exist",
			dir:       "testdata/envi",
			err:       errors.New("err"),
			resultLen: 0,
		},
		{
			name:      "success case",
			dir:       "testdata/env",
			err:       nil,
			resultLen: 5,
		},
	}

	for _, tCase := range tc {
		t.Run(tCase.name, func(t *testing.T) {
			val, err := ReadDir(tCase.dir)

			if tCase.err == nil {
				require.Lenf(t, val, tCase.resultLen, "expected length - %v, actual - %v", tCase.resultLen, len(val))
				return
			}

			require.Error(t, err, "function must return error")
		})
	}
}
