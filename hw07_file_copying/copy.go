package main

import (
	"errors"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	sourceStat, err := os.Stat(fromPath)
	if err != nil {
		return err
	}
	if sourceStat.Size() == 0 {
		return ErrUnsupportedFile
	}
	if sourceStat.Size() < offset {
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

	destination, err := os.OpenFile(toPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer func() {
		if err := destination.Close(); err != nil {
			panic(err)
		}
	}()

	if limit == 0 {
		_, err = io.Copy(destination, source)
	} else {
		_, err = io.CopyN(destination, source, limit)
	}
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	return nil
}
