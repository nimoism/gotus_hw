package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	env, err := ReadDir("testdata/env")
	var exitCode int
	for _, tst := range []struct {
		cmd  []string
		out  string
		code int
	}{
		{
			cmd:  []string{"sh", "-c", "echo $HELLO $BAR"},
			out:  "\"hello\" bar\n",
			code: 0,
		},
		{
			cmd:  []string{"wrong-cmd"},
			out:  "",
			code: ExitCodeIOError,
		},
		{
			cmd:  []string{},
			out:  "",
			code: ExitCodeCommandNotFound,
		},
	} {
		tst := tst
		var buf bytes.Buffer
		var r *os.File
		require.NoError(t, err)
		func() {
			var w *os.File
			r, w, err = os.Pipe()
			defer w.Close()
			origStdout := os.Stdout
			defer func() { os.Stdout = origStdout }()
			os.Stdout = w
			require.NoError(t, err)
			exitCode = RunCmd(tst.cmd, env)
		}()
		_, err = buf.ReadFrom(r)
		require.NoError(t, err)
		require.Equal(t, tst.code, exitCode)
		require.Equal(t, tst.out, buf.String())

	}
}
