package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	wrongEnvChars = "="
	trimChars     = " \n\t"
)

var ErrWrongEnvValue = errors.New("wrong env value")

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
	d, err := os.Open(dir)
	if err != nil {
		return nil, fmt.Errorf("read dir error: %w", err)
	}
	defer d.Close()
	var fileNames []string
	fileNames, err = d.Readdirnames(0)
	if err != nil {
		return nil, fmt.Errorf("read dir error: %w", err)
	}
	envs := make(Environment, len(fileNames))
	var value string
	for _, fn := range fileNames {
		value, err = func(fn string) (string, error) {
			f, err := os.Open(filepath.Join(dir, fn))
			if err != nil {
				return "", fmt.Errorf("read dir error: %w", err)
			}
			defer f.Close()
			stat, err := f.Stat()
			if err != nil {
				return "", fmt.Errorf("read dir error: %w", err)
			}
			if !stat.Mode().IsRegular() || stat.Mode().IsDir() {
				return "", nil
			}
			scanner := bufio.NewScanner(f)
			scanner.Split(bufio.ScanLines)
			if !scanner.Scan() {
				return "", nil
			}
			value = scanner.Text()
			return value, nil
		}(fn)
		if err != nil {
			return nil, fmt.Errorf("read dir error: %w", err)
		}
		if value, err = normalizeValue(value); err != nil {
			return nil, err
		}
		if value != "" {
			envs[fn] = value
		}
	}
	return envs, nil
}

func normalizeValue(value string) (string, error) {
	if value == "" {
		return value, nil
	}
	if strings.ContainsAny(value, wrongEnvChars) {
		return "", ErrWrongEnvValue
	}
	value = strings.TrimRight(value, trimChars)
	value = strings.ReplaceAll(value, "\x00", "\n")
	return value, nil
}
