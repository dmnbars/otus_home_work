package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrSameFile              = errors.New("from and to files are the same")
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

const defaultStep = 1024

func Copy(fromPath, toPath string, offset, limit int64) error {
	if fromPath == toPath {
		return ErrSameFile
	}
	sourceStat, err := os.Stat(fromPath)
	if err != nil {
		return err
	}
	sourceSize := sourceStat.Size()
	if sourceSize == 0 {
		return ErrUnsupportedFile
	}
	if sourceSize < offset {
		return ErrOffsetExceedsFileSize
	}

	source, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer func() {
		if err := source.Close(); err != nil {
			panic(err)
		}
	}()
	_, err = source.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	destination, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer func() {
		if err := destination.Close(); err != nil {
			panic(err)
		}
	}()

	if limit == 0 {
		limit = sourceSize
	} else if limit > sourceSize {
		limit = sourceSize
	}
	bar := pb.Start64(limit)
	defer bar.Finish()

	var step int64 = defaultStep
	if step > limit {
		step = limit
	}

	var copied int64
	for {
		if copied >= limit {
			break
		}
		n, err := io.CopyN(destination, source, step)
		copied += n
		bar.Add64(n)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
	}

	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	return nil
}
