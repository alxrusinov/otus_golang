package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	pwd, _ := os.Getwd()

	tc := []struct {
		name       string
		cmd        []string
		env        Environment
		returnCode int
	}{
		{
			name: "success case",
			env: Environment{
				"BAR":   {"bar", false},
				"EMPTY": {"", false},
				"FOO":   {"   foo\nwith new line false", false},
				"HELLO": {"hello", false}, "UNSET": {"", true},
			},
			cmd: []string{
				"/bin/bash",
				pwd + "/testdata/echo.sh",
				"arg1=1",
				"arg2=2",
			},
			returnCode: 0,
		},
		{
			name: "wrong cmd",
			env: Environment{
				"BAR":   {"bar", false},
				"EMPTY": {"", false},
				"FOO":   {"   foo\nwith new line false", false},
				"HELLO": {"hello", false}, "UNSET": {"", true},
			},
			cmd: []string{
				"/bin/bash",
				pwd + "/testdata/echo.shf",
				"arg1=1",
				"arg2=2",
			},
			returnCode: 127,
		},
	}

	for _, tCase := range tc {
		t.Run(tCase.name, func(t *testing.T) {
			returnCode := RunCmd(tCase.cmd, tCase.env)

			require.Equalf(t,
				tCase.returnCode, returnCode,
				"expected returnCode - %v, actual - %v",
				tCase.returnCode,
				returnCode)
		})
	}
}
