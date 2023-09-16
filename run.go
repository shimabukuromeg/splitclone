package split

import (
	"flag"
	"fmt"
	"io"
	"os"
)

type CLI struct {
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader
}

const (
	DefaultLineCount  = 1000
	DefaultChunkCount = 0
	DefaultByteCount  = 0
)

var (
	commandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	lineCount   int
	chunkCount  int
	byteCount   int64
)

func (cli *CLI) Run(args []string) int {
	// 指定したオプション数
	optionCount := 0

	commandLine.IntVar(&lineCount, "l", DefaultLineCount, "line_count [file]")
	commandLine.IntVar(&chunkCount, "n", DefaultChunkCount, "chunk_count [file]")
	commandLine.Int64Var(&byteCount, "b", DefaultByteCount, "byte_count [file]")

	if err := commandLine.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "can not parse flags: %v\n", err)
	}

	/*
		flagをパースした結果から分割方法を確定するロジック
		TODO: もっといい感じに実装できそうな気がする。
	*/
	var spliter Spliter = LineSpliter{LineCount: DefaultLineCount}
	if lineCount != DefaultLineCount {
		optionCount++
		spliter = LineSpliter{LineCount: lineCount}
	}
	if chunkCount != DefaultChunkCount {
		optionCount++
		spliter = ChunkSpliter{ChunkCount: chunkCount}
	}
	if byteCount != DefaultByteCount {
		optionCount++
		spliter = ByteSpliter{ByteCount: byteCount}
	}

	/*
		optionが複数指定されてたらエラーにする処理。
		optionが複数指定されてるか判断する方法が、デフォルト値以外の値になってるオプションの数で判断してしまってるけど
		普通にコマンドの引数で渡されてるフラグをチェックしたい。
		flagパッケージでデフォルト値を指定するので、最初から値が入ってて、入ってない場合指定されてない、みたいな判定ができない。
	*/
	if optionCount > 1 {
		// 指定されたメッセージを標準エラー出力に出力する。
		fmt.Fprintln(os.Stderr, "複数のオプションは選択できません")
		commandLine.Usage()
		return 1
	}

	/*
		ファイル名を取得してる。一番最初決めうちにしてしまってるけどもうちょっとちゃんとしたチェックした方が良さそう？
	*/
	var filename string
	if args := commandLine.Args(); len(args) > 0 {
		filename = args[0]
	} else {
		// 指定されたメッセージを標準エラー出力に出力する。
		fmt.Fprintln(os.Stderr, "ファイルを指定してください")
		commandLine.Usage()
		return 0
	}

	/*
		指定されたファイル名から、ファイルを読み込んで io.Reader型の変数に格納
		io.Reader型とは、「何かを読み込む機能を持つものをまとめて扱うために抽象化されたもの」
	*/
	var reader io.Reader
	switch filename {
	case "":
		reader = os.Stdin
	default:
		f, err := os.Open(filename)
		if err != nil {
			// 指定されたメッセージを標準エラー出力に出力する。
			fmt.Fprintln(os.Stderr, err)
		}
		defer f.Close()
		reader = f
	}

	var outputDirName string = ""

	if err := spliter.Split(reader, outputDirName); err != nil {
		// 指定されたメッセージを標準エラー出力に出力する。
		fmt.Fprintf(os.Stderr, "fail split file: %v\n", err)
		return 1
	}

	return 0
}
