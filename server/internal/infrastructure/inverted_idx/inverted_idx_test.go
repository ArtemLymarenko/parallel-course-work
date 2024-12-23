package invertedIdx

import (
	"fmt"
	"parallel-course-work/pkg/mock"
	fileManager "parallel-course-work/server/internal/infrastructure/file_manager"
	"slices"
	"sync"
	"testing"
	"time"
)

var logs = mock.NewLogger()

func TestParseContent(t *testing.T) {
	reader := fileManager.New(logs)
	invIdx := New(reader, logs)
	invIdx.Build("test_files/")
	text := `Hello,<br /><br />World,  World,  World, iast's a<br /> beautiful "chopper! " "?,,,.. . right now the `
	parseText := invIdx.parseText(text)

	result := []string{"iast", "beautiful", "chopper", "right", "hello", "world"}
	fmt.Println(parseText)
	for _, w := range parseText {
		if !slices.Contains(result, w) {
			t.Errorf("Parsed text does not contain %s", w)
		}
	}
}

func TestParseContent2(t *testing.T) {
	reader := fileManager.New(logs)
	invIdx := New(reader, logs)
	invIdx.Build("test_files/")
	text := `Another text`
	search := invIdx.Search(text)
	fmt.Println(search)
}

func TestAddFile(t *testing.T) {
	reader := fileManager.New(logs)
	invIdx := New(reader, logs)
	err := invIdx.AddFile("C:/Users/Artem/Desktop/MyFile.txt")
	if err != nil {
		t.Error(err)
	}
	search := invIdx.Search("Привіт")
	fmt.Println(search)
}

func TestInvertedIndexBuild(t *testing.T) {
	reader := fileManager.New(logs)
	invIdx := New(reader, logs)
	invIdx.Build("test_files/")
	if invIdx.storage.GetSize() == 0 {
		t.Errorf("InvertedIndex storage is empty")
	}
}

func TestInvertedIndexSearch(t *testing.T) {
	reader := fileManager.New(logs)
	invIdx := New(reader, logs)
	const pref = "test_files/"
	invIdx.Build(pref)
	files := invIdx.Search("always")
	fmt.Println(files)
	result1 := []string{pref + "0_2.txt", pref + "3_4.txt"}
	for _, w := range files {
		if !slices.Contains(result1, w) {
			t.Errorf("Parsed text does not contain %s", w)
		}
	}

	files = invIdx.Search("chemistry between Kutcher")
	fmt.Println(files)
	result2 := []string{pref + "0_2.txt", pref + "2_3.txt"}
	for _, w := range files {
		if !slices.Contains(result2, w) {
			t.Errorf("Parsed text does not contain %s", w)
		}
	}
}

func TestInvertedIndexAddFile(t *testing.T) {
	reader := fileManager.New(logs)
	invIdx := New(reader, logs)

	files := []string{"test_files/0_2.txt", "test_files/2_3.txt", "test_files/3_4.txt"}
	wg := sync.WaitGroup{}
	wg.Add(2 * len(files))
	for _, f := range files {
		go func() {
			err := invIdx.AddFile(f)
			if err != nil {
				t.Error(err)
			}
			wg.Done()
		}()

	}
	time.Sleep(1 * time.Second)
	fmt.Println("added")
	for _, f := range files {
		go func() {
			err := invIdx.RemoveFile(f)
			if err != nil {
				t.Error(err)
			}
			wg.Done()
		}()
	}

	wg.Wait()

	get, ok := invIdx.storage.Get("once")
	fmt.Println(get, ok)
	get, ok = invIdx.storage.Get("seems")
	fmt.Println(get, ok)
	get, ok = invIdx.storage.Get("kevin")
	fmt.Println(get, ok)
}
