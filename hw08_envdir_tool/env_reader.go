package main

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

func prepareValue(value []byte) []string {
	splinted := bytes.Split(value, []byte("\n"))
	result := make([]string, len(splinted))

	for i, v := range splinted {
		replaced := bytes.ReplaceAll(v, []byte{0}, []byte("\n"))
		trimmed := bytes.TrimRight(replaced, " \t")
		stringifies := string(trimmed)

		result[i] = stringifies
	}

	return result
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	envs := make(Environment)

	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		var envValue EnvValue

		if strings.Contains(info.Name(), "=") {
			return errors.New("file name cannot contain =")
		}

		if info.Size() == 0 {
			envValue.NeedRemove = true
			envs[info.Name()] = envValue
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}

		defer file.Close()

		value := make([]byte, info.Size())

		_, err = file.Read(value)

		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}

		preparedValue := prepareValue(value)

		envValue.Value = preparedValue[0]

		envs[info.Name()] = envValue

		return nil
	})
	if err != nil {
		return nil, err
	}

	return envs, nil
}
