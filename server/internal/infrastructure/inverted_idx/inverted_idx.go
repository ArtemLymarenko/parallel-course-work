package invertedIdx

import (
	"log"
	"os"
	"parallel-course-work/pkg/set"
	syncMap "parallel-course-work/pkg/sync_map"
	"regexp"
	"strings"
	"sync"
)

type SyncMap[V any] interface {
	Put(string, V)
	Get(string) (V, bool)
	Remove(string)
	GetSize() int64
	Modify(key string, cb func(modify V) interface{}) (bool, interface{})
}

type FileManager interface {
	Read(filePath string) ([]byte, error)
}

type Logger interface {
	Log(...interface{})
}

type InvertedIndex struct {
	storage        SyncMap[*set.Set[string]]
	processedFiles *set.Set[string]
	processedLock  *sync.RWMutex
	commonWords    *set.Set[string]
	fileManager    FileManager
	logger         Logger
}

func New(resourceDir string, fileManager FileManager, logger Logger) *InvertedIndex {
	commonWords := set.NewSet[string]()
	words := []string{
		"the", "be", "to", "of", "and", "a", "in", "that", "have", "i",
		"it", "for", "not", "on", "with", "he", "as", "you", "do", "at",
		"this", "but", "his", "by", "from", "they", "we", "say", "her", "she",
		"or", "an", "will", "my", "one", "all", "would", "there", "their", "what",
		"so", "up", "out", "if", "about", "who", "get", "which", "go", "me", "now", "him", "is", "are",
		"was", "were", "its",
	}

	for _, word := range words {
		commonWords.Add(word)
	}

	invIndex := &InvertedIndex{
		storage:        syncMap.NewSyncHashMap[*set.Set[string]](32, 16),
		processedFiles: set.NewSet[string](),
		commonWords:    commonWords,
		fileManager:    fileManager,
		processedLock:  &sync.RWMutex{},
		logger:         logger,
	}

	invIndex.Build(resourceDir)

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

func (i *InvertedIndex) Build(resourceDir string) {
	files, err := os.ReadDir(resourceDir)
	if err != nil {
		log.Fatalf("could not read the directory: %v", err)
	}

	idx := 0
	processedFile := make([]string, len(files))
	for _, file := range files {
		if err = i.AddFile(resourceDir + file.Name()); err != nil {
			i.logger.Log(err)
			continue
		}
		processedFile[idx] = file.Name()
		idx++
	}

	i.setProcessedFiles(processedFile[:idx])
}

func (i *InvertedIndex) AddFile(filePath string) error {
	fileContent, err := i.fileManager.Read(filePath)
	if err != nil {
		return err
	}

	parsedFileContent := i.parseText(string(fileContent))
	for _, word := range parsedFileContent {
		modified, _ := i.storage.Modify(word, func(fileSet *set.Set[string]) interface{} {
			fileSet.Add(filePath)
			return nil
		})
		if !modified {
			newSet := set.NewSet[string]()
			newSet.Add(filePath)
			i.storage.Put(word, newSet)
		}
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
	fileContent, err := i.fileManager.Read(filePath)
	if err != nil {
		return err
	}

	parsedFileContent := i.parseText(string(fileContent))
	for _, word := range parsedFileContent {
		modified, isSetEmpty := i.storage.Modify(word, func(fileSet *set.Set[string]) interface{} {
			fileSet.Remove(filePath)
			return fileSet.IsEmpty()
		})

		if modified {
			if isEmpty, ok := isSetEmpty.(bool); ok && isEmpty {
				i.storage.Remove(word)
			}
		}
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

	result := make([]string, len(filesCount))
	idx := 0
	for fileName, count := range filesCount {
		if count == maxCount {
			result[idx] = fileName
			idx++
		}
	}

	return result[:idx]
}
