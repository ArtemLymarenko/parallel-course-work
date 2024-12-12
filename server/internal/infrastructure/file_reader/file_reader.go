package fileReader

import (
	"bufio"
	"errors"
	"os"
)

type fileReader struct{}

func NewReader() *fileReader {
	return &fileReader{}
}

func (r *fileReader) Read(dir, fileName string) ([]byte, error) {
	file, err := os.Open(dir + "/" + fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	result := make([]byte, fileInfo.Size())
	fileReader := bufio.NewReader(file)
	chunk := make([]byte, 1024)
	offset := 0
	for {
		n, err := fileReader.Read(chunk)
		if err != nil && err.Error() != "EOF" {
			return nil, errors.New("failed to read file")
		}

		if n == 0 {
			break
		}

		copy(result[offset:], chunk[:n])
		offset += n
	}

	return result, nil
}
