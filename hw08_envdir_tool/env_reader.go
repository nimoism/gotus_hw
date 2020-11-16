package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	wrongEnvNameChars = "="
	trimChars         = " \n\t"
)

type Environment map[string]string

// Strings represents Environment as string slice like {"key0=value0", "key1=value1"}.
func (e Environment) Strings() []string {
	kvs := make([]string, 0, len(e))
	for key, value := range e {
		kvs = append(kvs, fmt.Sprintf("%v=%v", key, value))
	}
	return kvs
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	filesInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read dir error: %w", err)
	}
	envs := make(Environment, len(filesInfos))
	for _, fi := range filesInfos {
		if !fi.Mode().IsRegular() || fi.Mode().IsDir() {
			continue
		}
		if strings.ContainsAny(fi.Name(), wrongEnvNameChars) {
			continue
		}
		value, err := func(fi os.FileInfo) (string, error) {
			f, err := os.Open(filepath.Join(dir, fi.Name()))
			if err != nil {
				return "", fmt.Errorf("read %v file error: %w", fi.Name(), err)
			}
			defer f.Close()
			r := bufio.NewReader(f)
			var (
				buf        strings.Builder
				bytesValue []byte
				isPrefix   = true
			)
			for isPrefix {
				if bytesValue, isPrefix, err = r.ReadLine(); err != nil && !errors.Is(err, io.EOF) {
					return "", fmt.Errorf("read %v file error: %w", fi.Name(), err)
				}
				buf.Write(bytesValue)
			}
			return buf.String(), nil
		}(fi)
		if err != nil {
			return nil, err
		}
		value = normalizeValue(value)
		if value != "" {
			envs[fi.Name()] = value
		}
	}
	return envs, nil
}

func normalizeValue(value string) string {
	if value == "" {
		return value
	}
	value = strings.TrimRight(value, trimChars)
	value = strings.ReplaceAll(value, "\x00", "\n")
	return value
}
