package repl

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/c-bata/go-prompt"
	split "github.com/shimabukuromeg/splitclone"
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
	defer fmt.Fprintln(r.writer, "👋 Good Bye :)")

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		return 0
	}

	fmt.Println("📂 Selected file")

	p := prompt.New(
		func(s string) {},
		fileCompleter,
		prompt.OptionPrefix(currentDir+"> "),
		prompt.OptionPrefixTextColor(prompt.Blue),
	)

	var file string
	for {
		file = p.Input()
		if file != "" {
			break
		}
		fmt.Println("❌ No file selected! Please choose a file.")
	}

	fmt.Println("✅ Selected file:", file)

	// TODO: 分割する方法を選ぶ（行数・分割数・バイト数）
	p2 := prompt.New(
		func(s string) {},
		completer,
		prompt.OptionPrefix("Please choose a split method > "),
		prompt.OptionPrefixTextColor(prompt.Blue),
	)

	var mode string
	for {
		mode = p2.Input()
		if mode != "" {
			break
		}
		fmt.Println("❌ No split method selected! Please choose split method.")
	}

	fmt.Println("Your split method:", mode)

	p3 := prompt.New(
		func(s string) {},
		dummyCompleter,
		prompt.OptionPrefix("count (TYPE_NUMBER) => "),
		prompt.OptionPrefixTextColor(prompt.Green),
	)

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
			count, err := strconv.Atoi(p3.Input())
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

	fmt.Println("✅ Complete")

	return 1
}

// dummyCompleter は常に空のサジェスチョンリストを返します。
func dummyCompleter(in prompt.Document) []prompt.Suggest {
	return []prompt.Suggest{}
}

func completer(in prompt.Document) []prompt.Suggest {
	// 入力の末尾がスペースかどうかをチェック
	if strings.HasSuffix(in.Text, " ") {
		return []prompt.Suggest{} // スペースの後は何もサジェスチョンしない
	}

	s := []prompt.Suggest{
		{Text: "number", Description: "Split by specified number"},
		{Text: "line", Description: "Split by specified number of lines"},
		{Text: "byte", Description: "Split by specified size in bytes"},
	}

	return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)
}

func (r *REPL) printSplash() {
	fmt.Fprintln(r.writer, defaultSplashText)
}

func fileCompleter(d prompt.Document) []prompt.Suggest {
	files, _ := listFiles()
	var s []prompt.Suggest
	for _, file := range files {
		s = append(s, prompt.Suggest{Text: file})
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
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
