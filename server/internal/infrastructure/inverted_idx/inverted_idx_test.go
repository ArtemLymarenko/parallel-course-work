package invertedIdx

import (
	"fmt"
	fileReader "parallel-course-work/server/internal/infrastructure/file_reader"
	"slices"
	"testing"
)

func TestParseContent(t *testing.T) {
	reader := fileReader.NewReader()
	invIdx := New(reader)
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
	reader := fileReader.NewReader()
	invIdx := New(reader)
	fileNames := []string{"0_2.txt", "1_3.txt", "2_3.txt", "3_4.txt"}
	invIdx.Build(fileNames)
}

func TestInvertedIndexSearch(t *testing.T) {
	reader := fileReader.NewReader()
	invIdx := New(reader)
	fileNames := []string{"0_2.txt", "1_3.txt", "2_3.txt", "3_4.txt"}
	invIdx.Build(fileNames)

	files := invIdx.Search("always")
	fmt.Println(files)
	result1 := []string{"0_2.txt", "3_4.txt"}
	for _, w := range files {
		if !slices.Contains(result1, w) {
			t.Errorf("Parsed text does not contain %s", w)
		}
	}

	files = invIdx.Search("Seems like he makes movies specifically")
	fmt.Println(files)
	result2 := []string{"3_4.txt"}
	for _, w := range files {
		if !slices.Contains(result2, w) {
			t.Errorf("Parsed text does not contain %s", w)
		}
	}
}
