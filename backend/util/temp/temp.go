package temp

import (
	"errors"
	"io"
	"os"
)

type TempFile struct {
	file *os.File
}

func (t TempFile) File() *os.File {
	return t.file
}

func (t TempFile) Close() error {
	var joinedErr error
	if err := t.file.Close(); err != nil {
		joinedErr = err
	}
	if err := os.Remove(t.file.Name()); err != nil {
		joinedErr = errors.Join(joinedErr, err)
	}

	return nil
}

func CreateTempFile(r io.Reader, filename string) (TempFile, error) {
	inputFile, err := os.CreateTemp("", filename)
	if err != nil {
		return TempFile{}, err
	}
	_, err = io.Copy(inputFile, r)
	if err != nil {
		return TempFile{}, err
	}
	if err := inputFile.Sync(); err != nil {
		return TempFile{}, err
	}
	if _, err := inputFile.Seek(0, 0); err != nil {
		return TempFile{}, err
	}

	return TempFile{
		file: inputFile,
	}, nil
}
