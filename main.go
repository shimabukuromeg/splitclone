package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func OpenAndProcessFile(fileName string, lineCount int) error {
	f, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("failed open file: %w", err)
	}
	defer f.Close()

	if err := SplitTextIntoFiles(f, lineCount); err != nil {
		return err
	}

	return nil
}

func SplitTextIntoFiles(reader io.Reader, lineCount int) error {
	scanner := bufio.NewScanner(reader)
	count := 0
	fileIndex := 1
	var file *os.File

	for scanner.Scan() {
		if count%lineCount == 0 {
			if file != nil {
				file.Close()
				file = nil
			}

			// 書き出しファイル
			filename := fmt.Sprintf("example%d", fileIndex)
			var err error
			file, err = os.Create(filename)
			if err != nil {
				panic(err)
			}
			fileIndex++

		}
		fmt.Fprintln(file, scanner.Text())

		count++
	}

	if file != nil {
		file.Close()
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("reading input: %w", err)
	}

	return nil
}

var (
	lineCount  int
	chunkCount int
	byteCount  int
)

func init() {
	flag.IntVar(&lineCount, "l", 1000, "line_count [file]")
	flag.IntVar(&chunkCount, "n", 1000, "chunk_count [file]")
	flag.IntVar(&byteCount, "n", 1000, "byte_count [file]")
}

func main() {
	flag.Parse()
	args := flag.Args()

	// 分割対象のファイルは１つだけ指定
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "invalid args value: %v\n", len(args))
		flag.Usage()
	}

	err := OpenAndProcessFile(args[0], lineCount)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fail open and process file: %v\n", err)
	}

}
