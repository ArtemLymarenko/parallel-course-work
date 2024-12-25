package service

import (
	"fmt"
	"log"
	"time"
)

type InvertedIndex interface {
	AddFile(filePath string) error
	HasFileProcessed(filePath string) bool
}

type FileManager interface {
	GetFilesWithCond(dir string, cond func(fileName string) bool) ([]string, error)
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
		time.Sleep(period)
		files, err := iis.fileManager.GetFilesWithCond(directory, func(filePath string) bool {
			return !iis.invertedIdx.HasFileProcessed(filePath)
		})

		if err != nil {
			iis.logger.Log(err)
			continue
		}

		addedFiles := 0
		for _, filePath := range files {
			err := iis.invertedIdx.AddFile(filePath)
			if err != nil {
				log.Println(err)
				continue
			} else {
				addedFiles++
			}
		}
		msg := fmt.Sprintf("inverted index was updated successfully. added files: %v", addedFiles)
		iis.logger.Log(msg)
	}
}
