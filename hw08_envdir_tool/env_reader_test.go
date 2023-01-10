package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("not exists dir", func(t *testing.T) {
		envs, err := ReadDir("/testdata/not_exists_dir")
		assert.Error(t, err)
		assert.Nil(t, envs)
	})

	t.Run("empty dir", func(t *testing.T) {
		dirPath := createDir(t)

		envs, err := ReadDir(dirPath)
		assert.NoError(t, err)
		assert.Equal(t, Environment{}, envs)
	})

	tests := map[string]struct {
		file     string
		expected Environment
	}{
		"one line file": {
			file: "HELLO",
			expected: Environment{"HELLO": EnvValue{
				Value: "\"hello\"",
			}},
		},
		"file with second line": {
			file: "BAR",
			expected: Environment{"BAR": EnvValue{
				Value: "bar",
			}},
		},
		"empty file": {
			file:     "EMPTY",
			expected: Environment{"EMPTY": EnvValue{}},
		},
		"file with size 0": {
			file: "UNSET",
			expected: Environment{"UNSET": EnvValue{
				NeedRemove: true,
			}},
		},
		"file with terminal null": {
			file: "FOO",
			expected: Environment{"FOO": EnvValue{
				Value: "   foo\nwith new line",
			}},
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			dirPath, filePath := copyFile(t, tc.file)
			defer func() {
				require.NoError(t, os.Remove(filePath))
				require.NoError(t, os.Remove(dirPath))
			}()

			envs, err := ReadDir(dirPath)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, envs)
		})
	}

	t.Run("file with right tab", func(t *testing.T) {
		dirPath := createDir(t)

		file, err := os.Create(dirPath + "/TAB")
		require.NoError(t, err)

		_, err = file.WriteString("tab\t")
		require.NoError(t, err)
		require.NoError(t, file.Close())

		defer func() {
			require.NoError(t, os.Remove(dirPath+"/TAB"))
			require.NoError(t, os.Remove(dirPath))
		}()

		expected := Environment{"TAB": EnvValue{
			Value: "tab",
		}}

		envs, err := ReadDir(dirPath)
		assert.NoError(t, err)

		assert.Equal(t, expected, envs)
	})

	t.Run("multiple files", func(t *testing.T) {
		expected := Environment{
			"BAR": EnvValue{
				Value: "bar",
			},
			"EMPTY": EnvValue{},
			"FOO": EnvValue{
				Value: "   foo\nwith new line",
			},
			"HELLO": EnvValue{
				Value: "\"hello\"",
			},
			"UNSET": EnvValue{
				NeedRemove: true,
			},
		}

		envs, err := ReadDir("testdata/env")
		assert.NoError(t, err)

		assert.Equal(t, expected, envs)
	})
}

func createDir(t *testing.T) string {
	t.Helper()

	dirPath, err := os.MkdirTemp("/tmp", "test.")
	require.NoError(t, err)

	return dirPath
}

func copyFile(t *testing.T, src string) (string, string) {
	t.Helper()

	dirPath := createDir(t)
	filePath := dirPath + "/" + src
	data, err := os.ReadFile("testdata/env/" + src)
	require.NoError(t, err)

	err = os.WriteFile(filePath, data, 0o644)
	require.NoError(t, err)

	return dirPath, filePath
}
