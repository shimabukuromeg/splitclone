package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Mode interface {
	Split(io.Reader, string) error
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

func (l Line) Split(reader io.Reader, outputDirName string) error {
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
			file, err = os.Create(filepath.Join(outputDirName, filename))
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

func (b Byte) Split(reader io.Reader, outputDirName string) error {
	// 指定のバイト数が0より小さい場合はエラー
	if b.UnitByteCount <= 0 {
		return fmt.Errorf("invalid unitByteCount value: %d", b.UnitByteCount)
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

	chunkCount = int(totalSize) / int(b.UnitByteCount)
	if int(totalSize)%int(b.UnitByteCount) != 0 {
		chunkCount++
	}

	i := 1
	value := totalSize

	for value > 0 {
		bytesToWrite := b.UnitByteCount
		if value < b.UnitByteCount {
			bytesToWrite = value
		}
		filename := fmt.Sprintf("part-byte-%d", i)
		file, err := os.Create(filepath.Join(outputDirName, filename))
		if err != nil {
			return err
		}
		_, err = io.CopyN(file, actualReader, bytesToWrite)
		if err != nil {
			file.Close()
			return err
		}
		file.Close()
		value -= b.UnitByteCount
		i++
	}

	return nil
}

func (c Chunk) Split(reader io.Reader, outputDirName string) error {
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
		filename := fmt.Sprintf("part-num-%d", i+1)
		file, err := os.Create(filepath.Join(outputDirName, filename))
		if err != nil {
			return err
		}

		bytesToWrite := unitByteCount
		if i == c.UnitChunkCount-1 {
			bytesToWrite += extraBytes // 最後のファイルに余りの部分を付与
		}

		_, err = io.CopyN(file, actualReader, bytesToWrite)
		if err != nil {
			file.Close()
			return err
		}
		file.Close()
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

	// 指定したオプション数
	optionCount := 0
	// 分割モード
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

	// lineCount, chunkCount, byteCount の値が、デフォルト値以外の値になっているものが2つ以上あったらエラーにする
	if optionCount > 1 {
		fmt.Fprintln(os.Stderr, "Please specify only one option")
		flag.Usage()
		return
	}

	var filename string
	if args := flag.Args(); len(args) > 0 {
		filename = args[0]
	}

	var reader io.Reader
	switch filename {
	case "":
		reader = os.Stdin
	default:
		f, err := os.Open(filename)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		defer f.Close()
		reader = f
	}

	var outputDirName string = ""

	if err := mode.Split(reader, outputDirName); err != nil {
		fmt.Fprintf(os.Stderr, "fail split file: %v\n", err)
	}
}
