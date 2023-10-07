package repl

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/c-bata/go-prompt"
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

	file := p.Input()

	fmt.Println("✅ Selected file:", file)

	return 1
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
