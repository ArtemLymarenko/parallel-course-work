package invertedIdx

import (
	"fmt"
	"parallel-course-work/pkg/mock"
	fileManager "parallel-course-work/server/internal/infrastructure/file_manager"
	"slices"
	"testing"
)

var logs = mock.NewLogger()

func TestParseContent(t *testing.T) {
	reader := fileManager.New(logs)
	invIdx := New("../../../resources/", reader, logs)
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

func TestInvertedIndexBuild(t *testing.T) {
	reader := fileManager.New(logs)
	invIdx := New("../../../resources/", reader, logs)

	if invIdx.storage.GetSize() == 0 {
		t.Errorf("InvertedIndex storage is empty")
	}
}

func TestInvertedIndexSearch(t *testing.T) {
	reader := fileManager.New(logs)
	invIdx := New("../../../resources/", reader, logs)

	files := invIdx.Search("always")
	fmt.Println(files)
	result1 := []string{"0_2.txt", "3_4.txt", "555.txt"}
	for _, w := range files {
		if !slices.Contains(result1, w) {
			t.Errorf("Parsed text does not contain %s", w)
		}
	}

	files = invIdx.Search("chemistry between Kutcher")
	fmt.Println(files)
	result2 := []string{"1_3.txt"}
	for _, w := range files {
		if !slices.Contains(result2, w) {
			t.Errorf("Parsed text does not contain %s", w)
		}
	}
}
