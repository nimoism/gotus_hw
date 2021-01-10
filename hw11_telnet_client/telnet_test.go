package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, ioutil.NopCloser(in), out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString("hello\n")
			err = client.Send()
			require.NoError(t, err)

			err = client.Receive()
			require.NoError(t, err)
			require.Equal(t, "world\n", out.String())
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\n", string(request)[:n])

			n, err = conn.Write([]byte("world\n"))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})

	realStdin := os.Stdin
	defer func() { os.Stdin = realStdin }()
	stdin, fakeStdin, err := os.Pipe()
	require.NoError(t, err)
	os.Stdin = fakeStdin

	realStderr := os.Stderr
	defer func() { os.Stderr = realStderr }()
	stderr, fakeStderr, err := os.Pipe()
	require.NoError(t, err)
	os.Stderr = fakeStderr

	t.Run("server receive disconnect", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			out := &bytes.Buffer{}

			errScanner := bufio.NewScanner(stderr)

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, stdin, out)
			require.NoError(t, client.Connect())

			require.True(t, errScanner.Scan())
			require.Equal(t, fmt.Sprintf("...Connected to %v", l.Addr().String()), errScanner.Text())

			err = client.Receive()
			require.NoError(t, err)

			require.True(t, errScanner.Scan())
			require.Equal(t, "...Connection was closed by peer", errScanner.Text())

		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)

			require.NoError(t, conn.Close())
		}()

		wg.Wait()
	})

	t.Run("server send disconnect", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			out := &bytes.Buffer{}
			errScanner := bufio.NewScanner(stderr)

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, stdin, out)
			require.NoError(t, client.Connect())

			require.True(t, errScanner.Scan())
			require.Equal(t, fmt.Sprintf("...Connected to %v", l.Addr().String()), errScanner.Text())

			err = client.Send()
			require.NoError(t, err)

			require.True(t, errScanner.Scan())
			require.Equal(t, "...EOF", errScanner.Text())

		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			require.NoError(t, fakeStdin.Close())
		}()

		wg.Wait()
	})

}
