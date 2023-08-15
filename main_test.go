package main

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func countFileLines(t *testing.T, filename string) int {
	t.Helper() // この関数がヘルパーであることを示す

	file, err := os.Open(filename)
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("Error reading file: %v", err)
	}

	return lineCount
}

func TestLineSplit(t *testing.T) {
	input := "line1\nline2\nline3\n"
	reader := bytes.NewBufferString(input)
	mode := Line{UnitLineCount: 2}

	d := t.TempDir()

	// 分割実行
	if err := mode.Split(reader, d); err != nil {
		t.Errorf("Failed to split: %v", err)
	}

	// 生成されたファイル
	files, err := os.ReadDir(d)
	if err != nil {
		t.Errorf("!!! %+v", err)
	}

	// ファイルの期待値
	wantFiles := []struct {
		name      string
		lineCount int
	}{
		{
			name:      "part-line-1",
			lineCount: 2,
		},
		{
			name:      "part-line-2",
			lineCount: 1,
		},
	}

	// ファイルの件数の確認
	if len(files) != len(wantFiles) {
		t.Errorf("mismatch the number of files, got=%d, want=%d", len(files), len(wantFiles))
	}

	// ファイル名確認
	for i, f := range files {
		if f.Name() != wantFiles[i].name {
			t.Errorf("files[%d], %s != %s", i, f.Name(), wantFiles[i].name)
		}
	}

	// ファイルの行数を確認する
	for i, f := range files {
		if countFileLines(t, filepath.Join(d, f.Name())) != wantFiles[i].lineCount {
			t.Errorf("files[%d], %d != %d", i, countFileLines(t, filepath.Join(d, f.Name())), wantFiles[i].lineCount)
		}
	}

}

// func TestByteSplit(t *testing.T) {
// 	input := "abc"
// 	reader := bytes.NewBufferString(input)
// 	mode := Byte{UnitByteCount: 1}

// 	if err := mode.Split(reader); err != nil {
// 		t.Errorf("Failed to split: %v", err)
// 	}
// 	// 同様に、ファイルの出力を確認またはモックします。
// }

// func TestChunkSplit(t *testing.T) {
// 	input := "abcdef"
// 	reader := bytes.NewBufferString(input)
// 	mode := Chunk{UnitChunkCount: 2}

// 	if err := mode.Split(reader); err != nil {
// 		t.Errorf("Failed to split: %v", err)
// 	}
// 	// 同様に、ファイルの出力を確認またはモックします。
// }
