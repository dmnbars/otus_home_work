package main

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopy(t *testing.T) {
	t.Run("source file does not exist", func(t *testing.T) {
		err := Copy(
			"testdata/does_not_exists_input.txt",
			"/tmp/output.txt",
			0,
			0,
		)

		assert.ErrorIs(t, err, os.ErrNotExist)
	})
	t.Run("destination file already exist", func(t *testing.T) {
		err := Copy(
			"testdata/input.txt",
			"testdata/input.txt",
			0,
			0,
		)

		assert.ErrorIs(t, err, os.ErrExist)
	})
	t.Run("offset exceeds file size", func(t *testing.T) {
		err := Copy(
			"testdata/out_offset0_limit10.txt",
			"/tmp/output.txt",
			1000,
			0,
		)

		assert.ErrorIs(t, err, ErrOffsetExceedsFileSize)
	})
	t.Run("unsupported file", func(t *testing.T) {
		err := Copy(
			"/dev/urandom",
			"/tmp/output.txt",
			0,
			0,
		)

		assert.ErrorIs(t, err, ErrUnsupportedFile)
	})

	tests := []struct {
		input  string
		offset int64
		limit  int64
	}{
		{
			input:  "testdata/input.txt",
			offset: 0,
			limit:  0,
		},
		{
			input:  "testdata/input.txt",
			offset: 0,
			limit:  10,
		},
		{
			input:  "testdata/input.txt",
			offset: 0,
			limit:  1000,
		},
		{
			input:  "testdata/input.txt",
			offset: 0,
			limit:  10000,
		},
		{
			input:  "testdata/input.txt",
			offset: 100,
			limit:  1000,
		},
		{
			input:  "testdata/input.txt",
			offset: 6000,
			limit:  1000,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("copy, offset: %d, limit: %d", tc.offset, tc.limit), func(t *testing.T) {
			actualPath := fmt.Sprintf("/tmp/out_offset%d_limit%d.txt", tc.offset, tc.limit)
			expectedPath := fmt.Sprintf("testdata/out_offset%d_limit%d.txt", tc.offset, tc.limit)
			defer func() {
				err := os.Remove(actualPath)
				if err != nil {
					t.Errorf("failed to remove output file: %s", err)
				}
			}()

			err := Copy(
				"testdata/input.txt",
				actualPath,
				tc.offset,
				tc.limit,
			)

			assert.NoError(t, err)
			assertFilesEqual(t, expectedPath, actualPath)
		})
	}
}

func assertFilesEqual(t assert.TestingT, expectedPath, actualPath string) bool {
	expected, err := readAll(expectedPath)
	if err != nil {
		t.Errorf("failed to read expected file: %s", err)
	}

	actual, err := readAll(actualPath)
	if err != nil {
		t.Errorf("failed to read actual file: %s", err)
	}

	return assert.Equal(t, expected, actual)
}

func readAll(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return io.ReadAll(f)
}
