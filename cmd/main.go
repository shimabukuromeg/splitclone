package main

import (
	"flag"
	"os"

	split "splitclone"
)

var (
	lineCount  int
	chunkCount int
	byteCount  int64
)

func init() {
	flag.IntVar(&lineCount, "l", split.DefaultLineCount, "line_count [file]")
	flag.IntVar(&chunkCount, "n", split.DefaultChunkCount, "chunk_count [file]")
	flag.Int64Var(&byteCount, "b", split.DefaultByteCount, "byte_count [file]")
}

func main() {
	flag.Parse()

	cli := &split.CLI{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	}

	cli.Run(lineCount, chunkCount, byteCount)
}
