package main

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"strings"
	"unicode"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	envs := Environment{}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}

		if info.Size() == 0 {
			envs[entry.Name()] = EnvValue{
				Value:      "",
				NeedRemove: true,
			}
			continue
		}

		file, err := os.Open(dir + "/" + info.Name())
		if err != nil {
			return nil, err
		}

		reader := bufio.NewReader(file)
		line, _, err := reader.ReadLine()
		if err != nil {
			return nil, err
		}
		if err := file.Close(); err != nil {
			log.Printf("error while closing file: %s", err)
		}

		replace := bytes.ReplaceAll(line, []byte("\x00"), []byte("\n"))
		value := strings.TrimRightFunc(string(replace), unicode.IsSpace)

		envs[info.Name()] = EnvValue{
			Value: value,
		}
	}

	return envs, nil
}
