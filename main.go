package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
)

type Mode interface {
	Split(io.Reader) error
}

type Line struct {
	UnitLineCount int
}

type Byte struct {
	UnitByteCount int64
}

type Chunk struct {
	UnitChunkCount int
}

func (l Line) Split(reader io.Reader) error {
	scanner := bufio.NewScanner(reader)

	var lineCount int = 0
	fileIndex := 1
	var file *os.File

	for scanner.Scan() {
		if lineCount%l.UnitLineCount == 0 {
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

func (b Byte) Split(reader io.Reader) error {
	// 指定のバイト数が0より小さい場合はエラー
	if b.UnitByteCount <= 0 {
		return fmt.Errorf("invalid unitByteCount value: %d", b.UnitByteCount)
	}

	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanBytes)

	var byteCount int64 = 0
	fileIndex := 1
	var file *os.File

	for scanner.Scan() {
		if byteCount%b.UnitByteCount == 0 {
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

func (c Chunk) Split(reader io.Reader) error {
	// 指定の数が0より小さい場合はエラー
	if c.UnitChunkCount <= 0 {
		return fmt.Errorf("invalid unitChunkCount value: %d", c.UnitChunkCount)
	}

	// 読み込んだ合計のサイズ
	var totalSize int64
	var actualReader io.Reader = reader

	switch r := reader.(type) {
	case *os.File:
		fi, err := r.Stat()
		if err != nil {
			return err
		}
		totalSize = fi.Size()
		actualReader = r
	default:
		buf := new(bytes.Buffer)
		_, err := io.Copy(buf, reader)
		if err != nil {
			return err
		}
		totalSize = int64(buf.Len())
		actualReader = buf
	}

	unitByteCount := totalSize / int64(c.UnitChunkCount)
	extraBytes := totalSize % int64(c.UnitChunkCount) // 追加で余りを取得

	for i := 0; i < c.UnitChunkCount; i++ {
		partFileName := fmt.Sprintf("part-num-%d", i)
		partFile, err := os.Create(partFileName)
		if err != nil {
			return err
		}

		bytesToWrite := unitByteCount
		if i == c.UnitChunkCount-1 {
			bytesToWrite += extraBytes // 最後のファイルに余りの部分を付与
		}

		_, err = io.CopyN(partFile, actualReader, bytesToWrite)
		if err != nil {
			partFile.Close()
			return err
		}
		partFile.Close()

	}
	return nil
}

var (
	lineCount  int
	chunkCount int
	byteCount  int64
)

const defaultLineCount = 1000
const defaultChunkCount = 0
const defaultByteCount = 0

func init() {
	flag.IntVar(&lineCount, "l", defaultLineCount, "line_count [file]")
	flag.IntVar(&chunkCount, "n", defaultChunkCount, "chunk_count [file]")
	flag.Int64Var(&byteCount, "b", defaultByteCount, "byte_count [file]")
}

func main() {
	flag.Parse()
	args := flag.Args()

	// NOTE: lineCount, chunkCount, byteCount の値が、デフォルト値以外の値になっているものが2つ以上あったらエラーにする
	optionCount := 0
	var mode Mode = Line{UnitLineCount: defaultLineCount}
	if lineCount != defaultLineCount {
		optionCount++
		mode = Line{UnitLineCount: lineCount}
	}
	if chunkCount != defaultChunkCount {
		optionCount++
		mode = Chunk{UnitChunkCount: chunkCount}
	}
	if byteCount != defaultByteCount {
		optionCount++
		mode = Byte{UnitByteCount: byteCount}
	}
	if optionCount > 1 {
		fmt.Fprintln(os.Stderr, "Please specify only one option")
		flag.Usage()
		return
	}

	// 分割対象のファイルは１つだけ指定する。 TODO: 標準入力の場合も考慮する
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "invalid args value: %v\n", len(args))
		flag.Usage()
		return
	}

	f, err := os.Open(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "fail open and process file: %v\n", err)
	}
	defer f.Close()

	if err := mode.Split(f); err != nil {
		fmt.Fprintf(os.Stderr, "fail split file: %v\n", err)
	}
}
