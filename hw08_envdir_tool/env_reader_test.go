package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	envs, err := ReadDir("testdata/env/")
	require.NoError(t, err)
	expected := make(Environment)
	for _, kv := range [][2]string{
		{"BAR", "bar"},
		{"FOO", "   foo\nwith new line"},
		{"HELLO", `"hello"`},
	} {
		expected[kv[0]] = kv[1]
	}
	require.Equal(t, expected, envs)
}

func TestCleanValue(t *testing.T) {
	for _, tst := range []struct {
		value string
		exp   string
	}{
		{value: "", exp: ""},
		{value: "=", exp: "="},
		{value: "A=A", exp: "A=A"},
		{value: "AAA\n", exp: "AAA"},
		{value: "AAA\x00BB\n", exp: "AAA\nBB"},
		{value: "  AAA\x00BB  \n", exp: "  AAA\nBB"},
	} {
		value := normalizeValue(tst.value)
		require.Equal(t, tst.exp, value)
	}
}

func TestEnvironment_Strings(t *testing.T) {
	env := Environment{
		"BAR":   "bar",
		"FOO":   "   foo\nwith new line",
		"HELLO": `"hello"`,
	}
	exp := []string{
		"BAR=bar",
		"FOO=   foo\nwith new line",
		`HELLO="hello"`,
	}
	require.ElementsMatch(t, exp, env.Strings())
}
