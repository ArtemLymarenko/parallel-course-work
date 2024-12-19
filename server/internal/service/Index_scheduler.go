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
	directory   string
	period      time.Duration
	logger      Logger
}

func NewSchedulerService(
	invertedIdx InvertedIndex,
	fileManager FileManager,
	directory string,
	period time.Duration,
	logger Logger,
) *InvertedIndexScheduler {
	return &InvertedIndexScheduler{
		invertedIdx: invertedIdx,
		fileManager: fileManager,
		directory:   directory,
		period:      period,
		logger:      logger,
	}
}

func (iis *InvertedIndexScheduler) ScheduleAsync() {
	for {
		checkPoint := time.Now()
		time.Sleep(iis.period)
		files, err := iis.fileManager.GetFilesAfterTimeStamp(iis.directory, checkPoint)
		if err != nil {
			//log.Println(err)
			iis.logger.Log(err)
			continue
		}

		addedFiles := 0
		for _, filePath := range files {
			if !iis.invertedIdx.HasFileProcessed(filePath) {
				err := iis.invertedIdx.AddFile(iis.directory + filePath)
				if err != nil {
					log.Println(err)
					continue
				} else {
					addedFiles++
				}
			}
		}
		//log.Printf("inverted index was updated successfully at %v\n", checkPoint)
		msg := fmt.Sprintf("inverted index was updated successfully at %v\n", checkPoint)
		iis.logger.Log(msg)
		//log.Printf("successfully added files : %v\n", addedFiles)
		msg = fmt.Sprintf("successfully added files : %v\n", addedFiles)
		iis.logger.Log(msg)
	}
}
