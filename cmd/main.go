package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	split "splitclone"
)

var (
	lineCount  int
	chunkCount int
	byteCount  int64
)

const (
	defaultLineCount  = 1000
	defaultChunkCount = 0
	defaultByteCount  = 0
)

func init() {
	flag.IntVar(&lineCount, "l", defaultLineCount, "line_count [file]")
	flag.IntVar(&chunkCount, "n", defaultChunkCount, "chunk_count [file]")
	flag.Int64Var(&byteCount, "b", defaultByteCount, "byte_count [file]")
}

func main() {
	flag.Parse()

	cli := &split.CLI{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	}

	args := []string{"a", "b", "c"}

	cli.Run(args)

	// 指定したオプション数
	optionCount := 0

	var spliter split.Spliter = split.LineSpliter{LineCount: defaultLineCount}
	if lineCount != defaultLineCount {
		optionCount++
		spliter = split.LineSpliter{LineCount: lineCount}
	}
	if chunkCount != defaultChunkCount {
		optionCount++
		spliter = split.ChunkSpliter{ChunkCount: chunkCount}
	}
	if byteCount != defaultByteCount {
		optionCount++
		spliter = split.ByteSpliter{ByteCount: byteCount}
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

	if err := spliter.Split(reader, outputDirName); err != nil {
		fmt.Fprintf(os.Stderr, "fail split file: %v\n", err)
	}
}
