package invertedIdx

import (
	"fmt"
	"parallel-course-work/pkg/set"
	syncMap "parallel-course-work/pkg/sync_map"
	"regexp"
	"strings"
)

type SyncMap[V any] interface {
	Put(string, V)
	Get(string) (V, bool)
	Remove(string)
	GetSize() int64
}

type FileReader interface {
	Read(filePath string) ([]byte, error)
}

type InvertedIndex struct {
	//storage     map[string]*set.Set[string]
	//lock        sync.RWMutex
	storage     SyncMap[*set.Set[string]]
	commonWords *set.Set[string]
	fileReader  FileReader
}

func New(resourceDir string, fileReader FileReader) *InvertedIndex {
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
		storage:     syncMap.NewSyncHashMap[*set.Set[string]](4, 32),
		commonWords: commonWords,
		fileReader:  fileReader,
	}

	invIndex.Build(resourceDir, []string{"0_2.txt", "1_3.txt", "2_3.txt", "3_4.txt"})

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

func (i *InvertedIndex) Build(resourceDir string, fileNames []string) {
	for _, fileName := range fileNames {
		if err := i.AddFile(resourceDir, fileName); err != nil {
			fmt.Println(err)
			continue
		}
	}
}

func (i *InvertedIndex) AddFile(resourcesDir, fileName string) error {
	fileContent, err := i.fileReader.Read(resourcesDir + fileName)
	if err != nil {
		return err
	}

	parsedFileContent := i.parseText(string(fileContent))
	for _, word := range parsedFileContent {
		fileSet, exists := i.storage.Get(word)
		if exists {
			fileSet.Add(fileName)
		} else {
			newSet := set.NewSet[string]()
			newSet.Add(fileName)
			i.storage.Put(word, newSet)
		}
	}

	return nil
}

func (i *InvertedIndex) GetFileContent(resourcesDir, fileName string) ([]byte, error) {
	fileContent, err := i.fileReader.Read(resourcesDir + fileName)
	if err != nil {
		return nil, err
	}

	return fileContent, nil
}

func (i *InvertedIndex) RemoveFile(resourcesDir, fileName string) error {
	fileContent, err := i.fileReader.Read(resourcesDir + fileName)
	if err != nil {
		return err
	}

	parsedFileContent := i.parseText(string(fileContent))
	for _, word := range parsedFileContent {
		if fileSet, exists := i.storage.Get(word); exists {
			fileSet.Remove(fileName)
			if fileSet.IsEmpty() {
				i.storage.Remove(word)
			}
		}
	}

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
