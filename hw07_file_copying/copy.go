package main

import (
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
	ErrPathTheSame           = errors.New("from and to path the same")
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

	if fromPath == toPath {
		return ErrPathTheSame
	}

	fromFile, toFile, err := prepareFiles(fromPath, toPath)
	if err != nil {
		return err
	}

	defer fromFile.Close()
	defer toFile.Close()

	info, err := fromFile.Stat()
	if err != nil {
		return err
	}

	size := info.Size()

	if size == 0 {
		return ErrUnsupportedFile
	}

	if offset > size {
		return ErrOffsetExceedsFileSize
	}

	readingLimit := limit

	if readingLimit == 0 {
		readingLimit = size
	}

	if readingLimit > size-offset {
		readingLimit = size - offset
	}

	_, err = fromFile.Seek(offset, 0)

	if err != nil {
		return err
	}

	bar := pb.Full.Start64(readingLimit)
	barReader := bar.NewProxyReader(fromFile)
	defer bar.Finish()

	_, err = io.CopyN(toFile, barReader, readingLimit)
	if err != nil {
		return err
	}

	return nil
}
