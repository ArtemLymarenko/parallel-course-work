package fileReader

import (
	"bufio"
	"errors"
	"log"
	"os"
)

var ErrReadFile = errors.New("failed to read the file")

type fileReader struct{}

func New() *fileReader {
	return &fileReader{}
}

func (r *fileReader) Read(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Println("failed to close the file: ", err)
		}
	}()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	result := make([]byte, fileInfo.Size())
	reader := bufio.NewReader(file)
	chunk := make([]byte, 1024)
	offset := 0
	for {
		n, err := reader.Read(chunk)
		if err != nil && err.Error() != "EOF" {
			return nil, ErrReadFile
		}

		if n == 0 {
			break
		}

		copy(result[offset:], chunk[:n])
		offset += n
	}

	return result, nil
}
