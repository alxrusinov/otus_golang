package main

import (
	"bytes"
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrNegativeOffset        = errors.New("negative offset")
	ErrNegativeLimit         = errors.New("negative limit")
)

func prepareFiles(fromPath, toPath string) (fromFile, toFile *os.File, err error) {
	fromFile, fromError := os.Open(fromPath)

	if fromError != nil {
		return nil, nil, fromError
	}

	toFile, toError := os.Create(toPath)

	if toError != nil {
		return nil, nil, toError
	}

	return fromFile, toFile, nil
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	if offset < 0 {
		return ErrNegativeOffset
	}

	if limit < 0 {
		return ErrNegativeLimit
	}

	fromFile, toFile, err := prepareFiles(fromPath, toPath)
	if err != nil {
		return err
	}

	defer fromFile.Close()
	defer toFile.Close()

	readingLimit := limit

	info, errStat := fromFile.Stat()

	if err != nil {
		return errStat
	}

	size := info.Size()

	if size == 0 {
		return ErrUnsupportedFile
	}

	if offset > size {
		return ErrOffsetExceedsFileSize
	}

	if readingLimit == 0 {
		readingLimit = size
	}

	if readingLimit > size-offset {
		readingLimit = size - offset
	}

	buf := make([]byte, readingLimit)
	_, errRead := fromFile.ReadAt(buf, offset)

	if errRead != nil && errRead != io.EOF {
		return errRead
	}

	bf := bytes.NewReader(buf)

	bar := pb.Full.Start64(readingLimit)
	barReader := bar.NewProxyReader(bf)
	defer bar.Finish()

	_, errCopy := io.Copy(toFile, barReader)

	if err != nil {
		return errCopy
	}

	return nil
}
