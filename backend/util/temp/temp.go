package temp

import (
	"errors"
	"io"
	"os"
)

type TempFile struct {
	file   *os.File
	count  int
	closed bool
}

func (t *TempFile) UseFile() *os.File {
	t.count++
	return t.file
}

func (t *TempFile) Close() error {
	if t.closed {
		return nil
	}

	t.count--
	if t.count == 0 {
		var joinErr error
		if err := t.file.Close(); err != nil {
			joinErr = err
		}
		if err := os.Remove(t.file.Name()); err != nil {
			joinErr = errors.Join(joinErr, err)
		}
		t.closed = true
		return joinErr
	}
	return nil
}

func (t *TempFile) checkClosed() (closed bool, count int) {
	return t.closed, t.count
}

func CreateTempFile(r io.Reader, filename string) (*TempFile, error) {
	inputFile, err := os.CreateTemp("", filename)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(inputFile, r)
	if err != nil {
		return nil, err
	}
	if err := inputFile.Sync(); err != nil {
		return nil, err
	}
	if _, err := inputFile.Seek(0, 0); err != nil {
		return nil, err
	}

	return &TempFile{
		file:   inputFile,
		count:  0,
		closed: false,
	}, nil
}
