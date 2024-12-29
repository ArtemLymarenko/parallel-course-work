package invertedIdx

import (
	"errors"
	"github.com/ArtemLymarenko/parallel-course-work/pkg/set"
	"log"
	"regexp"
	"strings"
	"sync"
)

type FileManager interface {
	Read(filePath string) ([]byte, error)
	GetAllFiles(dir string) ([]string, error)
}

type Logger interface {
	Log(...interface{})
}

type InvertedIndex struct {
	storage        *SyncHashMap
	processedFiles *set.Set[string]
	processedLock  sync.RWMutex
	commonWords    *set.Set[string]
	fileManager    FileManager
	logger         Logger
}

func New(fileManager FileManager, logger Logger) *InvertedIndex {
	commonWords := set.NewSet[string]()
	words := []string{
		"the", "be", "to", "of", "and", "a", "in", "that", "have", "i",
		"it", "for", "not", "on", "with", "he", "as", "you", "do", "at",
		"this", "but", "his", "by", "from", "they", "we", "say", "her", "she",
		"or", "an", "will", "my", "one", "all", "would", "there", "their", "what",
		"so", "up", "out", "if", "about", "who", "get", "which", "go", "me",
		"now", "him", "is", "are", "was", "were", "its",
	}

	for _, word := range words {
		commonWords.Add(word)
	}

	invIndex := &InvertedIndex{
		storage:        NewSyncHashMap(32),
		processedFiles: set.NewSet[string](),
		commonWords:    commonWords,
		fileManager:    fileManager,
		processedLock:  sync.RWMutex{},
		logger:         logger,
	}

	return invIndex
}

func (i *InvertedIndex) parseText(content string) []string {
	text := strings.TrimSpace(strings.ToLower(content))

	reg, _ := regexp.Compile(`[,.!?(){}\[\]"]`)
	text = reg.ReplaceAllString(text, "")

	regBr, _ := regexp.Compile(`(?i)<br\s*/?>|'s|'ve|-|'re|n't|'d|'ll|'`)
	text = regBr.ReplaceAllString(text, " ")

	splitText := strings.Fields(text)

	uniqueWords := *set.NewSet[string]()
	for _, word := range splitText {
		if !i.commonWords.Has(word) {
			uniqueWords.Add(word)
		}
	}

	return uniqueWords.ToSlice()
}

func (i *InvertedIndex) Build(resourceDir string, threadCount int) {
	if threadCount < 1 {
		log.Fatalf("Thread count must be greater than zero")
	}

	filePaths, err := i.fileManager.GetAllFiles(resourceDir)
	if err != nil {
		log.Fatalf("could not read the directory: %v", err)
	}

	wg := sync.WaitGroup{}
	wg.Add(threadCount)

	totalChunks := len(filePaths) / threadCount
	startIdx, endIdx := 0, 0
	for threadIdx := range threadCount {
		startIdx = threadIdx * totalChunks
		if threadIdx == threadCount-1 {
			endIdx = len(filePaths)
		} else {
			endIdx = totalChunks * (threadIdx + 1)
		}

		filePathsChunk := filePaths[startIdx:endIdx]
		go func() {
			i.BuildFiles(filePathsChunk)
			wg.Done()
		}()
	}

	wg.Wait()
}

func (i *InvertedIndex) BuildFiles(filePaths []string) {
	idx := 0
	processedFile := make([]string, len(filePaths))
	for _, filePath := range filePaths {
		if err := i.AddFile(filePath); err != nil {
			i.logger.Log(err)
			continue
		}
		processedFile[idx] = filePath
		idx++
	}

	i.setProcessedFiles(processedFile[:idx])
}

func (i *InvertedIndex) AddFile(filePath string) error {
	if i.HasFileProcessed(filePath) {
		return errors.New("file has already been added to index")
	}

	fileContent, err := i.fileManager.Read(filePath)
	if err != nil {
		return err
	}

	parsedFileContent := i.parseText(string(fileContent))
	for _, word := range parsedFileContent {
		i.storage.Put(word, filePath)
	}

	i.processedLock.Lock()
	defer i.processedLock.Unlock()
	i.processedFiles.Add(filePath)
	return nil
}

func (i *InvertedIndex) HasFileProcessed(filePath string) bool {
	i.processedLock.Lock()
	defer i.processedLock.Unlock()
	return i.processedFiles.Has(filePath)
}

func (i *InvertedIndex) setProcessedFiles(files []string) {
	i.processedLock.Lock()
	defer i.processedLock.Unlock()
	for _, file := range files {
		i.processedFiles.Add(file)
	}
}

func (i *InvertedIndex) GetFileContent(filePath string) ([]byte, error) {
	fileContent, err := i.fileManager.Read(filePath)
	if err != nil {
		return nil, err
	}

	return fileContent, nil
}

func (i *InvertedIndex) RemoveFile(filePath string) error {
	if !i.HasFileProcessed(filePath) {
		return errors.New("nothing to remove")
	}

	fileContent, err := i.fileManager.Read(filePath)
	if err != nil {
		return err
	}

	parsedFileContent := i.parseText(string(fileContent))
	for _, word := range parsedFileContent {
		i.storage.Remove(word, filePath)
	}

	i.processedLock.Lock()
	defer i.processedLock.Unlock()
	i.processedFiles.Remove(filePath)

	return nil
}

func (i *InvertedIndex) Search(query string) []string {
	parsed := i.parseText(query)

	filesCount := make(map[string]int)
	maxCount := -1

	for _, word := range parsed {
		if fileSet, exists := i.storage.Get(word); exists {
			for fileName := range fileSet.Keys {
				filesCount[fileName] += 1
				if filesCount[fileName] > maxCount {
					maxCount = filesCount[fileName]
				}
			}
		}
	}

	var result []string
	for fileName, count := range filesCount {
		if count == maxCount {
			result = append(result, fileName)
		}
	}

	return result
}
