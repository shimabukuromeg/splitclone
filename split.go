package split

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/debug"
)

type Spliter interface {
	Split(io.Reader, string) error
}

type LineSpliter struct {
	LineCount int
}

type ByteSpliter struct {
	ByteCount int64
}

type ChunkSpliter struct {
	ChunkCount int
}

// スタックトレースで使う構造体
type ErrorWithStack struct {
	err   error
	stack []byte
}

func (e *ErrorWithStack) Error() string {
	return e.err.Error() + string(e.stack)
}

func WrapWithStack(err error) error {
	stk := debug.Stack()
	return &ErrorWithStack{
		err:   err,
		stack: stk,
	}
}

func (ls LineSpliter) Split(reader io.Reader, outputDirName string) error {
	scanner := bufio.NewScanner(reader)

	var lineCount int = 0
	fileIndex := 1
	var file *os.File

	for scanner.Scan() {
		if lineCount%ls.LineCount == 0 {
			if file != nil {
				file.Close()
				file = nil
			}

			// 書き出しファイル
			// filename := fmt.Sprintf("part-line-%d", fileIndex)
			filename := genFileName(fileIndex)
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
		// 文字列ベースで簡単なエラーの作り方
		// https://docs.google.com/presentation/d/1HW3wG8J_Q2536Iu__7HGr_mhurHajC7IOGjCnn3kZmg/edit#slide=id.g7e7c89dcc9_0_166
		return fmt.Errorf("reading input: %w", err)
	}

	return nil
}

func (bs ByteSpliter) Split(reader io.Reader, outputDirName string) error {
	// 指定のバイト数が0より小さい場合はエラー
	if bs.ByteCount <= 0 {
		return WrapWithStack(fmt.Errorf("invalid ByteCount value: %d\n ", bs.ByteCount))
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

	chunkCount := int(totalSize) / int(bs.ByteCount)
	if int(totalSize)%int(bs.ByteCount) != 0 {
		chunkCount++
	}

	i := 1
	value := totalSize

	for value > 0 {
		bytesToWrite := bs.ByteCount
		if value < bs.ByteCount {
			bytesToWrite = value
		}
		filename := genFileName(i)
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
		value -= bs.ByteCount
		i++
	}

	return nil
}

func (cs ChunkSpliter) Split(reader io.Reader, outputDirName string) error {
	// 指定の数が0より小さい場合はエラー
	if cs.ChunkCount <= 0 {
		return fmt.Errorf("invalid ChunkCount value: %d", cs.ChunkCount)
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

	ByteCount := totalSize / int64(cs.ChunkCount)
	extraBytes := totalSize % int64(cs.ChunkCount) // 追加で余りを取得

	for i := 0; i < cs.ChunkCount; i++ {
		filename := genFileName(i + 1)
		file, err := os.Create(filepath.Join(outputDirName, filename))
		if err != nil {
			return err
		}

		bytesToWrite := ByteCount
		if i == cs.ChunkCount-1 {
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
