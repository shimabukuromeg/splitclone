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

	fmt.Print(flag.Args())
	fmt.Print("\n=========\n")
	fmt.Printf("line: %d\n", lineCount)
	fmt.Printf("chunkCount: %d\n", chunkCount)
	fmt.Printf("byteCount: %d\n", byteCount)

	cli := &split.CLI{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	}

	args := []string{"a", "b", "c"}

	cli.Run(args)

	// 指定したオプション数
	optionCount := 0

	/*
		flagをパースした結果から分割方法を確定するロジック
		TODO: もっといい感じに実装できそうな気がする。
	*/
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

	/*
		optionが複数指定されてたらエラーにする処理。
		optionが複数指定されてるか判断する方法が、デフォルト値以外の値になってるオプションの数で判断してしまってるけど
		普通にコマンドの引数で渡されてるフラグをチェックしたい。
		flagパッケージでデフォルト値を指定するので、最初から値が入ってて、入ってない場合指定されてない、みたいな判定ができない。
	*/
	if optionCount > 1 {
		fmt.Fprintln(os.Stderr, "Please specify only one option")
		flag.Usage()
		return
	}

	/*
		ファイル名を取得してる。一番最初決めうちにしてしまってるけどもうちょっとちゃんとしたチェックした方が良さそう？
	*/
	var filename string
	if args := flag.Args(); len(args) > 0 {
		filename = args[0]
	}

	/*
		指定されたファイル名から、ファイルを読み込んで io.Reader型の変数に格納
		NOTE: io.Reader型とは、「何かを読み込む機能を持つものをまとめて扱うために抽象化されたもの」
	*/
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
