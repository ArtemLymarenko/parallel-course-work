package service

import (
	"fmt"
	"log"
	"time"
)

type InvertedIndex interface {
	Build(dir string)
	AddFile(filePath string) error
	HasFileProcessed(filePath string) bool
}

type FileManager interface {
	GetFilesAfterTimeStamp(dir string, after time.Time) ([]string, error)
}

type Logger interface {
	Log(...interface{})
}

type InvertedIndexScheduler struct {
	invertedIdx InvertedIndex
	fileManager FileManager
	logger      Logger
}

func NewSchedulerService(
	invertedIdx InvertedIndex,
	fileManager FileManager,
	logger Logger,
) *InvertedIndexScheduler {
	return &InvertedIndexScheduler{
		invertedIdx: invertedIdx,
		fileManager: fileManager,
		logger:      logger,
	}
}

func (iis *InvertedIndexScheduler) MonitorDirAsync(directory string, period time.Duration) {
	for {
		checkPoint := time.Now()
		time.Sleep(period)
		files, err := iis.fileManager.GetFilesAfterTimeStamp(directory, checkPoint)
		if err != nil {
			iis.logger.Log(err)
			continue
		}

		addedFiles := 0
		for _, filePath := range files {
			if !iis.invertedIdx.HasFileProcessed(filePath) {
				err := iis.invertedIdx.AddFile(directory + filePath)
				if err != nil {
					log.Println(err)
					continue
				} else {
					addedFiles++
				}
			}
		}
		msg := fmt.Sprintf("inverted index was updated successfully. added files: %v", addedFiles)
		iis.logger.Log(msg)
	}
}
