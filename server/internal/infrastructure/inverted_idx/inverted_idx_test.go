package invertedIdx

import (
	"encoding/csv"
	"fmt"
	"github.com/ArtemLymarenko/parallel-course-work/pkg/mock"
	"os"
	fileManager "server/internal/infrastructure/file_manager"
	"slices"
	"sync"
	"testing"
	"time"
)

var logs = mock.NewLogger()

func TestParseContent(t *testing.T) {
	reader := fileManager.New(logs)
	invIdx := New(reader, logs)
	invIdx.Build("test_files/", 1)
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
	invIdx.Build("test_files/", 1)
	text := `Another text`
	search := invIdx.Search(text)
	fmt.Println(search)
}

func TestInvertedIndexBuild(t *testing.T) {
	reader := fileManager.New(logs)
	invIdx := New(reader, logs)
	invIdx.Build("test_files/", 1)
	if invIdx.storage.GetSize() == 0 {
		t.Errorf("InvertedIndex storage is empty")
	}
}

func TestInvertedIndexSearch(t *testing.T) {
	reader := fileManager.New(logs)
	invIdx := New(reader, logs)
	const pref = "test_files/"
	invIdx.Build(pref, 1)
	files := invIdx.Search("always")
	fmt.Println(files)
	result1 := []string{
		pref + "0_2.txt", pref + "2_3.txt", pref + "3_4.txt",
		pref + "dir1/0_2.txt", pref + "dir1/2_3.txt", pref + "dir1/3_4.txt",
		pref + "dir2/0_2.txt", pref + "dir2/2_3.txt", pref + "dir2/3_4.txt",
	}
	for _, w := range files {
		if !slices.Contains(result1, w) {
			t.Errorf("Parsed text does not contain %s", w)
		}
	}

	files = invIdx.Search("chemistry between Kutcher")
	fmt.Println(files)
	result2 := []string{
		pref + "0_2.txt", pref + "2_3.txt", pref + "3_4.txt",
		pref + "dir1/0_2.txt", pref + "dir1/2_3.txt", pref + "dir1/3_4.txt",
		pref + "dir2/0_2.txt", pref + "dir2/2_3.txt", pref + "dir2/3_4.txt",
	}
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
	if ok {
		t.Errorf("value should not be found in storage %v", get)
	}
	get, ok = invIdx.storage.Get("seems")
	if ok {
		t.Errorf("value should not be found in storage %v", get)
	}
	get, ok = invIdx.storage.Get("kevin")
	if ok {
		t.Errorf("value should not be found in storage %v", get)
	}
}

func TestInvertedIndex_Build(t *testing.T) {
	reader := fileManager.New(logs)
	invIdx := New(reader, logs)

	invIdx.Build("test_files/", 12)
	expectSize := int64(236)
	if invIdx.storage.GetSize() != expectSize {
		t.Errorf("InvertedIndex storage is wrong size, expected %v", expectSize)
	}
}

func Test_Build_Multi(t *testing.T) {
	resourceDir := "../../../resources/data/"

	file, err := os.Create("test_results.csv")
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{"threadCount", "ms"})
	if err != nil {
		t.Fatalf("Failed to write header to CSV: %v", err)
	}

	runTest := func(threadCount int) {
		invIdx := New(fileManager.New(logs), logs)

		start := time.Now()
		invIdx.Build(resourceDir, threadCount)
		elapsed := time.Since(start)

		err := writer.Write([]string{
			fmt.Sprintf("%v", threadCount),
			fmt.Sprintf("%v", elapsed.Milliseconds()),
		})
		if err != nil {
			t.Fatalf("Failed to write result to CSV: %v", err)
		}

		t.Logf("Elapsed for %v threads is %v ms\n", threadCount, elapsed.Milliseconds())
	}

	threadCount := 1
	for threadCount < 1024 {
		t.Run(fmt.Sprintf("threads_%v", threadCount), func(t *testing.T) {
			t.Cleanup(func() {})
			runTest(threadCount)
			threadCount *= 2
		})
	}
}
