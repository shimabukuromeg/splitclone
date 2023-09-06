package split

import (
	"flag"
	"fmt"
	"io"
)

type CLI struct {
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader
}

func (cli *CLI) Run(args []string) error {
	fmt.Println("exec split command ")

	var test int
	flag.IntVar(&test, "t", 100, "test")
	flag.Parse()

	// fmt.Printf("test is %d\n", test)

	for _, v := range args {
		fmt.Println(v)
	}

	return nil
}
