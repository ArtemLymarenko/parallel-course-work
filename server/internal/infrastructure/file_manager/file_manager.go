package fileManager

import (
	"bufio"
	"errors"
	"os"
	"path"
	"path/filepath"
)

var ErrReadFile = errors.New("failed to read the file")

type Logger interface {
	Log(...interface{})
}

type fileManager struct {
	logger Logger
}

func New(logger Logger) *fileManager {
	return &fileManager{
		logger: logger,
	}
}

func (m *fileManager) Read(filePath string) ([]byte, error) {
	file, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = file.Close(); err != nil {
			m.logger.Log("failed to close the file: ", err)
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

func (m *fileManager) GetFilesWithCond(
	dir string,
	cond func(fileName string) bool,
) ([]string, error) {
	directory, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer directory.Close()

	files, err := directory.Readdir(-1)
	if err != nil {
		return nil, err
	}

	var newFiles []string
	for _, file := range files {
		if file.IsDir() {
			subDir := path.Join(dir, file.Name())
			subFiles, err := m.GetFilesWithCond(subDir, cond)
			if err != nil {
				return nil, err
			}
			newFiles = append(newFiles, subFiles...)
		} else {
			filePath := path.Join(dir, file.Name())
			if cond(filePath) {
				newFiles = append(newFiles, filePath)
			}
		}
	}

	return newFiles, nil
}

func (m *fileManager) GetAllFiles(dir string) ([]string, error) {
	directory, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer directory.Close()

	files, err := directory.Readdir(-1)
	if err != nil {
		return nil, err
	}

	var newFiles []string
	for _, file := range files {
		if file.IsDir() {
			subDir := path.Join(dir, file.Name())
			subFiles, err := m.GetAllFiles(subDir)
			if err != nil {
				return nil, err
			}
			newFiles = append(newFiles, subFiles...)
		} else {
			newFiles = append(newFiles, path.Join(dir, file.Name()))
		}
	}

	return newFiles, nil
}
