package repl

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	goprompt "github.com/ktr0731/go-prompt"
	split "github.com/shimabukuromeg/splitclone"
	"github.com/tj/go-spin"
)

type REPL struct {
	writer io.Writer
}

func NewRepl() (*REPL, error) {
	r := &REPL{writer: os.Stdout}
	return r, nil
}

func (r *REPL) Run() int {
	r.printSplash()
	defer fmt.Fprintln(r.writer, "\n👋 Good Bye :)")

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		return 1
	}

	fmt.Println("📂 Selected file")

	var ops []goprompt.Option

	var file string
	for {
		file, err = goprompt.Input(currentDir+"> ", fileCompleter, append(
			ops,
			goprompt.OptionPrefixTextColor(goprompt.Color(goprompt.Blue)),
		)...)
		if errors.Is(err, goprompt.ErrAbort) {
			return 1
		} else if err != nil {
			return 1
		}
		if file != "" {
			break
		}
		fmt.Println("❌ No file selected! Please choose a file.")
	}

	fmt.Printf("✅ Selected file: \033[32m%s\033[0m\n", file)

	// NOTE: 分割する方法を選ぶ（行数・分割数・バイト数）
	var mode string
	for {
		// mode = p2.Input()
		mode, err = goprompt.Input("Please choose a split method > ", completer, append(
			ops,
			goprompt.OptionPrefixTextColor(goprompt.Color(goprompt.Blue)),
		)...)
		if err != nil {
			return 1
		}
		if mode != "" {
			break
		}
		fmt.Println("❌ No split method selected! Please choose split method.")
	}

	fmt.Printf("✅ Your split method: \033[32m%s\033[0m\n", mode)

	f, err := os.Open(file)
	if err != nil {
		// 指定されたメッセージを標準エラー出力に出力する。
		fmt.Fprintln(os.Stderr, err)
	}
	defer f.Close()
	reader := f

	// TODO: number だけじゃなくて、line,byte　の条件も追加する
	if mode == "number" {
		for {
			number, err := goprompt.Input("count (TYPE_NUMBER) => ", dummyCompleter, append(
				ops,
				goprompt.OptionPrefixTextColor(goprompt.Color(goprompt.Green)),
			)...)
			if err != nil {
				return 1
			}

			count, err := strconv.Atoi(number)
			if err == nil {
				fmt.Println("Received number:", count)
				var spliter split.Spliter = split.ChunkSpliter{ChunkCount: count}
				if err := spliter.Split(reader, ""); err != nil {
					// 指定されたメッセージを標準エラー出力に出力する。
					fmt.Fprintf(os.Stderr, "fail split file: %v\n", err)
					return 1
				}
				break
			} else {
				fmt.Println("❌ Please enter a valid number.")
			}
		}
	}

	// NOTE: spinさせて雰囲気出した
	s := spin.New()
	for i := 0; i < 20; i++ {
		fmt.Printf("\r  \033[36msplitting\033[m %s ", s.Next())
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Printf("\n✅ Complete\n")

	return 1
}

// dummyCompleter は常に空のサジェスチョンリストを返します。
func dummyCompleter(in goprompt.Document) []goprompt.Suggest {
	return []goprompt.Suggest{}
}

func completer(in goprompt.Document) []goprompt.Suggest {
	// 入力の末尾がスペースかどうかをチェック
	if strings.HasSuffix(in.Text, " ") {
		return []goprompt.Suggest{} // スペースの後は何もサジェスチョンしない
	}

	s := []goprompt.Suggest{
		{Text: "number", Description: "Split by specified number"},
		{Text: "line", Description: "Split by specified number of lines"},
		{Text: "byte", Description: "Split by specified size in bytes"},
	}

	return goprompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)
}

func (r *REPL) printSplash() {
	fmt.Fprintln(r.writer, defaultSplashText)
}

func fileCompleter(d goprompt.Document) []goprompt.Suggest {
	files, _ := listFiles()
	var s []goprompt.Suggest
	for _, file := range files {
		s = append(s, goprompt.Suggest{Text: file})
	}
	return goprompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func listFiles() ([]string, error) {
	var files []string
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}
