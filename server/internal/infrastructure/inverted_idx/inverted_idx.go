package invertedIdx

import (
	"parallel-course-work/pkg/set"
	"regexp"
	"strings"
	"sync"
)

const resourcesDir = "D:/Development/GOLang/parallel-course-work/server/resources"

type FileReader interface {
	Read(dir, fileName string) ([]byte, error)
}

type InvertedIndex struct {
	storage     map[string]*set.Set[string]
	lock        sync.RWMutex
	commonWords *set.Set[string]
	fileReader  FileReader
}

func New(fileReader FileReader) *InvertedIndex {
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

	return &InvertedIndex{
		storage:     make(map[string]*set.Set[string]),
		lock:        sync.RWMutex{},
		commonWords: commonWords,
		fileReader:  fileReader,
	}
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

func (i *InvertedIndex) Build(fileNames []string) {
	for _, fileName := range fileNames {
		if err := i.AddFile(fileName); err != nil {
			continue
		}
	}
}

func (i *InvertedIndex) AddFile(fileName string) error {
	fileContent, err := i.fileReader.Read(resourcesDir, fileName)
	if err != nil {
		return err
	}

	parsedFileContent := i.parseText(string(fileContent))
	i.lock.Lock()
	for _, word := range parsedFileContent {
		if _, exists := i.storage[word]; !exists {
			i.storage[word] = set.NewSet[string]()
		}
		i.storage[word].Add(fileName)
	}
	i.lock.Unlock()

	return nil
}

func (i *InvertedIndex) Search(query string) []string {
	parsed := i.parseText(query)

	filesCount := make(map[string]int)
	maxCount := -1

	i.lock.Lock()
	for _, word := range parsed {
		if files, exists := i.storage[word]; exists {
			for fileName := range files.Keys {
				filesCount[fileName] += 1
				if filesCount[fileName] > maxCount {
					maxCount = filesCount[fileName]
				}
			}
		}
	}
	i.lock.Unlock()

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
