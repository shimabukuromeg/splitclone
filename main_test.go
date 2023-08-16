package main

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func countFileLines(t *testing.T, file *os.File) int {
	t.Helper()

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
func TestSplit(t *testing.T) {
	splitTests := []struct {
		name      string
		input     string
		mode      Mode
		wantFiles []struct {
			name      string
			byteCount int64
			lineCount int
		}
	}{
		{
			name:  "LineSplit",
			input: "line1\nline2\nline3\n",
			mode:  Line{UnitLineCount: 2},
			wantFiles: []struct {
				name      string
				byteCount int64
				lineCount int
			}{
				{name: "part-line-1", byteCount: 12, lineCount: 2},
				{name: "part-line-2", byteCount: 6, lineCount: 1},
			},
		},
		{
			name:  "ByteSplit",
			input: "abcdefghijklmn",
			mode:  Byte{UnitByteCount: 4},
			wantFiles: []struct {
				name      string
				byteCount int64
				lineCount int
			}{
				{name: "part-byte-1", byteCount: 4, lineCount: 1},
				{name: "part-byte-2", byteCount: 4, lineCount: 1},
				{name: "part-byte-3", byteCount: 4, lineCount: 1},
				{name: "part-byte-4", byteCount: 2, lineCount: 1},
			},
		},
		{
			name:  "ChunkSplit",
			input: "abcdefghijklmn",
			mode:  Chunk{UnitChunkCount: 2},
			wantFiles: []struct {
				name      string
				byteCount int64
				lineCount int
			}{
				{name: "part-num-1", byteCount: 7, lineCount: 1},
				{name: "part-num-2", byteCount: 7, lineCount: 1},
			},
		},
	}

	for _, tt := range splitTests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bytes.NewBufferString(tt.input)
			d := t.TempDir()

			if err := tt.mode.Split(reader, d); err != nil {
				t.Errorf("Failed to split: %v", err)
			}

			checkFiles(t, d, tt.wantFiles)
		})
	}
}

func BenchmarkSplit(b *testing.B) {
	splitTests := []struct {
		name      string
		input     string
		mode      Mode
		wantFiles []struct {
			name      string
			byteCount int64
			lineCount int
		}
	}{
		{
			name:  "LineSplit",
			input: "line1\nline2\nline3\n",
			mode:  Line{UnitLineCount: 2},
			wantFiles: []struct {
				name      string
				byteCount int64
				lineCount int
			}{
				{name: "part-line-1", byteCount: 12, lineCount: 2},
				{name: "part-line-2", byteCount: 6, lineCount: 1},
			},
		},
		{
			name:  "ByteSplit",
			input: "abcdefghijklmn",
			mode:  Byte{UnitByteCount: 4},
			wantFiles: []struct {
				name      string
				byteCount int64
				lineCount int
			}{
				{name: "part-byte-1", byteCount: 4, lineCount: 1},
				{name: "part-byte-2", byteCount: 4, lineCount: 1},
				{name: "part-byte-3", byteCount: 4, lineCount: 1},
				{name: "part-byte-4", byteCount: 2, lineCount: 1},
			},
		},
		{
			name:  "ChunkSplit",
			input: "abcdefghijklmn",
			mode:  Chunk{UnitChunkCount: 2},
			wantFiles: []struct {
				name      string
				byteCount int64
				lineCount int
			}{
				{name: "part-num-1", byteCount: 7, lineCount: 1},
				{name: "part-num-2", byteCount: 7, lineCount: 1},
			},
		},
	}

	for _, tt := range splitTests {
		b.Run(tt.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ { // ベンチマークのループ
				reader := bytes.NewBufferString(tt.input)
				d := b.TempDir()

				if err := tt.mode.Split(reader, d); err != nil {
					b.Errorf("Failed to split: %v", err)
				}
			}
		})
	}
}

// ファイルの内容をチェックする共通関数
func checkFiles(t *testing.T, dir string, wantFiles []struct {
	name      string
	byteCount int64
	lineCount int
}) {
	t.Helper()

	files, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("Failed to read dir: %v", err)
	}

	if len(files) != len(wantFiles) {
		t.Errorf("mismatch the number of files, got=%d, want=%d", len(files), len(wantFiles))
	}

	for i, f := range files {
		if f.Name() != wantFiles[i].name {
			t.Errorf("file name mismatch at index %d: got=%s, want=%s", i, f.Name(), wantFiles[i].name)
		}

		file, err := os.Open(filepath.Join(dir, f.Name()))
		if err != nil {
			t.Fatalf("Failed to open file: %v", err)
		}
		defer file.Close()

		fi, err := file.Stat()
		if err != nil {
			t.Fatalf("Failed to get file stats: %v", err)
		}

		if fi.Size() != wantFiles[i].byteCount {
			t.Errorf("file size mismatch at index %d: got=%d, want=%d", i, fi.Size(), wantFiles[i].byteCount)
		}

		lineCount := countFileLines(t, file)
		if lineCount != wantFiles[i].lineCount {
			t.Errorf("files[%d], %d != %d", i, lineCount, wantFiles[i].lineCount)
		}
	}
}
