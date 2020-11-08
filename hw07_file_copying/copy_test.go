package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const dataPath = "testdata"

var srcPath = filepath.Join(dataPath, "input.txt")

func TestCopyRegular(t *testing.T) {
	for _, tst := range []struct {
		offset int64
		limit  int64
	}{
		{offset: 0, limit: 0},
		{offset: 0, limit: 10},
		{offset: 0, limit: 1000},
		{offset: 0, limit: 10000},
		{offset: 100, limit: 1000},
		{offset: 6000, limit: 1000},
	} {
		tst := tst
		t.Run(fmt.Sprintf("%v_%v", offset, limit), func(t *testing.T) {
			dst, err := ioutil.TempFile(os.TempDir(), "go-copy")
			require.NoError(t, err)
			defer dst.Close()
			err = Copy(srcPath, dst.Name(), tst.offset, tst.limit)
			require.NoError(t, err)
			expPath := filepath.Join(dataPath, fmt.Sprintf("out_offset%d_limit%d.txt", tst.offset, tst.limit))
			fmt.Println(expPath)
			expContent, err := ioutil.ReadFile(expPath)
			require.NoError(t, err)
			dstContent, err := ioutil.ReadFile(dst.Name())
			require.NoError(t, err)
			require.Zero(t, bytes.Compare(expContent, dstContent))
		})
	}
}

func TestCopyNonRegular(t *testing.T) {
	dst, err := ioutil.TempFile(os.TempDir(), "go-copy")
	require.NoError(t, err)
	defer dst.Close()
	err = Copy("/dev/urandom", dst.Name(), 0, 1000)
	require.NoError(t, err)
	stat, err := dst.Stat()
	require.NoError(t, err)
	require.Equal(t, 1000, int(stat.Size()))
}

func TestCopyNonRegularNoLimitErr(t *testing.T) {
	dst, err := ioutil.TempFile(os.TempDir(), "go-copy")
	require.NoError(t, err)
	defer dst.Close()
	err = Copy("/dev/urandom", dst.Name(), 0, 0)
	require.EqualError(t, err, ErrUnsupportedFile.Error())
}
