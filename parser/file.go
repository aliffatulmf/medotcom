package parser

import (
	"errors"
	"os"
)

const (
	TEXT string = ".txt"
	JSON        = ".json"
)

var ErrFileTooLarge = errors.New("file too large")

func CheckFileSize(file *os.File, limit int64) error {
	stat, err := file.Stat()
	if err != nil {
		return err
	}

	if stat.Size() > limit {
		return ErrFileTooLarge
	}
	return nil
}
