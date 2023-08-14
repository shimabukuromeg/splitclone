package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

type Line struct {
	LineCount int
}

type Byte struct {
	BytesPerPart int64
}

func SplitFileByLineCount(reader io.Reader, unitLineCount int) error {
	scanner := bufio.NewScanner(reader)

	var lineCount int = 0
	fileIndex := 1
	var file *os.File

	for scanner.Scan() {
		if lineCount%unitLineCount == 0 {
			if file != nil {
				file.Close()
				file = nil
			}

			// 書き出しファイル
			filename := fmt.Sprintf("part-line-%d", fileIndex)
			var err error
			file, err = os.Create(filename)
			if err != nil {
				panic(err)
			}
			fileIndex++

		}
		fmt.Fprintln(file, scanner.Text())

		lineCount++
	}

	if file != nil {
		file.Close()
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("reading input: %w", err)
	}

	return nil
}

func SplitFileByBytes(reader io.Reader, unitByteCount int64) error {
	// 指定のバイト数が0より小さい場合はエラー
	if unitByteCount <= 0 {
		return fmt.Errorf("invalid unitByteCount value: %d", unitByteCount)
	}

	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanBytes)

	var byteCount int64 = 0
	fileIndex := 1
	var file *os.File

	for scanner.Scan() {
		if byteCount%unitByteCount == 0 {
			if file != nil {
				file.Close()
				file = nil
			}

			// 書き出しファイル
			filename := fmt.Sprintf("part-byte-%d", fileIndex)
			var err error
			file, err = os.Create(filename)
			if err != nil {
				panic(err)
			}
			fileIndex++

		}
		fmt.Fprint(file, scanner.Text())

		byteCount++
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
	unitLineCount int
	chunkCount    int
	byteCount     int64
)

func init() {
	flag.IntVar(&unitLineCount, "l", 1000, "line_count [file]")
	flag.IntVar(&chunkCount, "n", 0, "chunk_count [file]")
	flag.Int64Var(&byteCount, "b", 0, "byte_count [file]")
}

func main() {
	flag.Parse()
	args := flag.Args()

	fmt.Println(args)
	fmt.Println(chunkCount)
	fmt.Println(byteCount)

	// 分割対象のファイルは１つだけ指定
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "invalid args value: %v\n", len(args))
		flag.Usage()
		return
	}

	// TODO: lineCount, chunkCount, byteCount の値が、デフォルト値以外の値になっているものが2つ以上あったらエラーにする

	f, err := os.Open(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "fail open and process file: %v\n", err)
	}
	defer f.Close()

	// TODO: 標準入力を読み取ります
	// reader := bufio.NewReader(os.Stdin)

	// 指定したバイト数で分割
	if err := SplitFileByBytes(f, byteCount); err != nil {
		fmt.Fprintf(os.Stderr, "fail split file by bytes: %v\n", err)
	}

	// 指定した行数で分割
	if err := SplitFileByLineCount(f, unitLineCount); err != nil {
		fmt.Fprintf(os.Stderr, "fail open and process file: %v\n", err)
	}

	// TODO: 指定した数で分割

	if err != nil {
		fmt.Fprintf(os.Stderr, "fail open and process file: %v\n", err)
	}

}
