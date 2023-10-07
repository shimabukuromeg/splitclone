package repl

import (
	"fmt"
	"io"
	"os"
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

	defer fmt.Fprintln(r.writer, "Good Bye :)")
	return 1
}

func (r *REPL) printSplash() {
	fmt.Fprintln(r.writer, defaultSplashText)
}
